package node

import (
	"bufio"
	"container/list"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"github.com/WirvsVirus-DeMed/backend/db"
	"github.com/WirvsVirus-DeMed/backend/protobuf"
	"github.com/WirvsVirus-DeMed/backend/util"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"strconv"
	"sync"
	"time"
)

const DeMedVersion = "alpha_0.1"

type Node struct {
	publicIp          net.IP
	localPort         uint32
	connectionLimit   uint32
	discoveryInterval time.Duration

	config             *tls.Config
	listener           net.Listener
	Clients            list.List
	clientMutex        sync.Mutex
	PeerBlackList      list.List
	PeerBlackListMutex sync.Mutex

	IncomingMessageQueue  chan *MessageDescriptor
	BroadcastMessageQueue chan *protobuf.MessageFrame
}

func (this *Node) Init(localCertFile, localKeyFile, caCertFile, serverName string, port, connectionLimit, bufferedMessages uint32, discoveryInterval time.Duration) {
	this.Clients.Init()
	this.PeerBlackList.Init()
	this.publicIp = util.GetPublicIp()
	this.discoveryInterval = discoveryInterval
	this.connectionLimit = connectionLimit
	this.localPort = port
	this.IncomingMessageQueue = make(chan *MessageDescriptor, bufferedMessages)
	this.BroadcastMessageQueue = make(chan *protobuf.MessageFrame, bufferedMessages)

	cert, err := tls.LoadX509KeyPair(localCertFile, localKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	rootCaCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", caCertFile, err)
	}

	rootCAs := x509.NewCertPool()
	if ok := rootCAs.AppendCertsFromPEM(rootCaCert); !ok {
		log.Fatalf("Failed to append %q to RootCAs: %v", caCertFile, err)
	}

	this.config = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    rootCAs,
		ServerName:   serverName,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
}

func (this *Node) HandleConnection(info *Info) {
	defer info.CloseConnection()

	this.clientMutex.Lock()
	for e := this.Clients.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Info)
		if c.RemotePeer.IP.Equal(info.RemotePeer.IP) && c.RemotePeer.Port == info.RemotePeer.Port {
			this.clientMutex.Unlock()
			return
		}
	}
	this.Clients.PushBack(info)
	this.clientMutex.Unlock()

	info.r = bufio.NewReader(info.Socket)
	info.w = bufio.NewWriter(info.Socket)

	if !info.incoming {
		log.Printf("[*] Connected to %s\n", info.RemotePeer.String())

		cert, _ := x509.ParseCertificate(this.config.Certificates[0].Certificate[0])
		pubKey := cert.PublicKey.(*rsa.PublicKey)

		pk := protobuf.Hello_PublicKey{
			N: pubKey.N.Bytes(),
			E: int32(pubKey.E),
		}

		pb := protobuf.Hello{
			DeMedVersion: DeMedVersion,
			Ip:           this.publicIp,
			Port:         this.localPort,
			PubKey:       &pk,
		}

		info.SendHello(&pb)
	} else {
		log.Printf("[*] Incoming Connection %s\n", info.RemotePeer.String())
	}

	buffer := make([]byte, 4096)

	for {
		n, err := info.r.Read(buffer)
		if n == 0 {
			log.Printf("[x] %s connection lost\n", info.RemotePeer.String())
			this.clientMutex.Lock()
			for e := this.Clients.Front(); e != nil; e = e.Next() {
				c := e.Value.(*Info)

				if c == info {
					this.Clients.Remove(e)
				}
			}
			this.clientMutex.Unlock()
			return
		} else if err != nil {
			log.Println(err)
			pp := PeerPunish{}
			this.clientMutex.Lock()
			for e := this.Clients.Front(); e != nil; e = e.Next() {
				c := e.Value.(*Info)

				if c == info {
					this.Clients.Remove(e)
					pp = PeerPunish{
						Peer:  *c.Peer,
						until: time.Now().Add(10 * time.Minute),
					}
					break
				}
			}
			this.clientMutex.Unlock()

			this.PeerBlackListMutex.Lock()
			this.PeerBlackList.PushBack(pp)
			this.PeerBlackListMutex.Unlock()
			return
		}

		msgFrame := new(protobuf.MessageFrame)
		err = proto.Unmarshal(buffer[0:n], msgFrame)
		if err == nil {
			log.Printf("[>] recieved from %s: %v\n", info.RemotePeer.String(), msgFrame.PayloadType)
			info.LastSeen = time.Now()
			this.IncomingMessageQueue <- &MessageDescriptor{msgFrame: msgFrame, origin: info}
		} else {
			log.Printf("[*] Error while decoding MessageFrame from %s - Reason: %v\n", info.RemotePeer.String(), err)
		}
	}
}

func (this *Node) Listen() {
	var err = error(nil)
	this.listener, err = tls.Listen("tcp", ":"+strconv.FormatUint(uint64(this.localPort), 10), this.config)
	if err != nil {
		log.Fatalf("Error while starting listener: %v", err)
	}
	defer this.listener.Close()

	for {
		conn, err := this.listener.Accept()
		if err != nil {
			log.Printf("Error while accepting client: %v", err)
		}

		info := Info{
			RemotePeer:      conn.RemoteAddr().(*net.TCPAddr),
			Socket:          conn,
			LastSeen:        time.Now(),
			ConnectionState: Established,
			incoming:        true,
		}

		go this.HandleConnection(&info)
	}
}

func (this *Node) Connect(ip net.IP, port uint32) {
	cfg := tls.Config{
		Certificates: this.config.Certificates,
		RootCAs:      this.config.ClientCAs,
		ServerName:   "DeMed-Node",
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	tcpVersion := "tcp"
	if ip.To4() != nil {
		tcpVersion += "4"
	} else {
		tcpVersion += "6"
	}

	conn, err := tls.Dial(tcpVersion, "["+ip.String()+"]:"+strconv.Itoa(int(port)), &cfg)
	if err != nil {
		log.Printf("Could not connect to %v. Reason: %v", ip.String()+":"+strconv.Itoa(int(port)), err)
		this.PeerBlackListMutex.Lock()
		this.PeerBlackList.PushBack(PeerPunish{
			Peer: net.TCPAddr{
				IP:   ip,
				Port: int(port),
			},
			until: time.Now().Add(10 * time.Minute),
		})
		this.PeerBlackListMutex.Unlock()
		return
	}

	info := Info{
		RemotePeer:      conn.RemoteAddr().(*net.TCPAddr),
		Socket:          conn,
		LastSeen:        time.Now(),
		ConnectionState: Established,
		incoming:        false,
	}

	go this.HandleConnection(&info)
}

func (this *Node) BroadcastSender() {
	for {
		msg := <-this.BroadcastMessageQueue
		if msg == nil {
			return
		}

		frameData, _ := proto.Marshal(msg)

		this.clientMutex.Lock()
		for e := this.Clients.Front(); e != nil; e = e.Next() {
			c := e.Value.(*Info)
			c.w.Write(frameData)
			c.w.Flush()
		}
		this.clientMutex.Unlock()
	}
}

func (this *Node) PeerDiscovery() {
	for {
		this.PeerBlackListMutex.Lock()
		for e := this.PeerBlackList.Front(); e != nil; e = e.Next() {
			p := e.Value.(PeerPunish)

			if p.until.Before(time.Now()) {
				this.PeerBlackList.Remove(e)
			}
		}
		this.PeerBlackListMutex.Unlock()

		peers, _ := db.GetAllPeers(db.CurrentDb)

		blackListShadow := new(list.List)
		blackListShadow.Init()
		this.PeerBlackListMutex.Lock()
		for e := this.PeerBlackList.Front(); e != nil; e = e.Next() {
			pp := e.Value.(PeerPunish)
			blackListShadow.PushBack(pp)
		}
		this.PeerBlackListMutex.Unlock()

		this.clientMutex.Lock()
		for _, p := range peers {
			connect := true

			for e := blackListShadow.Front(); e != nil; e = e.Next() {
				pp := e.Value.(PeerPunish)

				if p.Port == uint32(pp.Peer.Port) && p.IP.Equal(pp.Peer.IP) {
					connect = false
					break
				}
			}
			if !connect {
				continue
			}

			for e := this.Clients.Front(); e != nil; e = e.Next() {
				c := e.Value.(*Info)

				if (c.Peer != nil && (c.Peer.IP.Equal(p.IP) && uint32(c.Peer.Port) == p.Port)) ||
					(c.RemotePeer.IP.Equal(p.IP) && uint32(c.RemotePeer.Port) == p.Port) {
					connect = false
					break
				}
			}

			if connect {
				go this.Connect(p.IP, p.Port)
			}
		}
		this.clientMutex.Unlock()

		time.Sleep(this.discoveryInterval)
	}
}

func (this *Node) HandleMessages() {
	for {
		md := <-this.IncomingMessageQueue
		if md == nil {
			return
		}

		switch PayloadType(md.msgFrame.PayloadType) {
		case Hello:
			hello := protobuf.Hello{}
			err := proto.Unmarshal(md.msgFrame.Payload, &hello)
			if err != nil {
				log.Printf("[*] Error while unmarshalling PeerResponse: %v\n", err)
			}

			md.origin.Peer = &net.TCPAddr{
				IP:   hello.Ip,
				Port: int(hello.Port),
			}

			md.origin.RemotePubKey = rsa.PublicKey{
				N: new(big.Int).SetBytes(hello.PubKey.N),
				E: int(hello.PubKey.E),
			}

			cert, _ := x509.ParseCertificate(this.config.Certificates[0].Certificate[0])
			pubKey := cert.PublicKey.(*rsa.PublicKey)

			clientsDirty := false
			this.clientMutex.Lock()

			if pubKey.N.Cmp(md.origin.RemotePubKey.N) == 0 && pubKey.E == md.origin.RemotePubKey.E {

				for e := this.Clients.Front(); e != nil; e = e.Next() {
					c := e.Value.(*Info)

					if c == md.origin {
						md.origin.CloseConnection()
						clientsDirty = true
						break
					}
				}
			}

			for e := this.Clients.Front(); e != nil; e = e.Next() {
				c := e.Value.(*Info)

				if c == md.origin || c.RemotePubKey.N == nil {
					continue
				}

				if md.origin.RemotePubKey.N.Cmp(c.RemotePubKey.N) == 0 &&
					md.origin.RemotePubKey.E == c.RemotePubKey.E ||
					(c.Peer != nil && (c.Peer.IP.Equal(md.origin.Peer.IP) &&
						c.Peer.Port == md.origin.Peer.Port)) ||
					c.RemotePeer.IP.Equal(md.origin.Peer.IP) {

					md.origin.CloseConnection()
					clientsDirty = true
					continue
				}
			}

			for e := this.Clients.Front(); e != nil; e = e.Next() {
				c := e.Value.(*Info)

				if md.origin == c {
					continue
				}
				if (c.Peer != nil && (c.Peer.IP.Equal(md.origin.Peer.IP) && c.Peer.Port == md.origin.Peer.Port)) ||
					c.RemotePeer.IP.Equal(md.origin.Peer.IP) {
					c.CloseConnection()
					this.Clients.Remove(e)
				}
			}

			blacklistShadow := new(list.List)
			if clientsDirty {
				blacklistShadow.Init()
				for e := this.Clients.Front(); e != nil; e = e.Next() {
					c := e.Value.(*Info)

					if c.ConnectionState == Closed {
						this.Clients.Remove(e)
						blacklistShadow.PushBack(PeerPunish{
							Peer:  *c.Peer,
							until: time.Now().Add(30 * time.Minute),
						})
					}
				}
			}
			this.clientMutex.Unlock()

			this.PeerBlackListMutex.Lock()
			for e := blacklistShadow.Front(); e != nil; e = e.Next() {
				c := e.Value.(PeerPunish)
				this.PeerBlackList.PushBack(c)
			}
			this.PeerBlackListMutex.Unlock()

			if md.origin.ConnectionState == Closed {
				continue
			}

			s := md.origin.RemotePubKey.N.String()

			log.Printf("[*] %v from %v >> Connection Established", s[0:4]+"..."+s[len(s)-5:len(s)-1], md.origin.Peer.String())

			md.origin.ConnectionState += Established
			p := db.Peer{
				IP:       md.origin.Peer.IP,
				Port:     uint32(md.origin.Peer.Port),
				LastSeen: time.Now(),
			}
			p.Add(db.CurrentDb)

			if md.origin.incoming {
				pk := protobuf.Hello_PublicKey{
					N: pubKey.N.Bytes(),
					E: int32(pubKey.E),
				}

				pb := protobuf.Hello{
					DeMedVersion: DeMedVersion,
					Ip:           this.publicIp,
					Port:         this.localPort,
					PubKey:       &pk,
				}

				md.origin.SendHello(&pb)
			}
			break

		case Goodbye:
			this.clientMutex.Lock()
			for e := this.Clients.Front(); e != nil; e = e.Next() {
				c := e.Value.(*Info)

				if c == md.origin {
					c.CloseConnection()
					this.Clients.Remove(e)
				}
			}
			this.clientMutex.Unlock()
			break

		case MessageBroadcast:
			log.Println("MessageBroadcast not implemented")
			break

		case RequestResource:
			log.Println("RequestResource not implemented")
			break

		case PeerRequest:
			this.clientMutex.Lock()
			peers := make([]*protobuf.PeerResponse_Peer, this.Clients.Len())

			i := 0
			for e := this.Clients.Front(); e != nil; e = e.Next() {
				c := e.Value.(*Info)
				peers[i] = &protobuf.PeerResponse_Peer{
					Ip:   c.RemotePeer.IP,
					Port: uint32(c.RemotePeer.Port),
				}

				i++
			}
			this.clientMutex.Unlock()

			msg := protobuf.PeerResponse{
				Peers: peers,
			}

			md.origin.SendPeerResponse(msg)
			break

		case PeerResponse:
			peerResp := protobuf.PeerResponse{}
			err := proto.Unmarshal(md.msgFrame.Payload, &peerResp)
			if err != nil {
				log.Printf("[*] Error while unmarshalling PeerResponse: %v\n", err)
			}

			for _, p := range peerResp.Peers {
				go this.Connect(p.Ip, p.Port)
			}

			break
		}
	}
}
