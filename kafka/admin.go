package kafka

import (
	"context"
	"log"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Admin struct {
	client *kadm.Client
}

func NewAdmin(brokers []string) *Admin {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		panic(err)
	}

	admin := kadm.NewClient(client)
	return &Admin{client: admin}
}

func (a *Admin) TopicExists(topic string) bool {
	ctx := context.Background()
	topicsMetadata, err := a.client.ListTopics(ctx)
	if err != nil {
		panic(err)
	}

	if topicsMetadata.Has(topic) {
		return true
	}

	return false
}

func (a *Admin) CreateTopic(topic string) {
	ctx := context.Background()
	resp, err := a.client.CreateTopics(ctx, 1, 1, nil, topic)
	if err != nil {
		panic(err)
	}

	for _, ctr := range resp {
		if ctr.Err != nil {
			log.Fatalf("Unable to create topic '%s': %s", ctr.Topic, ctr.Err)
		} else {
			log.Printf("Created topic '%s'\n", ctr.Topic)
		}
	}
}

func (a *Admin) Close() {
	a.client.Close()
}
