package subscriber

import (
	"context"
	"omo-msa-activity/model"

	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	proto "github.com/xtech-cloud/omo-msp-notification/proto/notification"
)

type Notification struct {
}

var (
	daoRecord   *model.RecordDAO
	microServer server.Server
	subscribers map[string]server.Subscriber
)

func (this *Notification) Handle(_ctx context.Context, _message *proto.SimpleMessage) error {
	md, exists := metadata.FromContext(_ctx)

	operatorType := ""
	operatorLabel := ""
	topic := ""
	if exists {
		if _, ok := md["Operator-Type"]; ok {
			operatorType = md["Operator-Type"]
		}
		if _, ok := md["Operator-Label"]; ok {
			operatorLabel = md["Operator-Label"]
		}
		if _, ok := md["Micro-Topic"]; ok {
			topic = md["Micro-Topic"]
		}
	}

	if "" == topic {
		return nil
	}

	daoRecord.Insert(model.Record{
		Notification:  topic,
		OperatorLabel: operatorLabel,
		OperatorType:  operatorType,
		Action:        _message.Action,
		Head:          _message.Head,
		Body:          _message.Body,
	})
	return nil
}

func Setup(_server server.Server) {
	microServer = _server
	subscribers = make(map[string]server.Subscriber)
	// Register subscriber
	dao := model.NewChannelDAO(model.DefaultConn)
	channels, err := dao.ListAll()
	if err != nil {
		panic(err)
	}

	logger.Infof("found %d channels", len(channels))
	for _, channel := range channels {
		registerSubscriber(channel.Notification, new(Notification))
	}

	daoRecord = model.NewRecordDAO(model.DefaultConn)
}

func Add(_notification string) {
	//TODO 服务运行后再添加没有效果
	/*
		err := registerSubscriber(_notification, new(Notification))
		if nil != err {
			logger.Error(err)
		}
	*/
}

func Remove(_notification string) {
	//TODO 服务运行后再删除没有效果
	/*
		if _, ok := subscribers[_notification]; ok {
			//subscribers[_notification].Subscriber().Unsubscribe()
			delete(subscribers, _notification)
		}
	*/
}

func registerSubscriber(_topic string, _h interface{}, _opts ...server.SubscriberOption) error {
	sub := microServer.NewSubscriber(_topic, _h, _opts...)
	subscribers[_topic] = sub
	return microServer.Subscribe(sub)
}
