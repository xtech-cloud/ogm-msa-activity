package handler

import (
	"context"
	"omo-msa-activity/model"
	"omo-msa-activity/subscriber"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-activity/proto/activity"
)

type Channel struct{}

func (this *Channel) Subscribe(_ctx context.Context, _req *proto.ChannelSubRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Channel.Subscribe, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Notification {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "notification is required"
		return nil
	}

	dao := model.NewChannelDAO(model.DefaultConn)
	channel, err := dao.Find(_req.Notification)
	if nil != err {
		return err
	}

	if "" != channel.Notification {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "channel is exists"
		return nil
	}

	err = dao.Insert(&model.Channel{
		Notification: _req.Notification,
	})
	if nil != err {
		return err
	}
	subscriber.Add(_req.Notification)
	return nil
}

func (this *Channel) Unsubscribe(_ctx context.Context, _req *proto.ChannelUnsubRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Channel.Unsubscribe, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Notification {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "notification is required"
		return nil
	}

	dao := model.NewChannelDAO(model.DefaultConn)
	err := dao.Delete(_req.Notification)
	if nil != err {
		return err
	}
	subscriber.Remove(_req.Notification)
	return nil
}

func (this *Channel) Fetch(_ctx context.Context, _req *proto.ChannelFetchRequest, _rsp *proto.ChannelFetchResponse) error {
	logger.Infof("Received Channel.Fetch, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	if _req.Offset > 0 {
		offset = _req.Offset
	}

	count := int64(0)
	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewChannelDAO(model.DefaultConn)
	channels, err := dao.List(offset, count)
	if nil != err {
		return err
	}

	total, err := dao.Count()
	if nil != err {
		return err
	}

	_rsp.Channel = make([]*proto.ChannelEntity, len(channels))
	_rsp.Total = total

	for idx, v := range channels {
		_rsp.Channel[idx] = &proto.ChannelEntity{
			Notification: v.Notification,
			Alias:        v.Alias,
		}
	}

	return nil
}
