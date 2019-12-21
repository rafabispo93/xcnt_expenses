package event

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	kafka "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type kafkaMock struct {
	FetchMessageStub   func(ctx context.Context) (kafka.Message, error)
	CommitMessagesStub func(ctx context.Context, msgs ...kafka.Message) error
	WriteMessagesStub  func(ctx context.Context, msgs ...kafka.Message) error
}

func (m *kafkaMock) FetchMessage(ctx context.Context) (kafka.Message, error) {
	if m.FetchMessageStub != nil {
		return m.FetchMessageStub(ctx)
	}
	return kafka.Message{}, errors.New("not mocked")
}

func (m *kafkaMock) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	if m.CommitMessagesStub != nil {
		return m.CommitMessagesStub(ctx, msgs...)
	}
	return errors.New("not mocked")
}

func (m *kafkaMock) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	if m.WriteMessagesStub != nil {
		return m.WriteMessagesStub(ctx, msgs...)
	}
	return errors.New("not mocked")
}

func TestKafkaReadMessage(t *testing.T) {
	t.Parallel()

	expected := []byte(gofakeit.Paragraph(1, 2, 3, "."))

	var acked bool
	m := &kafkaMock{
		FetchMessageStub: func(ctx context.Context) (kafka.Message, error) {
			return kafka.Message{
				Value: expected,
			}, nil
		},
		CommitMessagesStub: func(ctx context.Context, msgs ...kafka.Message) error {
			acked = true
			return nil
		},
	}

	d := kafkaReader(m)
	var done bool
	require.NoError(t, d.ReadMessage(context.TODO(), func(_ context.Context, actual []byte) error {
		assert.Equal(t, expected, actual)
		done = true
		return nil
	}))
	assert.True(t, acked)
	assert.True(t, done)
}

func TestKafkaWriteMessage(t *testing.T) {
	t.Parallel()

	expected := []byte(gofakeit.Paragraph(1, 2, 3, "."))

	var sent bool
	m := &kafkaMock{
		WriteMessagesStub: func(ctx context.Context, msgs ...kafka.Message) error {
			assert.Equal(t, expected, msgs[0].Value)
			sent = true
			return nil
		},
	}

	d := kafkaWriter(m)
	require.NoError(t, d.WriteMessage(context.TODO(), expected))
	assert.True(t, sent)
}
