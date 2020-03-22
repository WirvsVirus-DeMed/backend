package node

import "github.com/WirvsVirus-DeMed/backend/protobuf"

type MessageDescriptor struct {
	msgFrame *protobuf.MessageFrame
	origin   *Info
}
