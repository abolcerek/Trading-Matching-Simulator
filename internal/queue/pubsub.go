package queue

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange string, key string, val T) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ctx := context.Background()
	mandatory := false
	immediate := false
	pubStruct := amqp.Publishing{
		ContentType: "application/json",
		DeliveryMode: amqp.Persistent,
		Body: bytes,
	}
	err = ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, pubStruct)
	if err != nil {
		return err
	}
	return nil
}

func SubscribeJSON[T any](ch *amqp.Channel, exchange string, queueName string, key string, handler func(T) string) error {
	err := ch.Qos(10, 0, false)
	if err != nil {
		return err
	}
	messages, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for msg := range messages {
			var buffer T
			err = json.Unmarshal(msg.Body, &buffer)
			if err != nil {
				msg.Nack(false, false)
				continue
			}
			ackType := handler(buffer)
			switch ackType {
			case "ack":
				msg.Ack(false)
			case "nackRequeue":
				msg.Nack(false, true)
			case "nackDiscard":
				msg.Nack(false, false)
			default:
				fmt.Println("Unknown Ack type")
			}
		}
	}()
return nil
}

func DeclareAndBind(ch *amqp.Channel, exchange string, queueName string, key string) (amqp.Queue, error) {
	durable := true
	autoDelete := false
	exclusive := false
	queue, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		return amqp.Queue{}, err
	}
	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return amqp.Queue{}, err
	}
	return queue, nil
}
