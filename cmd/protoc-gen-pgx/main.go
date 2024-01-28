package main

import (
	_ "embed"
	"log/slog"
	"os"
	"slices"
	"text/template"

	"github.com/defaulterrr/protoc-gen-pgx/internal/plugin"
	"google.golang.org/protobuf/compiler/protogen"
)

//go:embed generated.template
var generatedTemplate string

func main() {
	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	slog.SetDefault(slog.New(textHandler))

	template, err := template.New("pgx").Parse(generatedTemplate)
	if err != nil {
		panic(err)
	}
	slog.Debug("Hello world!")

	protogen.Options{}.Run(func(p *protogen.Plugin) error {
		slog.Debug("Protogen plugin called with following files to be generated", "files", p.Request.FileToGenerate)

		for path, file := range p.FilesByPath {
			annotatedMessages := []*protogen.Message{}
			if !slices.Contains(p.Request.FileToGenerate, path) {
				slog.Debug("Skipping file", "file", path)
				continue
			}

			slog.Debug("Processing file", "file", path)
			annotatedMessagesInFile := plugin.FindAnnotatedMessages(file.Messages)

			annotatedMessages = append(annotatedMessages, annotatedMessagesInFile...)

			names := []string{}
			for _, msg := range annotatedMessages {
				names = append(names, string(msg.Desc.Name()))
			}

			slog.Debug("Will generate pgx decorators for following types", "types", names)
			plugin.GenerateDecoratorsForMessages(p, path, file, template, annotatedMessages...)
		}

		return nil
	})
}
