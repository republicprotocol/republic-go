package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/republicprotocol/republic-go/grpc"
	"github.com/urfave/cli"
	netGrpc "google.golang.org/grpc"
)

func main() {
	// Create new cli application
	app := cli.NewApp()

	// Define sub-commands
	app.Commands = []cli.Command{
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "get status of the node with given address",
			Action: func(c *cli.Context) error {
				return GetStatus(c.Args())
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func GetStatus(args []string) error {
	// Check arguments length
	if len(args) != 1 {
		return fmt.Errorf("please provide target network address (e.g. 0.0.0.0:8080)")
	}
	address := args[0]

	// Dial to the target node
	conn, err := netGrpc.Dial(address, netGrpc.WithInsecure())
	if err != nil {
		return err
	}
	c := grpc.NewStatusServiceClient(conn)

	// Call status rpc
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rsp, err := c.Status(ctx, &grpc.StatusRequest{})
	if err != nil {
		return err
	}

	fmt.Println(rsp.Address, rsp.Bootstrapped, rsp.Peers)

	return nil
}
