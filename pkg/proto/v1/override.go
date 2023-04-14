package proto

import (
	"google.golang.org/protobuf/proto"
)

// Clone ...
func (m *Message) Clone() *Message {
	return proto.Clone(m).(*Message)
}
