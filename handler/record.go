package handler

import (
	"context"
	"omo-msa-activity/model"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-activity/proto/activity"
)

type Record struct{}

func (this *Record) Fetch(_ctx context.Context, _req *proto.RecordFetchRequest, _rsp *proto.RecordFetchResponse) error {
	logger.Infof("Received Record.Fetch, req is %v", _req)
	_rsp.Status = &proto.Status{}

	query := model.RecordQuery{
		StartTime:    _req.StartTime,
		EndTime:      _req.EndTime,
		Notification: _req.Notification,
		Action:       _req.Action,
	}

	dao := model.NewRecordDAO(model.DefaultConn)
	records, err := dao.Query(query)
	if nil != err {
		return err
	}

	_rsp.Record = make([]*proto.RecordEntity, len(records))

	for idx, v := range records {
		_rsp.Record[idx] = &proto.RecordEntity{
			Channel: &proto.ChannelEntity{
				Notification: v.Notification,
			},
			Operator: &proto.Operator{
				Label: v.OperatorLabel,
				Type:  v.OperatorType,
			},
			Action: v.Action,
			Head:   v.Head,
			Body:   v.Body,
			Time:   v.Embedded.CreatedAt.Unix(),
		}
	}

	return nil
}
