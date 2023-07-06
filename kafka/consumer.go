package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type MessageProcessor func(ctx context.Context, msg []byte) error

type Subscription struct {
	Topic    string
	Listener MessageProcessor
}

type Consumer struct {
	brokers      []string
	groupID      string
	subscription []Subscription
}

func NewConsumer(brokers []string, groupID string, subscription []Subscription) *Consumer {
	return &Consumer{
		brokers:      brokers,
		groupID:      groupID,
		subscription: subscription,
	}
}

func (c *Consumer) RunSubscription(ctx context.Context, errCh chan error) {
	for _, s := range c.subscription {
		go func(subs Subscription) {
			errCh <- c.subscribe(ctx, subs)
		}(s)
	}
}

func (c *Consumer) subscribe(ctx context.Context, subs Subscription) error {
	admin := NewAdmin(c.brokers)
	defer admin.Close()
	if !admin.TopicExists(subs.Topic) {
		admin.CreateTopic(subs.Topic)
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(c.brokers...),
		kgo.ConsumerGroup(c.groupID),
		kgo.ConsumeTopics(subs.Topic),
	)
	if err != nil {
		return err
	}
	defer client.Close()
	log.Printf("subscribe to Topic: %s", subs.Topic)

	for {
		fetches := client.PollFetches(ctx)
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			fmt.Printf("Topic %s: ", subs.Topic)
			err := subs.Listener(ctx, record.Value)
			if err != nil {
				// anticipate if consume then error
				log.Println("error: ", err)
			}
		}
	}
}
