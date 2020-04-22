package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kaginawa/kaginawa-sdk-go/kaginawa"
)

func main() {
	endpoint := flag.String("e", "", "endpoint (https://...)")
	key := flag.String("k", "", "api key")
	cid := flag.String("c", "", "custom id")
	flag.Parse()

	// Prepare the API client
	client, err := kaginawa.NewClient(*endpoint, *key)
	if err != nil {
		println(err)
		os.Exit(2)
	}
	if len(*cid) == 0 {
		println("most specify a custom id")
		os.Exit(2)
	}

	// Gathering reports
	nodes, err := client.ListNodesByCustomID(context.Background(), *cid)
	if err != nil {
		println(err)
		os.Exit(1)
	}
	fmt.Printf("found %d node(s)\n", len(nodes))
	sshServerSet := map[string]struct{}{}
	for _, node := range nodes {
		if node.SSHRemotePort > 0 {
			fmt.Printf("🔵 %s %s (%v) %s:%d\n",
				node.ID, node.CustomID, node.Timestamp(), node.SSHServerHost, node.SSHRemotePort)
			sshServerSet[node.SSHServerHost] = struct{}{}
		} else {
			fmt.Printf("⚪ %s %s (%v) (ssh forwarding disabled or not connected)\n",
				node.ID, node.CustomID, node.Timestamp())
		}
	}

	// Gathering ssh server information
	if len(sshServerSet) > 0 {
		fmt.Printf("found %d ssh server(s)\n", len(sshServerSet))
		for sshServer := range sshServerSet {
			sv, err := client.FindSSHServerByHostname(context.Background(), sshServer)
			if err != nil {
				println(err)
				os.Exit(1)
			}
			fmt.Printf("⚫ %s@%s:%d key=%v pw=%v\n", sv.User, sv.Host, sv.Port, len(sv.Key) > 0, len(sv.Password) > 0)
		}
	}
}
