package plugin

import (
	"encoding/json"
	"path/filepath"

	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	gendoc "github.com/pseudomuto/protoc-gen-doc"
	"github.com/pseudomuto/protokit"
	"google.golang.org/protobuf/proto"
)

// SupportedFeatures describes a flag setting for supported features.
var SupportedFeatures = uint64(plugin_go.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

// Plugin describes a protoc code generate plugin. It's an implementation of Plugin from github.com/pseudomuto/protokit
type Plugin struct{}

// Generate compiles the documentation and generates the CodeGeneratorResponse to send back to protoc. It does this
// by rendering a template based on the options parsed from the CodeGeneratorRequest.
func (p *Plugin) Generate(r *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {

	result := protokit.ParseCodeGenRequest(r)

	resp := new(plugin_go.CodeGeneratorResponse)
	fdsGroup := groupProtosByDirectory(result, false)

	for dir, fds := range fdsGroup {
		tmpl := gendoc.NewTemplate(fds)
		for _, template := range IsolateFileAsTemplate(tmpl) {

			data, err := json.MarshalIndent(template, "", "  ")
			if err != nil {
				return nil, err
			}
			resp.File = append(resp.File, &plugin_go.CodeGeneratorResponse_File{
				Name:    proto.String(filepath.Join(dir, template.Files[0].Name) + ".json"),
				Content: proto.String(string(data)),
			})

		}

	}

	resp.SupportedFeatures = proto.Uint64(SupportedFeatures)

	return resp, nil
}

func IsolateFileAsTemplate(file *gendoc.Template) []*gendoc.Template {
	tmps := make([]*gendoc.Template, 0)
	for _, v := range file.Files {
		tmps = append(tmps, &gendoc.Template{
			Files:   []*gendoc.File{v},
			Scalars: file.Scalars,
		})
	}

	return tmps
}

func groupProtosByDirectory(fds []*protokit.FileDescriptor, sourceRelative bool) map[string][]*protokit.FileDescriptor {
	fdsGroup := make(map[string][]*protokit.FileDescriptor)

	for _, fd := range fds {
		dir := ""
		if sourceRelative {
			dir, _ = filepath.Split(fd.GetName())
		}
		if dir == "" {
			dir = "./"
		}
		fdsGroup[dir] = append(fdsGroup[dir], fd)
	}
	return fdsGroup
}
