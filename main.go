package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
)

const (
	usage = "pubsubhub <GCP_PROJECT_ID>"
)

const (
	exitOK int = iota
	exitError
)

func subscribe(ctx context.Context, s *pubsub.Subscription) error {
	err := s.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Printf("[subscription:%s message:%s] %s\n", s.ID(), m.ID, string(m.Data))
		m.Ack()
		fmt.Printf("[subscription:%s message:%s] ACK\n", s.ID(), m.ID)
	})
	if err != nil {
		return errors.Wrapf(err, "subscribe error: %s", s.ID())
	}

	return nil
}

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

		c, err := s.Config(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitError
		}

		// skip push subscription
		if c.PushConfig.Endpoint != "" {
			continue
		}

		ss = append(ss, s)
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, s := range ss {
		s := s
		g.Go(func() error {
			fmt.Printf("Start subscribing %s...\n", s.ID())

			return subscribe(ctx, s)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	return exitOK
}

func main() {
	run(os.Args)
}
