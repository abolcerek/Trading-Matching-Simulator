package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Connection, error) {
	const conn = "amqp://guest:guest@127.0.0.1:5672/"
	connection, err := amqp.Dial(conn)
	if err != nil {
		return &amqp.Connection{}, err
	}
	outboundExchange := "events"
	inboundExchange := "orders"
	outboundExchangeKind := "fanout"
	inboundExchangeKind := "direct"
	writerQueueName := "writer"
	wsQueueName := "ws"
	engineQueueName := "engine"
	ch, err := connection.Channel()
	if err != nil {
		return nil,  err
	}
	err = ch.ExchangeDeclare(outboundExchange, outboundExchangeKind, true, false, false, false, nil)
	if err != nil {
		return nil,  err
	}
	err = ch.ExchangeDeclare(inboundExchange, inboundExchangeKind, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	_, err = DeclareAndBind(ch, outboundExchange, writerQueueName, "")
	if err != nil {
		return nil, err
	}
	_, err = DeclareAndBind(ch, outboundExchange, wsQueueName, "")
	if err != nil {
		return nil, err
	}
	_, err = DeclareAndBind(ch, inboundExchange, engineQueueName, "")
	if err != nil {
		return nil, err
	}
	ch.Close()
	return connection, nil
}