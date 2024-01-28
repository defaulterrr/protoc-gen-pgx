package plugin

import (
	pb "github.com/defaulterrr/protoc-gen-pgx/pb/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func FindAnnotatedMessages(messages []*protogen.Message) []*protogen.Message {
	annotatedMessages := []*protogen.Message{}

	for _, message := range messages {
		if shouldGenerate(message.Desc.Options().(*descriptorpb.MessageOptions)) {
			annotatedMessages = append(annotatedMessages, message)
		}
	}

	return annotatedMessages
}

func shouldGenerate(in *descriptorpb.MessageOptions) bool {
	if in == nil {
		return false
	}

	v := proto.GetExtension(in, pb.E_ShouldGenerate)
	if v == nil {
		return false
	}

	shouldGenerate, ok := v.(bool)
	if !ok {
		return false
	}
	return shouldGenerate
}
