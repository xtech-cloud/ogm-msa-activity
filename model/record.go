package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Record struct {
	Embedded      gorm.Model `gorm:"embedded"`
	Notification  string     `gorm:"column:notification;type:varchar(256)"`
	OperatorLabel string     `gorm:"column:operator_label;type:varchar(256)"`
	OperatorType  string     `gorm:"column:operator_type;type:varchar(256)"`
	Action        string     `gorm:"column:action;type:varchar(256)"`
	Head          string     `gorm:"column:head;type:varchar(256)"`
	Body          string     `gorm:"column:body;type:TEXT"`
}

func (Record) TableName() string {
	return "msa_activity_record"
}

type RecordQuery struct {
	StartTime    int64
	EndTime      int64
	Notification string
	Action       string
}

type RecordDAO struct {
	db *gorm.DB
}

func NewRecordDAO(_conn *Conn) *RecordDAO {
	dao := &RecordDAO{}
	if nil != _conn {
		dao.db = _conn.db
	}
	return dao
}

func (this *RecordDAO) Insert(_record Record) error {
	db := this.db
	if nil == db {
		sqlDB, err := openSqlDB()
		if nil != err {
			return err
		}
		defer closeSqlDB(sqlDB)
		db = sqlDB
	}

	return db.Create(&_record).Error
}

func (this *RecordDAO) Query(_query RecordQuery) ([]*Record, error) {
	db := this.db
	if nil == db {
		sqlDB, err := openSqlDB()
		if nil != err {
			return nil, err
		}
		defer closeSqlDB(sqlDB)
		db = sqlDB
	}

	var records []*Record
	if "" != _query.Notification {
		db = db.Where("notification = ?", _query.Notification)
	}

	if "" != _query.Action {
		db = db.Where("actionn = ?", _query.Action)
	}

	if 0 != _query.StartTime {
		start := time.Unix(_query.StartTime, 0)
		db = db.Where("created_at > ?", start)
	}
	if 0 != _query.EndTime {
		end := time.Unix(_query.EndTime, 0)
		db = db.Where("created_at < ?", end)
	}
	res := db.Order("created_at desc").Find(&records)
	return records, res.Error
}
func (this *RecordDAO) Count() (int64, error) {
	db := this.db
	if nil == db {
		sqlDB, err := openSqlDB()
		if nil != err {
			return 0, err
		}
		defer closeSqlDB(sqlDB)
		db = sqlDB
	}

	count := int64(0)
	res := db.Model(&Record{}).Count(&count)
	return count, res.Error
}
