package plugin

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Meta struct {
	TableName           string
	MainTypeCapitalized string
	MainTypeLowercased  string
	Types               []Fields
}

type Fields struct {
	Name               string
	GoType             string
	MainTypeLowercased string
}

func GenerateFileForType(
	typeToGenerate Meta,
	template *template.Template,
) ([]byte, error) {
	var buf bytes.Buffer

	if err := template.Execute(&buf, typeToGenerate); err != nil {
		return nil, fmt.Errorf("failed to generate type using given template: %w", err)
	}

	return buf.Bytes(), nil
}

func MetaFromProtobufType(message *protogen.Message) (Meta, error) {
	tableName := PostgresTableNameForMessage(message)

	capitalized := cases.Title(language.English).String(tableName)
	lowercased := cases.Lower(language.English).String(tableName)

	opts := Meta{
		TableName:           tableName,
		MainTypeCapitalized: capitalized,
		MainTypeLowercased:  lowercased,
		Types:               []Fields{},
	}

	for _, field := range message.Fields {
		name := string(field.Desc.Name())
		gotype, err := PostgresTypeFromProtobufType(field)
		if err != nil {
			return opts, fmt.Errorf("failed to map message %s to postgres type: %w", message.Desc.Name(), err)
		}
		opts.Types = append(opts.Types, Fields{
			Name:               name,
			GoType:             gotype,
			MainTypeLowercased: lowercased,
		})
	}

	return opts, nil
}

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
	return strings.ToLower(string(message.Desc.Name()))
}
