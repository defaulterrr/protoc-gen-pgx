package plugin

import (
	"fmt"
	"strings"

	pb "github.com/defaulterrr/protoc-gen-pgx/pb/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func PostgresTypeFromProtobufType(pbtype *protogen.Field) (string, error) {
	switch pbtype.Desc.Kind() {
	case protoreflect.BoolKind:
		return "bool", nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32", nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32", nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64", nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64", nil
	case protoreflect.FloatKind:
		return "float32", nil
	case protoreflect.DoubleKind:
		return "float64", nil
	case protoreflect.StringKind:
		return "string", nil
	case protoreflect.BytesKind:
		return "[]byte", nil
	default:
		return "", fmt.Errorf("unsupported type: %s", pbtype.Desc.Kind().String())
	}
}

func PostgresTableNameForMessage(message *protogen.Message) string {
	tableName, ok := getTableNameFromOptions(message.Desc.Options().(*descriptorpb.MessageOptions))
	if ok {
		return tableName
	}
	return strings.ToLower(string(message.Desc.Name()))
}

func getTableNameFromOptions(in *descriptorpb.MessageOptions) (string, bool) {
	if in == nil {
		return "", false
	}

	v := proto.GetExtension(in, pb.E_TableName)
	if v == nil {
		return "", false
	}

	shouldGenerate, ok := v.(string)
	return shouldGenerate, ok
}
