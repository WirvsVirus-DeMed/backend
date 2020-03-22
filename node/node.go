package node

import (
	"WirvsVirus/DeMed/protobuf"
	"WirvsVirus/DeMed/util"
	"bufio"
	"container/list"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
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
	publicIp        net.IP
	localPort       uint32
	connectionLimit uint32

	config      *tls.Config
	listener    net.Listener
	Clients     list.List
	clientMutex sync.Mutex

	IncomingMessageQueue  chan *MessageDescriptor
	BroadcastMessageQueue chan *protobuf.MessageFrame
}

func (this *Node) Init(localCertFile, localKeyFile, caCertFile, serverName string, port, connectionLimit, bufferedMessages uint32) {
	this.Clients.Init()
	this.publicIp = util.GetPublicIp()
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
		//RootCAs:      rootCAs,
		ClientCAs:  rootCAs,
		ServerName: serverName,
		ClientAuth: tls.RequireAndVerifyClientCert,
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
			this.clientMutex.Lock()
			for e := this.Clients.Front(); e != nil; e = e.Next() {
				c := e.Value.(*Info)

				if c == info {
					this.Clients.Remove(e)
					return
				}
			}
			this.clientMutex.Unlock()
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
			log.Fatalf("Error while accepting client: %v", err)
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
		//ClientCAs:		this.config.ClientCAs,
		ServerName: "DeMed-Node",
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	conn, err := tls.Dial("tcp", ip.String()+":"+strconv.Itoa(int(port)), &cfg)
	if err != nil {
		log.Fatalf("Could not connect to %v. Reason: %v", ip.String()+":"+strconv.Itoa(int(port)), err)
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

				if md.origin.RemotePubKey.N.Cmp(c.RemotePubKey.N) == 0 && md.origin.RemotePubKey.E == c.RemotePubKey.E {
					md.origin.CloseConnection()
					clientsDirty = true
					break
				}
			}

			if clientsDirty {
				for e := this.Clients.Front(); e != nil; e = e.Next() {
					c := e.Value.(*Info)

					if c.ConnectionState == Closed {
						this.Clients.Remove(e)
					}
				}
			}
			this.clientMutex.Unlock()

			if md.origin.ConnectionState == Closed {
				continue
			}

			s := md.origin.RemotePubKey.N.String()

			log.Printf("[*] %v from %v", s[0:4]+"..."+s[len(s)-5:len(s)-1], md.origin.Peer.String())

			go func(info *Info, n *Node) {
				n.clientMutex.Lock()

				for e := n.Clients.Front(); e != nil; e = e.Next() {
					c := e.Value.(*Info)

					if info == c {
						continue
					}
					if (c.Peer != nil && (c.Peer.IP.Equal(info.Peer.IP) && c.Peer.Port == info.Peer.Port)) ||
						c.RemotePeer.IP.Equal(info.Peer.IP) {
						c.CloseConnection()
						n.Clients.Remove(e)
					}
				}
				n.clientMutex.Unlock()
			}(md.origin, this)

			md.origin.ConnectionState += Established
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
