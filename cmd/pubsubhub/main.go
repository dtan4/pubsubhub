package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
)

const (
	usage = "pubsubhub <GCP_PROJECT_ID>"
)

const (
	defaultTimeout = 1 * time.Minute
)

const (
	exitOK int = iota
	exitError
)

type pubSubClient struct {
	client *pubsub.Client
}

func NewPubSubClient(client *pubsub.Client) *pubSubClient {
	return &pubSubClient{
		client: client,
	}
}

func (c *pubSubClient) ListSubscriptions(ctx context.Context) ([]*pubsub.Subscription, error) {
	ss := []*pubsub.Subscription{}
	iter := c.client.Subscriptions(ctx)

	for {
		s, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return []*pubsub.Subscription{}, errors.Wrap(err, "cannot iterate subscriptions")
		}

		c, err := s.Config(ctx)
		if err != nil {
			return []*pubsub.Subscription{}, errors.Wrapf(err, "cannot get subscription config of %s", s.ID())
		}

		// skip push subscription
		if c.PushConfig.Endpoint != "" {
			continue
		}

		ss = append(ss, s)
	}

	return ss, nil
}

func subscribe(ctx context.Context, s *pubsub.Subscription) error {
	err := s.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		log.Printf("[subscription:%s message:%s] body:%q", s.ID(), m.ID, string(m.Data))
		m.Ack()
		log.Printf("[subscription:%s message:%s] ACK", s.ID(), m.ID)
	})
	if err != nil {
		return errors.Wrapf(err, "subscribe error: %s", s.ID())
	}

	log.Printf("[subscription:%s] Stopping...", s.ID())

	return nil
}

func run(args []string) int {
	if len(args) != 2 {
		log.Println(usage)
		return exitError
	}
	projectID := args[1]

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err)
		return exitError
	}
	ps := NewPubSubClient(client)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	for {
		select {
		case <-sigCh:
			log.Println("Terminating...")
			return exitOK
		case <-ctx.Done():
			return exitOK
		default:
			log.Println("Retrieving subscriptions...")

			ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
			defer cancel()

			ss, err := ps.ListSubscriptions(ctx)
			if err != nil {
				log.Println(err)
				return exitError
			}

			g, ctx := errgroup.WithContext(ctx)

			for _, s := range ss {
				s := s
				g.Go(func() error {
					log.Printf("Start subscribing %s...", s.ID())

					return subscribe(ctx, s)
				})
			}

			if err := g.Wait(); err != nil {
				log.Println(err)
				return exitError
			}
		}
	}
}

func main() {
	run(os.Args)
}
