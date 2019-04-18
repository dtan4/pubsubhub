package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

const (
	usage = "pubsubhub <GCP_PROJECT_ID>"
)

const (
	exitOK int = iota
	exitError
)

func run(args []string) int {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, usage)
		return exitError
	}
	projectID := args[1]

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	ss := []*pubsub.Subscription{}
	iter := client.Subscriptions(ctx)

	for {
		s, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitError
		}

		ss = append(ss, s)
	}

	for _, s := range ss {
		fmt.Println(s.ID())
	}

	return exitOK
}

func main() {
	run(os.Args)
}
