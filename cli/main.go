package main

import (
	"os"

	"github.com/dpb587/bosh-compiled-releases/cli/cmd"
	"github.com/jessevdk/go-flags"
)

func main() {
	var parser = flags.NewParser(&struct{}{}, flags.Default)

	parser.AddCommand("serve", "Start an HTTP server to answer compiled release lookups", "", &cmd.Serve{})
	parser.AddCommand("rewrite-manifest", "Rewrite a manifest to reference compiled releases", "", &cmd.RewriteManifest{})
	parser.AddCommand("file-add-compiled-release", "Add a compiled release reference to a file repository", "", &cmd.FileAddCompiledRelease{})

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
