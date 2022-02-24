package gadget

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type Gadget struct {
	GadgetName     string
	GadgetSynopsis string
	GadgetUsage    string
	Run            func()
	InitFlags      func(f *flag.FlagSet)
}

func (g *Gadget) Name() string             { return g.GadgetName }
func (g *Gadget) Synopsis() string         { return g.GadgetSynopsis }
func (g *Gadget) Usage() string            { return g.GadgetUsage }
func (g *Gadget) SetFlags(f *flag.FlagSet) { g.InitFlags(f) }
func (g *Gadget) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	g.Run()
	// TODO: Check if subcommand succeeded or not
	return subcommands.ExitSuccess
}

//func New(name string, synopsis string, usage string, setFlags func(f *flag.FlagSet), run func()) Gadget {
//return Gadget{
//name:     name,
//synopsis: synopsis,
//usage:    usage,
//Run:      run,
//}
//}
