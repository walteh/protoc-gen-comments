// protoc-gen-doc is used to generate documentation from comments in your proto files.
//
// It is a protoc plugin, and can be invoked by passing `--doc_out` and `--doc_opt` arguments to protoc.
//
// Example: generate HTML documentation
//
//	protoc --doc_out=. --doc_opt=html,index.html protos/*.proto
//
// Example: use a custom template
//
//	protoc --doc_out=. --doc_opt=custom.tmpl,docs.txt protos/*.proto
//
// For more details, check out the README at https://github.com/pseudomuto/protoc-gen-doc
package main

import (
	"flag"
	"runtime/debug"

	"github.com/pseudomuto/protokit"

	"log"

	"github.com/walteh/protoc-gen-comments/plugin"
)

var versionFlag = false

func init() {
	flag.BoolVar(&versionFlag, "version", false, "print the version and exit")
}

func main() {

	if versionFlag {
		v, ok := debug.ReadBuildInfo()
		if !ok {
			log.Printf("unknown\n")
			return
		}
		log.Printf("%s %s\n", v.Main.Path, v.Main.Version)
		return
	}

	if err := protokit.RunPlugin(new(plugin.Plugin)); err != nil {
		log.Fatal(err)
	}
}

// HandleFlags checks if there's a match and returns true if it was "handled"
