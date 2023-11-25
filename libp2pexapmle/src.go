package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	host, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	for _, addr := range host.Addrs() {
		fmt.Printf("Listening on %s/p2p/%s\n", addr, host.ID())

	}

	//Get ID from command line flags
	dest := flag.String("dest", "", "target peer to dial")
	flag.Parse()

	destMultiaddr, err := multiaddr.NewMultiaddr(*dest)
	if err != nil {
		panic(err)
	}
	pi, err := peer.AddrInfoFromP2pAddr(destMultiaddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Hello World, my hosts ID is ", host.ID())
	err = host.Connect(context.Background(), *pi)
	if err != nil {
		panic(err)
	}

	host.SetStreamHandler("/chat/1.0.0", handleStream)

	host.Peerstore().AddAddr(pi.ID, pi.Addrs[0], peerstore.PermanentAddrTTL)
	s, err := host.NewStream(context.Background(), pi.ID, "/chat/1.0.0")
	if err != nil {
		panic(err)
	}
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go writeData(rw)
	go readData(rw)

	select {} //wait forever

}

func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}
