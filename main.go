package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

var (
	rabbitDSN   = kingpin.Flag("rabbit-dsn", "RabbitMQ DSN.").String()
	queuesCount = kingpin.Flag("queues-count", "The count of queues.").Default("20000").Int()
	exchange    = kingpin.Flag("exchange", "RabbitMQ Exchange.").Default("test").String()
)

func main() {
	kingpin.Parse()

	conn, err := amqp.Dial(*rabbitDSN)

	if err != nil {
		log.Fatalf("cannot (re)dial: %v: %q", err, rabbitDSN)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("cannot create channel: %v", err)
	}

	err = channel.ExchangeDeclare(
		*exchange,      // name
		"headers", // type
		true,    // durable
		false, // auto-deleted
		false,   // internal
		false,   // noWait
		nil,       // arguments
	)

	if err != nil {
		log.Fatalf("cannot declare fanout exchange: %v", err)
	}

	var i = 0;
	for ;i < *queuesCount; i++ {
		queue, err := channel.QueueDeclare(
			"",  // name of the queue
			false, // durable
			true,  // delete when unused
			true,  // exclusive
			false, // noWait
			nil,   // arguments
		)
		if err != nil {
			log.Fatalf("cannot declare queue: %v", err)
		}

		headers := make(amqp.Table)
		headers["x-match"] = "any"
		err = channel.QueueBind(
			queue.Name,    // name of the queue
			"",            // bindingKey
			*exchange, // sourceExchange
			false,         // noWait
			headers, // arguments
		)
		if err != nil {
			log.Fatalf("could not bind queue to exchange: %v", err)
		}
	}
	fmt.Println("done", i)

	if err = conn.Close(); err != nil {
		log.Fatalf("cannot close amqp connection: %v", err)
	}
}
