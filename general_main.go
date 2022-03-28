// +build !linux,!windows

package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/seandheath/gogogadget/internal/download"
	"github.com/seandheath/gogogadget/internal/gadget"
	"github.com/seandheath/gogogadget/internal/pivot"
	"github.com/seandheath/gogogadget/internal/server"
	"github.com/seandheath/gogogadget/internal/search"
)

// Add your gadget here
var gadgets = [...]*gadget.Gadget{
	&download.Gadget,
	&pivot.Gadget,
	&server.Gadget,
	&search.Gadget,
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	//subcommands.Register(subcommands.FlagsCommand(), "")
	//subcommands.Register(subcommands.CommandsCommand(), "")

	for _, g := range gadgets {
		subcommands.Register(g, "")
	}

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))

}
