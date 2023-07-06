package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Topics struct {
	CryptoPublished string
}

type Publisher struct {
	client *kgo.Client
	topic  Topics
}

func NewPublisher(brokers []string, topics Topics) *Publisher {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		panic(err)
	}

	return &Publisher{client: client, topic: topics}
}

func (p *Publisher) PostCrypto(ctx context.Context, c CryptoModel) error {
	b, err := c.Serialize()
	if err != nil {
		return err
	}

	record := &kgo.Record{Topic: p.topic.CryptoPublished, Value: b}
	if err := p.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return err
	}

	return nil
}

func (p *Publisher) Close() {
	p.client.Close()
}
