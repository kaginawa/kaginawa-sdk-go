package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kaginawa/kaginawa-sdk-go"
)

func main() {
	var scanner = bufio.NewScanner(os.Stdin)
	nextLine := func() string { scanner.Scan(); return scanner.Text() }
	fatal := func(err error) { println(err.Error()); os.Exit(1) }

	// Build client
	fmt.Print("endpoint > ")
	endpoint := nextLine()
	fmt.Print("api key > ")
	apiKey := nextLine()
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}
	client, err := kaginawa.NewClient(endpoint, apiKey)
	if err != nil {
		fatal(err)
	}

	// Collect target information
	fmt.Print("target id > ")
	id := nextLine()
	report, err := client.FindNode(context.Background(), id)
	if err != nil {
		fatal(err)
	}
	fmt.Println(report.ID, report.Hostname)

	// Login prompt
	fmt.Print("user > ")
	user := nextLine()
	fmt.Print("password > ")
	password := nextLine()

	// Command prompt
	for {
		fmt.Print("command (type \"exit\" to exit) > ")
		command := nextLine()
		if command == "exit" || command == "quit" {
			break
		}
		result, err := client.Command(context.Background(), id, command, user, "", password, 0)
		if err != nil {
			println(err.Error())
			if strings.Contains(err.Error(), "authenticate") {
				os.Exit(1)
			}
		} else {
			fmt.Println(result)
		}
	}
}
