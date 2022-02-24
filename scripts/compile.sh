#!/bin/bash

help() {
    echo "Select what architecture you want to build for (or leave empty)"
    echo "Available OS/Arch:"
    go tool dist list
    echo "Usage: ./compile.sh [OS ARCH]"
}

trap 'help' ERR

if [[ ! -e "build" ]]; then
    mkdir build
fi

ARGC=$#
GOOS=""
GOARCH=""
FNAME="gogogadget"
if [ $ARGC -gt "0" ]; then
    if [ $ARGC -eq 2 ]; then
        GOOS=$1
        GOARCH=$2
        FNAME="$FNAME$GOOS$GOARCH"
        echo "Building GoGoGadget for $GOOS/$GOARCH"
    else 
        help
    fi
fi

env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -a -o build/$FNAME
