# Go Go Gadget

GoGoGadget provides a set of tools useful for cyber security testing packaged in a statically compiled binary. Have you ever been on an engagement and compromised a system only to find the system provides none of the tools you need? Even worse the system is running on ARM or PPC and you have no time to set up a cross-compiler to get something working... GoGoGadget is here to solve all your problems. Want to move files off the machine? `gogogoadget server` will start a web server. Need to pull something from your attack box? `gogogadget download` is what you want. Did you find a juicy target on the inside of the network? Check out `gogogadget pivot`.

# Current Gadgets
- download - a wget style download utility
- escalate - user escalation on Linux using dirtypipe
- pivot - recieve traffic on a port and forward to another host
- portscan - an nmap-like tcp scanner
- screenshot - take screenshots of any displays open on the device
- search - a grep style utility for searching file contents
- server - a web server allowing file downloading and uploading
- telnet - a telnet client

# Build

GoGoGadget uses the standard Go toolchain to target a bunch of different architectures and processor types. To see a full list of available operating systems and processor architectures type `go tool dist list`. Provide the operating system and processor architecture on the command line by typing `GOOS=<operating system> GOARCH=<architecture> go build`. 

# Usage

Once you have GoGoGadget on a target machine you can use any of the gadgets by typing `gogogadget <gadget>`. For specific usage information you can read help for each gadget using `gogogadget help <gadget>`.


# Reduced Size Binary

To reduce the size of the GoGoGadget binary we can strip the debugging symbols from the binary by appending the following build flag:

`-ldflags="-s -w"`

`-s` - Omit the symbol table and debug information

`-w` - Omit the DWARF symbol table

Additionally the binary can be compressed using the Ultimate Packer for Executables which compresses the binary into a self-decompressing binary.

`upx -9 gogogadget`

`-9` - Maximum compression

To ensure that upx supports the architecture that you're targeting type `upx -h` to see all supported architectures.