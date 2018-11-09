package manageMq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sctek.com/typhoon/th-platform-gateway/common"
	"sctek.com/typhoon/th-platform-gateway/sms"
	"strings"
)

var GlobalMq *MessageManageMq

type MessageManageMq struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	topics  string //topic fanout direct
	done    chan error
}

func newMessageManageMa() *MessageManageMq {
	mesMq := &MessageManageMq{
		conn:    nil,
		channel: nil,
		tag:     "manage_message_consumer",
		topics:  "",
		done:    make(chan error),
	}
	return mesMq
}

func InitMq()  {
	GlobalMq = newMessageManageMa()
	err:=GlobalMq.Connect(common.Config.ManageMq.Uri)
	if err!=nil{
		panic(err)
	}
	err=GlobalMq.ExchangeDeclare(common.Config.ManageMq.Exchange,common.Config.ManageMq.ExchangeType)
	if err!=nil{
		panic(err)
	}
	err = GlobalMq.QueueDeclare(common.Config.ManageMq.QueueName,common.Config.ManageMq.Key,common.Config.ManageMq.Exchange)
	if err!=nil{
		panic(err)
	}
}

func (m *MessageManageMq) Connect(uri string) (err error) {
	common.Log.Errorf("dialing %q", uri)
	m.conn, err = amqp.Dial(uri)
	if err != nil {
		fmt.Println("链接失败！！", uri)
		return err
	}
	common.Log.Errorf("got Connection, getting Channel")
	m.channel, err = m.conn.Channel()
	if err != nil {
		err = fmt.Errorf("Channel: %s", err)
		common.Log.Errorln(err)
		return err
	}
	go func() {
		fmt.Printf("closing: %s", <-m.conn.NotifyClose(make(chan *amqp.Error)))
	}()
	return nil
}

func (m *MessageManageMq) ExchangeDeclare(exchange, exchangeType string) error {
	common.Log.Errorf("got Channel, declaring Exchange (%q)\r\n", exchange)
	if err := m.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		err = fmt.Errorf("Exchange Declare: %s", err)
		common.Log.Errorln(err)
		return err
	}
	return nil
}

func (m *MessageManageMq) QueueDeclare(qName, key, exchange string) error {
	common.Log.Errorf("declared Exchange, declaring Queue %q", qName)
	queue, err := m.channel.QueueDeclare(
		qName, // name of the queue
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		err = fmt.Errorf("Queue Declare: %s", err)
		common.Log.Errorln(err)
		return err
	}
	common.Log.Errorf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = m.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		err = fmt.Errorf("Queue Bind: %s", err)
		common.Log.Errorln(err)
		return err
	}
	return nil
}

// 发布消息
func (m *MessageManageMq) Publish(topic, msg string) (err error) {
	common.Log.Errorf("publish exchangeType=%q,msg=%q\r\n", topic, msg)
	if m.topics == "" || !strings.Contains(m.topics, topic) {
		err = m.channel.ExchangeDeclare(topic, "fanout", true, false, false, true, nil)
		if err != nil {
			return err
		}
		m.topics += "  " + topic + "  "
	}

	err = m.channel.Publish(topic, topic, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	//发布消息失败
	if err!=nil{
		ExampleLoggerOutput("消息推送到MQ失败！！")
	}
	ExampleLoggerOutput("消息"+msg+"发送成功！！")
	return nil
}

//关闭mq
func (m *MessageManageMq) Shutdown() error {
	// will close() the deliveries channel
	if err := m.channel.Cancel(m.tag, true); err != nil {
		return fmt.Errorf("manageMq cancel failed: %s", err)
	}

	if err := m.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer common.Log.Infof("AMQP shutdown Message Manage  OK")

	// wait for handle() to exit
	return <-m.done
}

//接收消息

func (m *MessageManageMq) ReceiveMessage(queueName string) error {
	common.Log.Errorf("Queue bound to Exchange, starting Consume (consumer tag %q)", m.tag)
	deliveries, err := m.channel.Consume(
		queueName, // name
		m.tag,     // consumerTag,
		false,     // noAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		err = fmt.Errorf("Queue Consume: %s", err)
		common.Log.Errorln(err)
		return err
	}
	go Handle(deliveries, m.done)

	return nil
}

// 测试连接是否正常
func (m*MessageManageMq)Ping() (err error) {
	common.Log.Infoln("测试rabbitMq 是否已经链接")
	if m.channel == nil {
		err:= errors.New("RabbitMQ is not initialize")
		return err
	}

	err = m.channel.ExchangeDeclare("ping.ping", "topic", false, true, false, true, nil)
	if err != nil {
		return err
	}

	msgContent := "ping.ping"

	err = m.channel.Publish("ping.ping", "ping.ping", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msgContent),
	})

	if err != nil {
		return err
	}

	err = m.channel.ExchangeDelete("ping.ping", false, false)

	return err
}

//消息发送
func Handle(deliveries <-chan amqp.Delivery, done chan error) {
	log.Println("此处一直阻塞等待获取 mq 中的消息")
	for d := range deliveries {
		common.Log.Infof(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		//go new(Consumer).UnmarshalMQBody(d.Body)
		//修改数据库的字段值么还是插入一条数据
		//解码 发送
		body :=&struct {
			Phone string `json:"phone"`
			Message string `json:"message"`
		}{}
		common.Log.Infoln("短息解码")
		err:=json.Unmarshal(d.Body,body)
		if err!=nil{
			common.Log.Errorln(err)
			return
		}
		err =new(sms.SMSMessage).SendMobileMessage(body.Phone,body.Message)
		if err!=nil{
			//发送失败如何处理
			//默认丢弃
		}
		d.Ack(false)
	}
	common.Log.Errorln("handle: deliveries channel closed")
	done <- nil
}

//临时封装临时用一下
func ExampleLoggerOutput(info string) {
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "INFO: ", log.Lshortfile)

		infoMessage = func(info string) {
			logger.Output(2, info)
		}
	)

	infoMessage(info)

	fmt.Print(&buf)
	// Output:
	// INFO: example_test.go:36: Hello world
}
