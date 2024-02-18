package main

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jakewright/home-automation/tools/libraries/imports"
)

type firehoseDataEvent struct {
	TypeName  string
	EventName string
}

type firehoseData struct {
	PackageName string
	Events      []*firehoseDataEvent
}

const firehoseTemplateText = `// Code generated by jrpc. DO NOT EDIT.

package {{ .PackageName }}

import (
	"encoding/json"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/firehose"
)

{{ range .Events }}
	// Publish publishes the event to the Firehose
	func (m *{{ .TypeName }}) Publish(ctx context.Context, p firehose.Publisher) error {
		if err := m.Validate(); err != nil {
			return err
		}

		return p.Publish(ctx, "{{ .EventName }}", m)
	}

	// {{ .TypeName }}Handler implements the necessary functions to be a Firehose handler
	type {{ .TypeName }}Handler func(*{{ .TypeName }}) firehose.Result

	// HandleEvent handles the Firehose event
	func (h {{ .TypeName }}Handler) HandleEvent(ctx context.Context, decode firehose.Decoder) firehose.Result {
		var body {{ .TypeName }}
		if err := decode(&body); err != nil {
			return firehose.Discard(oops.WithMessage(err, "failed to unmarshal payload"))
		}
		return h(&body)
	}
{{ end }}
`

type firehoseGenerator struct {
	baseGenerator
}

func (g *firehoseGenerator) Template() (*template.Template, error) {
	return template.New("firehose_template").Parse(firehoseTemplateText)
}

func (g *firehoseGenerator) PackageDir() string {
	return packageDirExternal
}

func (g *firehoseGenerator) Data(_ *imports.Manager) (interface{}, error) {
	var events []*firehoseDataEvent
	for _, m := range g.file.Messages {
		eventName, ok := m.Options["event_name"].(string)
		if !ok {
			continue
		}

		alias, parts := m.Lineage()
		if alias != "" {
			return nil, fmt.Errorf("unexpected alias in local message lineage %q", m.QualifiedName)
		}
		typeName := strings.Join(parts, "_") // Replicate what typesGenerator{} does

		events = append(events, &firehoseDataEvent{
			TypeName:  typeName,
			EventName: eventName,
		})
	}

	if len(events) == 0 {
		return nil, nil
	}

	return &firehoseData{
		PackageName: externalPackageName(g.options),
		Events:      events,
	}, nil
}

func (g *firehoseGenerator) Filename() string {
	return "firehose.go"
}
