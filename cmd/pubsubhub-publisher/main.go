package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

const (
	usage = "pubsubhub-publisher <GCP_PROJECT_ID> <PUBSUB_TOPIC> <INTERVAL_MILLISECOND>"
)

const (
	defaultSeed = 123456789
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

func (c *pubSubClient) PublishContinuously(ctx context.Context, topic string, interval time.Duration, done <-chan int) error {
	t := c.client.Topic(topic)
	if t == nil {
		return errors.Errorf("topic %q not found", topic)
	}

	r := rand.New(rand.NewSource(defaultSeed))

	i := 0

	for {
		select {
		case <-done:
			return nil
		default:
		}

		body := fmt.Sprintf("%d - %d", i, r.Int())

		res := t.Publish(ctx, &pubsub.Message{
			Data: []byte(body),
		})

		id, err := res.Get(ctx)
		if err != nil {
			return errors.Wrap(err, "cannot publish message")
		}

		log.Printf("published id:%s body:%q", id, body)

		time.Sleep(interval)
		i++
	}
}

func run(args []string) int {
	if len(args) != 4 {
		log.Println(usage)
		return exitError
	}
	projectID, topic := args[1], args[2]

	i, err := strconv.Atoi(args[3])
	if err != nil {
		log.Println(err)
		return exitError
	}
	interval := time.Duration(i) * time.Millisecond

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err)
		return exitError
	}
	ps := NewPubSubClient(client)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	done := make(chan int)

	go func() {
		<-sigCh
		log.Println("Terminating...")
		close(done)
	}()

	if err := ps.PublishContinuously(ctx, topic, interval, done); err != nil {
		log.Println(err)
		return exitError
	}

	return exitOK
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	os.Exit(run(os.Args))
}
