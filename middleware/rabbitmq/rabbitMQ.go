package rabbitmq

import (
	"fmt"
	"log"
	"tiktok_demo/config"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn  *amqp.Connection
	mqurl string
}

var Rmq *RabbitMQ

// InitRabbitMQ 初始化RabbitMQ的连接和通道。
func InitRabbitMQ() {
	Config := config.GetConfig()
	log.Println(Config)
	MQURL := "amqp://" + Config.RabbitMQ.DefaultUser + ":" + Config.RabbitMQ.DefaultPass + "@" + Config.RabbitMQ.Host + ":" + Config.RabbitMQ.Port + "/"
	log.Println(MQURL)
	Rmq = &RabbitMQ{
		mqurl: MQURL,
	}
	dial, err := amqp.Dial(Rmq.mqurl)
	Rmq.failOnErr(err, "创建连接失败")
	Rmq.conn = dial

}

// 连接出错时，输出错误信息。
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s\n", err, message)
		panic(fmt.Sprintf("%s:%s\n", err, message))
	}
}

// 关闭mq通道和mq的连接。
func (r *RabbitMQ) destroy() {
	err := r.conn.Close()
	if err != nil {
		return
	}
}
