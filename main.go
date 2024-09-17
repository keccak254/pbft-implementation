package main

import (
    "os"
    "github.com/keccak254/pbft-implementation/network"
)

func main() {
    nodeID := os.Args[1]
    server := network.NewServer(nodeID)
    server.Start()
}