package node

import (
	"bufio"
	"crypto/rsa"
	"github.com/WirvsVirus-DeMed/backend/protobuf"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"log"
	"net"
	"time"
)

type Info struct {
	RemotePeer *net.TCPAddr
	Peer       *net.TCPAddr
	Socket     net.Conn
	r          *bufio.Reader
	w          *bufio.Writer

	LastSeen        time.Time
	ConnectionState ConnectionState
	incoming        bool
	RemotePubKey    rsa.PublicKey
}

func (this *Info) CloseConnection() {
	if this.ConnectionState == Closed {
		return
	}

	this.ConnectionState = Closed
	this.Socket.Close()

	if this.Peer != nil {
		log.Println("[*] Client disconnected - " + this.Peer.String())
	} else {
		log.Println("[*] Client disconnected - " + this.RemotePeer.String())
	}
}

func (this *Info) SendPeerRequest() {
	this.SendMessageFrame(nil, PeerRequest)
}
func (this *Info) SendPeerResponse(message protobuf.PeerResponse) {
	data, _ := proto.Marshal(&message)

	this.SendMessageFrame(data, PeerResponse)
}
func (this *Info) SendHello(message *protobuf.Hello) {
	data, _ := proto.Marshal(message)

	this.SendMessageFrame(data, Hello)
}
func (this *Info) SendGoodbye() {
	this.SendMessageFrame(nil, Goodbye)
}
func (this *Info) SendMessageBroadcast(message *protobuf.MessageBroadcast) {
	data, _ := proto.Marshal(message)

	this.SendMessageFrame(data, MessageBroadcast)
}
func (this *Info) SendRequestResource(message *protobuf.RequestBroadcast) {
	data, _ := proto.Marshal(message)

	this.SendMessageFrame(data, RequestResource)
}
func (this *Info) SendMessageFrame(payloadBytes []byte, payloadType PayloadType) {
	pbFrame := protobuf.MessageFrame{
		PayloadType: uint32(payloadType),
		Time:        ptypes.TimestampNow(),
		Payload:     payloadBytes,
	}

	frameData, _ := proto.Marshal(&pbFrame)

	if this.ConnectionState != Closed {
		_, _ = this.w.Write(frameData)
		_ = this.w.Flush()
	}
}
