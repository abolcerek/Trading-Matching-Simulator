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
	exchange := "events"
	exchangeKind := "fanout"
	writerQueueName := "writer"
	wsQueueName := "ws"
	ch, err := connection.Channel()
	if err != nil {
		return nil,  err
	}
	err = ch.ExchangeDeclare(exchange, exchangeKind, true, false, false, false, nil)
	if err != nil {
		return nil,  err
	}
	_, err = DeclareAndBind(ch, exchange, writerQueueName, "")
	if err != nil {
		return nil, err
	}
	_, err = DeclareAndBind(ch, exchange, wsQueueName, "")
	if err != nil {
		return nil, err
	}
	ch.Close()
	return connection, nil
}