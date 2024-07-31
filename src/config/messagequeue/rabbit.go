package messagequeue

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.gin.order/src/internal/pojo"
	"log"
)

var (
	RabbitConn    *amqp.Connection
	RabbitChannel *amqp.Channel
	err           error
)

func Mqinit(mq *pojo.RabbitmqConf) {
	log.Println(mq.Url)
	RabbitConn, err = amqp.Dial(mq.Url)
	if err != nil {
		log.Println("连接RabbitMQ失败", err)
		panic(err)
	}
	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		log.Println("获取RabbitMQ channel失败")
		panic(err)
	}
	fmt.Println("rabbitmq初始化连接成功")
}

type Consumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}
type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewProducer(queueName string) (*Producer, error) {
	declare, err := RabbitChannel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &Producer{
		conn:    RabbitConn,
		channel: RabbitChannel,
		queue:   declare,
	}, nil
}

func NewConsumer(queueName string) (*Consumer, error) {
	_, err := RabbitChannel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:      RabbitConn,
		channel:   RabbitChannel,
		queueName: queueName,
	}, nil
}

func (p *Producer) Publish(message string) error {
	err := p.channel.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return err
}

func (p *Producer) Close() {
	p.channel.Close()
	p.conn.Close()
}

func (c *Consumer) Consume(handler func(message string)) error {
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Println(err, "++++++++++++++++++++++++++++++")
		return err
	}

	go func() {
		for d := range msgs {
			handler(string(d.Body))
		}
	}()

	return nil
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}
