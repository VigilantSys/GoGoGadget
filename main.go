package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/vigilantsys/gogogadget/internal/download"
	"github.com/vigilantsys/gogogadget/internal/escalate"
	"github.com/vigilantsys/gogogadget/internal/gadget"
	"github.com/vigilantsys/gogogadget/internal/pivot"
	"github.com/vigilantsys/gogogadget/internal/search"
	"github.com/vigilantsys/gogogadget/internal/server"
)

// Add your gadget here
var gadgets = [...]*gadget.Gadget{
	&download.Gadget,
	&pivot.Gadget,
	&server.Gadget,
	&escalate.Gadget,
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
