package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "github.com/keccak254/pbft-implementation/network"
)

func main() {
    nodeID := flag.String("id", "", "Node ID")
    flag.Parse()

    if *nodeID == "" {
        log.Fatal("Please provide a node ID using the -id flag")
    }

    node := network.NewNode(*nodeID)

    http.HandleFunc("/req", node.HandleRequest)
    http.HandleFunc("/preprepare", node.HandlePrePrepare)
    http.HandleFunc("/prepare", node.HandlePrepare)
    http.HandleFunc("/commit", node.HandleCommit)


    port := "1111"
    switch *nodeID {
    case "MS":
        port = "1112"
    case "Google":
        port = "1113"
    case "IBM":
        port = "1114"
    }

    fmt.Printf("Node %s listening on port %s\n", *nodeID, port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}