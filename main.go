package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
)

func main() {
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/8007"),
		libp2p.EnableHolePunching(),
		libp2p.EnableNATService(),
		libp2p.EnableRelay(),
	)

	host.SetStreamHandler("/chat/1.0.0", handleStream)

	for _, addr := range host.Addrs() {
		fmt.Printf("Listening on %s/p2p/%s\n", addr, host.ID())

	}

	fmt.Println("Hello World, i am main host and my multiaddr is ", fmt.Sprintf("%s/p2p/%s", host.Addrs()[0], host.ID().String()))
	fmt.Println("Try to run other code with this multtiaddr and dest flag")
	if err != nil {
		panic(err)
	}
	fmt.Println("Waiting for connection")
	fmt.Println("Source host peers: ", host.Network().Peers())

	select {} //wait forever

}

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
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
