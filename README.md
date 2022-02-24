# Go Go Gadget

This program provides a set of tools useful for cyber security testing in a statically compiled binary. Have you ever been on an engagement and compromised a system only to find the system provides none of the tools you need? Even worse the system is running on ARM or PPC and you have no time to set up a cross-compiler to get something working... GoGoGadget is here to solve all your problems. Want to move files off the machine? `gogogoadget server` will start a web server. Need to pull something from your attack box? `gogogadget download` is what you want. Did you find a juicy target on the inside of the network? Check out `gogogadget pivot`. The `compile.sh` script assists with cross compiling the binary for a target OS and architecture. To see a list of available OSes and ARCHitectures run `compile.sh help`.

# Current Tools

GoGoGadget will be consistently growing as I identify new tools I wish I had available on compromised boxes. Please feel free to contribute!

## Server

The server provides a simple file server. You can specify a folder to serve and a port to listen on.

## Download

Simple file downloader - point it to a URL and it will download it.

## Pivot

Take traffic in on one port and send it out to another. 
