package grpc

import (
	"github.com/gogo/gateway"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"io"
)

// JsonPb marshaller wraps the better gogo-gateway marshaller in a way it can be used as golang JsonPb marshaller.
type JsonPbMarshaller struct {
	jsonpb.Marshaler
	gateway.JSONPb
}

func NewJsonPbMarshaller() *JsonPbMarshaller {
	return &JsonPbMarshaller{
		Marshaler: jsonpb.Marshaler{},
		JSONPb: gateway.JSONPb{
			EmitDefaults: true,
			Indent:       "  ",
			OrigName:     true,
		},
	}
}

func (m *JsonPbMarshaller) Marshal(out io.Writer, pb proto.Message) error {
	bt, err := m.JSONPb.Marshal(pb)
	if err != nil {
		return err
	}
	if _, err := out.Write(bt); err != nil {
		return err
	}
	return nil
}
