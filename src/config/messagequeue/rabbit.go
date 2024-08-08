package messagequeue

import (
	"github.com/streadway/amqp"
	"go.gin.order/src/config"
	"log"
)

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queue     amqp.Queue
	queueName string
}

func NewRabbitMQ() (*RabbitMQ, error) {
	RabbitConn, err := amqp.Dial(config.MQURL)
	if err != nil {
		log.Println("连接RabbitMQ失败", err)
		return nil, err
		//panic(err)
	}
	RabbitChannel, err := RabbitConn.Channel()
	if err != nil {
		log.Println("获取RabbitMQ channel失败")
		return nil, err
		//panic(err)
	}
	mq := &RabbitMQ{
		conn:    RabbitConn,
		channel: RabbitChannel,
	}
	mq.DeclareExchange("fanout_exchange", "fanout")
	mq.DeclareExchange("direct_exchange", "direct")
	mq.DeclareExchange("topic_exchange", "topic")
	return mq, nil
}

// 声明交换机
func (r *RabbitMQ) DeclareExchange(name, kind string) error {
	return r.channel.ExchangeDeclare(
		name,  // name
		kind,  // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// 声明队列
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// 将队列绑定到交换机上，将信息分布到各个队列
func (r *RabbitMQ) BindQueue(queueName, routingKey, exchangeName string) error {
	return r.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
}

// 发布消息到交换机
func (r *RabbitMQ) Publish(exchange, routingKey string, body []byte) error {
	return r.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

// 从队列中消费信息
func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	deliveryChan, err := r.channel.Consume(
		queueName,
		"",
		true, // Auto-acknowledge set to false for manual acknowledgment
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return deliveryChan, nil
}

func (r *RabbitMQ) Acknowledge(deliveryTag uint64) error {
	return r.channel.Ack(deliveryTag, false)
}

func (r *RabbitMQ) Cancel(consumerTag string) error {
	return r.channel.Cancel(consumerTag, false)
}

func (r *RabbitMQ) Close() {
	if err := r.channel.Close(); err != nil {
		log.Printf("Failed to close channel: %v", err)
	}
	if err := r.conn.Close(); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}
