package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/filecoin-project/mir"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func main() {
	m := map[t.ModuleID]modules.Module{
		"mempool": simplemempool.NewModule(),
	}

	// create a Mir node
	node, err := mir.NewNode("client", &mir.NodeConfig{Logger: logging.NilLogger}, m, nil, nil)
	if err != nil {
		panic(fmt.Sprintf("error creating a Mir node: %v", err))
	}

	// wait for user input in a separate goroutine
	go readDataAndInjectToNode(node)

	// run the node
	err = node.Run(context.Background())
	if err != nil {
		panic(fmt.Sprintf("error running node: %v", err))
	}
}

func readDataAndInjectToNode(node *mir.Node) {
	// Read the user input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Type in the data and press Enter: ")

	scanner.Scan()
	if scanner.Err() != nil {
		panic(fmt.Sprintf("error reading user data: %v", scanner.Err()))
	}

	text := scanner.Text()

	// Inject the request
	request := &requestpb.Request{Data: []byte(text)}
	event := events.NewClientRequests("storage-client", []*requestpb.Request{request})
	err := node.InjectEvents(context.Background(), events.ListOf(event))
	if err != nil {
		panic(fmt.Sprintf("error injecting user request: %v", err))
	}
}
