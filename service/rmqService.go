package service

import (
	"errors"
	"github.com/koding/multiconfig"
	Log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sctek.com/typhoon/th-platform-gateway/common"
	"time"
)

//连接结构
type Connect struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

//信道结构
type Channel struct {
	Name     string `json:"name"`
	Connect  string `json:"connect"`
	QosCount int    `json:"qos_count"`
	QosSize  int    `json:"qos_size"`
}

//交换机绑定结构
type EBind struct {
	Destination string `json:"destination"`
	Key         string `json:"key"`
	NoWait      bool   `json:"no_wait"`
}

//交换机结构
type Exchange struct {
	Name        string                 `json:"name"`
	Channel     string                 `json:"channel"`
	Type        string                 `json:"type" default:"direct"`
	Durable     bool                   `json:"durable" `
	AutoDeleted bool                   `json:"auto_deleted"`
	Internal    bool                   `json:"internal" `
	NoWait      bool                   `json:"no_wait" `
	Bind        []EBind                `json:"ebind"`
	Args        map[string]interface{} `json:"args"`
}

//队列绑定结构
type QBind struct {
	ExchangeName string `json:"exchange_name"`
	Key          string `json:"key"`
	NoWait       bool   `json:"no_wait"`
}

//队列结构
type Queue struct {
	Name       string                 `json:"name"`
	Channel    string                 `json:"channel"`
	Durable    bool                   `json:"durable"`
	AutoDelete bool                   `json:"auto_delete"`
	Exclusive  bool                   `json:"exclusive"`
	NoWait     bool                   `json:"no_wait"`
	Bind       []QBind                `json:"qbind"`
	Args       map[string]interface{} `json:"args"`
}

//发送者配置
type Pusher struct {
	Name         string `json:"name"`
	Channel      string `json:"channel"`
	Exchange     string `json:"exchange"`
	Key          string `json:"key" `
	Mandatory    bool   `json:"mandatory" `
	Immediate    bool   `json:"immediate" `
	ContentType  string `json:"content_type"`
	DeliveryMode uint8  `json:"delivery_mode"`
}

//接收者配置
type Popup struct {
	Name      string `json:"name"`
	QName     string `json:"q_name"`
	Channel   string `json:"channel"`
	Consumer  string `json:"consumer"`
	AutoACK   bool   `json:"auto_ack"`
	Exclusive bool   `json:"exclusive" `
	NoLocal   bool   `json:"no_local" `
	NoWait    bool   `json:"no_wait" `
}

//配置文件结构
type mqCfg struct {
	Connects  []Connect  `json:"connects"`
	Channels  []Channel  `json:"channels"`
	Exchanges []Exchange `json:"exchanges"`
	Queue     []Queue    `json:"queue"`
	Pusher    []Pusher   `json:"pusher"`
	Popup     []Popup    `json:"popup"`
}

var _Cfg *mqCfg = new(mqCfg)                                                     //配置文件对象
var _ConnectPool map[string]*amqp.Connection = make(map[string]*amqp.Connection) //连接名称:连接对象
var _ChannelPool map[string]*amqp.Channel = make(map[string]*amqp.Channel)       //信道名称:信道对象
var _ExchangePool map[string]string = make(map[string]string)                    //交换机名称:所属信道名称
var _QueuePool map[string]string = make(map[string]string)                       //队列名称:所属信道名称
var _Pusher map[string]Pusher = make(map[string]Pusher)                          //Pusher名称:Pusher配置
var _Poper map[string]Popup = make(map[string]Popup)                             //Poper名称:Poper配置

//读取配置文件
func loadCfg() (err error) {
	if err = _Cfg.load(); err != nil {
		return err
	}
	return nil
}

func (c *mqCfg) load() error {

	t := &multiconfig.TagLoader{}
	j := &multiconfig.JSONLoader{Path: common.CPath.MustValue("rmqPath", "path", "rmq.json")}
	m := multiconfig.MultiLoader(t, j)
	err := m.Load(c)
	return err
}

func CloseConnect(name string) (err error) {
	if _, ok := _ConnectPool[name]; ok {
		_ConnectPool[name].Close()
	}
	return nil
}

//创建连接
func CreateConnect(v Connect) (err error) {
	var connect *amqp.Connection
	if connect, err = amqp.Dial(v.Addr); err != nil {
		return err
	} else {
		if _, ok := _ConnectPool[v.Name]; !ok {
			_ConnectPool[v.Name] = connect
		} else {
			return errors.New("连接已存在\n")
		}
	}
	return nil
}

//初始化连接
func initConnect() (err error) {
	for _, v := range _Cfg.Connects {
		if err = CreateConnect(v); err != nil {
			return err
		}
	}
	return nil
}

//关闭信道
func CloseChannel(name string) (err error) {
	if _, ok := _ChannelPool[name]; ok {
		_ChannelPool[name].Close()
	}
	return nil
}

//创建信道
func CreateChannel(v Channel) (err error) {
	if _, ok := _ConnectPool[v.Connect]; !ok {
		return errors.New("连接不存在\n")
	}
	var channel *amqp.Channel
	if channel, err = _ConnectPool[v.Connect].Channel(); err != nil {
		return err
	} else {
		if _, ok := _ChannelPool[v.Name]; !ok {
			_ChannelPool[v.Name] = channel
		} else {
			return errors.New("信道已存在\n")
		}

	}
	//prefetchCount：消费者未确认消息的个数。
	//prefetchSize ：消费者未确认消息的大小。
	//global ：是否全局生效，true表示是。全局生效指的是针对当前connect里的所有channel都生效。
	if err = channel.Qos(1, 0, false); err != nil {
		return nil
	}
	return nil
}

//初始化信道
func initChannel() (err error) {
	for _, v := range _Cfg.Channels {
		if err = CreateChannel(v); err != nil {
			return err
		}
	}
	return nil
}

//删除交换机
func DeleteExchange(name string, ifUnused bool, noWait bool) (err error) {
	if _, ok := _ExchangePool[name]; ok {
		if _, ok := _ChannelPool[_ExchangePool[name]]; ok {
			if err = _ChannelPool[_ExchangePool[name]].ExchangeDelete(
				name, ifUnused, noWait); err != nil {
				return err
			}
		}
		delete(_ExchangePool, name)
	}
	return nil
}

//创建交换机
func CreateExchange(v Exchange) (err error) {
	if _, ok := _ChannelPool[v.Channel]; !ok {
		return errors.New("信道不存在\n")
	}
	//name:交换器的名称，对应图中exchangeName。
	//kind:也叫作type，表示交换器的类型。有四种常用类型：direct、fanout、topic、headers。
	//durable:是否持久化，true表示是。持久化表示会把交换器的配置存盘，当RMQ Server重启后，会自动加载交换器。
	//autoDelete:是否自动删除，true表示是。至少有一条绑定才可以触发自动删除，当所有绑定都与交换器解绑后，会自动删除此交换器。
	//internal:是否为内部，true表示是。客户端无法直接发送msg到内部交换器，只有交换器可以发送msg到内部交换器。
	//noWait:是否非阻塞，true表示是非阻塞。
	//args:直接写nil
	if err = _ChannelPool[v.Channel].ExchangeDeclare(v.Name, v.Type,
		false, true, false, false, nil); err != nil {
		return err
	} else {
		if _, ok := _ExchangePool[v.Name]; !ok {
			_ExchangePool[v.Name] = v.Channel
		} else {
			return errors.New("交换机已存在")
		}
	}
	//源交换器根据路由键&绑定键把msg转发到目的交换器。
	//destination：目的交换器，通常是内部交换器。
	//key：对应图中BandingKey，表示要绑定的键。
	//source：源交换器。
	//nowait：是否非阻塞，true表示是非阻塞
	//args：直接写nil
	for _, b := range v.Bind {
		if err = _ChannelPool[v.Channel].ExchangeBind(b.Destination, b.Key, v.Name, b.NoWait, nil); err != nil {
			return err
		}
	}
	return nil
}

//初始化交换机
func initExchange() (err error) {
	for _, v := range _Cfg.Exchanges {
		if err = CreateExchange(v); err != nil {
			return err
		}
	}
	return nil
}

//删除队列
func DeleteQueue(name string, ifUnused bool, ifEmpty bool, noWait bool) (err error) {
	if _, ok := _QueuePool[name]; ok {
		if _, ok := _ChannelPool[_QueuePool[name]]; ok {
			if _, err = _ChannelPool[_QueuePool[name]].QueueDelete(
				name, ifUnused, ifEmpty, noWait); err != nil {
				return err
			}
		}
		delete(_QueuePool, name)
	}
	return nil
}

//创建队列
func CreateQueue(v Queue) (err error) {
	if _, ok := _ChannelPool[v.Channel]; !ok {
		return errors.New("信道不存在\n")
	}

	//处理x-message-ttl的类型，json里面写的是int，go读出来的是double
	//if _, ok := v.Args["x-message-ttl"]; ok {
	//	t := int32(v.Args["x-message-ttl"].(float64))
	//	delete(v.Args, "x-message-ttl")
	//	v.Args["x-message-ttl"] = t
	//}

	//name：队列名称
	//durable：是否持久化，true为是。持久化会把队列存盘，服务器重启后，不会丢失队列以及队列内的信息。（注：1、不丢失是相对的，如果宕机时有消息没来得及存盘，还是会丢失的。2、存盘影响性能。）
	//autoDelete：是否自动删除，true为是。至少有一个消费者连接到队列时才可以触发。当所有消费者都断开时，队列会自动删除。
	//exclusive：是否设置排他，true为是。如果设置为排他，则队列仅对首次声明他的连接可见，并在连接断开时自动删除。（注意，这里说的是连接不是信道，相同连接不同信道是可见的）。
	//nowait：是否非阻塞，true表示是非阻塞。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。
	//args：直接写nil
	if _, err = _ChannelPool[v.Channel].QueueDeclare(v.Name, false,
		true, false, false, nil); err != nil {
		return err
	} else {
		if _, ok := _QueuePool[v.Name]; !ok {
			_QueuePool[v.Name] = v.Channel
		} else {
			return errors.New("队列已存在")
		}
	}
	//name：队列名称
	//key：表示要绑定的键。
	//exchange：交换器名称
	//nowait：是否非阻塞，true表示是非阻塞。
	//args：直接写nil。

	for _, b := range v.Bind {
		if err = _ChannelPool[v.Channel].QueueBind(v.Name, b.Key, b.ExchangeName, b.NoWait, nil); err != nil {
			return err
		}
	}
	return nil
}

//初始化队列
func initQueue() (err error) {
	for _, v := range _Cfg.Queue {
		if err = CreateQueue(v); err != nil {
			return err
		}
	}
	return nil
}

//创建Pusher
func CreatePusher(v Pusher) (err error) {
	if _, ok := _Pusher[v.Name]; !ok {
		_Pusher[v.Name] = v
	} else {
		return errors.New("Pusher已存在")
	}
	return nil
}

//删除Pusher
func DeletePusher(name string) (err error) {
	if _, ok := _Poper[name]; ok {
		delete(_Poper, name)
	}
	return nil
}

//初始化Pusher
func initPusher() (err error) {
	for _, v := range _Cfg.Pusher {
		if err = CreatePusher(v); err != nil {
			return err
		}
	}
	return err
}

//创建Poper
func CreatePoper(v Popup) (err error) {
	if _, ok := _Poper[v.Name]; !ok {
		_Poper[v.Name] = v
	} else {
		return errors.New("Poper已存在")
	}
	return err
}

//删除Poper
func DeletePoper(name string) (err error) {
	if _, ok := _Poper[name]; ok {
		delete(_Poper, name)
	}
	return nil
}

//初始化Poper
func initPoper() (err error) {
	for _, v := range _Cfg.Popup {
		if err = CreatePoper(v); err != nil {
			return err
		}
	}
	return nil
}

//关闭
func Fini() (err error) {
	for _, conn := range _ConnectPool {
		for _, ch := range _ChannelPool {
			if err = ch.Close(); err != nil {
				return err
			}
		}
		if err = conn.Close(); err != nil {
			return err
		}
	}
	//清空所有缓存
	_Cfg = new(mqCfg)                                //配置文件对象
	_ConnectPool = make(map[string]*amqp.Connection) //连接名称:连接对象
	_ChannelPool = make(map[string]*amqp.Channel)    //信道名称:信道对象
	_ExchangePool = make(map[string]string)          //交换机名称:所属信道名称
	_QueuePool = make(map[string]string)             //队列名称:所属信道名称
	_Pusher = make(map[string]Pusher)                //Pusher名称:Pusher配置
	_Poper = make(map[string]Popup)                  //Popup名称:Popup配置

	return nil
}

//初始化
func Init() (err error) {
	if err = loadCfg(); err != nil {
		return err
	}
	if err = initConnect(); err != nil {
		return err
	}
	if err = initChannel(); err != nil {
		return err
	}
	if err = initExchange(); err != nil {
		return err
	}
	if err = initQueue(); err != nil {
		return err
	}
	if err = initPusher(); err != nil {
		return err
	}
	if err = initPoper(); err != nil {
		return err
	}
	return err
}

func Receive() error {
	if err := Pop("Poper", callback); err != nil {
		return err
	}

	if err := Pop("weChat", weChatCallback); err != nil {
		return err
	}
	//
	//if err := Pop("dlxPoper", dlxCallback); err != nil {
	//	return err
	//}
	if err := Pop("DBDataId", otherCallback); err != nil {
		return err
	}
	return nil
}

//exchange：要发送到的交换机名称，。
//key：路由键，对应图中RoutingKey。
//mandatory：直接false。
//immediate ：直接false。
//向交换机推送一条消息
func Push(name string, key string, msg []byte) (err error) {
	if _, ok := _Pusher[name]; !ok {
		return errors.New("Pusher不存在")
	}

	cfg := _Pusher[name]
	if key != "" {
		cfg.Key = key
	}
	if _, ok := _ChannelPool[cfg.Channel]; !ok {
		return errors.New("Channel不存在")
	}

	if err = _ChannelPool[cfg.Channel].Publish(cfg.Exchange, cfg.Key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: msg}); err != nil {
		return err
	}
	return nil
}

type MSG struct {
	Body    []byte
	Tag     uint64
	Channel string
	Poper   string
}

func (m MSG) Ack(multiple bool) (err error) {
	if _, ok := _ChannelPool[m.Channel]; !ok {
		return errors.New("Ack失败,Channel无效")
	} else {
		_ChannelPool[m.Channel].Ack(m.Tag, multiple)
	}
	return nil
}

//处理消息(顺序处理,如果需要多线程可以在回调函数中做手脚)
func handleMsg(msgs <-chan amqp.Delivery, callback func(MSG), channel string, popupName string) {
	//fmt.Println("等待消息中......")
	for d := range msgs {
		var msg MSG = MSG{
			Body:    d.Body,
			Tag:     d.DeliveryTag,
			Channel: channel,
			Poper:   popupName,
		}

		callback(msg)
		err := d.Ack(false)
		if err != nil {
			Log.Errorln(err)
		}
		time.Sleep(time.Nanosecond)
	}
}

//queue:队列名称。
//consumer:消费者标签，用于区分不同的消费者。
//autoAck:是否自动回复ACK，true为是，回复ACK表示高速服务器我收到消息了。建议为false，手动回复，这样可控性强。
//exclusive:设置是否排他，排他表示当前队列只能给一个消费者使用。
//noLocal:如果为true，表示生产者和消费者不能是同一个connect。
//nowait：是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。
//args：直接写nil，没研究过，不解释。
//注意下返回值：返回一个<- chan Delivery类型，遍历返回值，有消息则往下走， 没有则阻塞。
//---------------------
//从队列获取消息 -- 推模式
func Pop(name string, callback func(MSG)) (err error) {
	if _, ok := _Poper[name]; !ok {
		return errors.New("Poper不存在")
	}
	cfg := _Poper[name]
	if _, ok := _ChannelPool[cfg.Channel]; !ok {
		return errors.New("Channel不存在")
	}
	var msgs <-chan amqp.Delivery
	if msgs, err = _ChannelPool[cfg.Channel].Consume(cfg.QName, cfg.Consumer,
		false, false, false, false, nil); err != nil {
		return err
	}
	go handleMsg(msgs, callback, cfg.Channel, name)

	return nil
}
