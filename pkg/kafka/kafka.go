package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type KafkaConfig struct {
	Brokers []string
}

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(config KafkaConfig, topic string) *Producer {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	if err != nil {
		log.Fatalf("failed to create sarama producer: %v", err)
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}
}

func (p *Producer) Produce(ctx context.Context, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message: %v", err)
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

type Consumer struct {
	group sarama.ConsumerGroup
	topic string
}

func NewConsumer(config KafkaConfig, topic, groupID string) *Consumer {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	group, err := sarama.NewConsumerGroup(config.Brokers, groupID, saramaConfig)
	if err != nil {
		log.Fatalf("failed to create sarama consumer group: %v", err)
	}

	return &Consumer{
		group: group,
		topic: topic,
	}
}

func (c *Consumer) Consume(ctx context.Context, handler sarama.ConsumerGroupHandler) {
	for {
		if err := c.group.Consume(ctx, []string{c.topic}, handler); err != nil {
			log.Printf("Error from consumer group: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func (c *Consumer) Close() error {
	return c.group.Close()
}
