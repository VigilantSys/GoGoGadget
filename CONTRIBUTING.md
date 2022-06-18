This project uses the `cobra-cli` tool to add new gadgets. To use it you can follow these steps:

1. Install cobra-cli

`go install github.com/spf13/cobra-cli@latest`

2. Add a gadget using the following command

`cobra-cli --config cobra.yaml add <gadget name>`

A new file will be added to the `cmd/` directory with `<gadget name>.go`. Add working code to the `Run` function, add flags in the `init` function. See the other gadgets for examples.
