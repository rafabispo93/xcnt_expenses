package event

import (
	"context"

	"github.com/pkg/errors"
	kafka "github.com/segmentio/kafka-go"
)

type kafkaReceiver interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
}

func kafkaReader(r kafkaReceiver) MessageReader {
	return MessageReaderFunc(func(ctx context.Context, cb MessageCallback) error {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to fetch message from Kafka")
		}
		if err := cb(ctx, msg.Value); err != nil {
			return err
		}

		return errors.Wrap(r.CommitMessages(ctx, msg), "failed to commit message in Kafka")
	})
}

// NewKafkaReader initializes a kafka.Reader for use in a Consumer.
func NewKafkaReader(host, group string, topic Topic) MessageReader {
	return kafkaReader(kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{host},
		GroupID: group,
		Topic:   topic.String(),
	}))
}

type kafkaSender interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

func kafkaWriter(s kafkaSender) MessageWriter {
	return MessageWriterFunc(func(ctx context.Context, msgs ...[]byte) error {
		kmm := make([]kafka.Message, len(msgs))
		for i, msg := range msgs {
			kmm[i] = kafka.Message{Value: msg}
		}
		return errors.Wrap(s.WriteMessages(ctx, kmm...), "failed to write message to Kafka")
	})
}

// NewKafkaWriter initializes a kafka.Writer for use in a Producer.
func NewKafkaWriter(host, name string, topic Topic) MessageWriter {
	return kafkaWriter(kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{host},
		Topic:   topic.String(),
		Dialer: &kafka.Dialer{
			ClientID: name,
		},
	}))
}
