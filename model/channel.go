package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Channel struct {
	Notification string `gorm:"column:notification;type:varchar(256);not null;unique;primary_key"`
	Alias        string `gorm:"column:alias;type:varchar(256)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Channel) TableName() string {
	return "msa_activity_channel"
}

type ChannelDAO struct {
	db *gorm.DB
}

func NewChannelDAO(_conn *Conn) *ChannelDAO {
	dao := &ChannelDAO{}
	if nil != _conn {
		dao.db = _conn.db
	}
	return dao
}

func (this *ChannelDAO) Find(_notification string) (*Channel, error) {
	db := this.db
	if nil == db {
		sqlDB, err := openSqlDB()
		if nil != err {
			return nil, err
		}
		defer closeSqlDB(sqlDB)
		db = sqlDB
	}

	var channel Channel
	result := db.Where("notification = ?", _notification).First(&channel)
	if gorm.IsRecordNotFoundError(result.Error) {
		return &Channel{}, nil
	}
	return &channel, result.Error
}

func (this *ChannelDAO) Insert(_channel *Channel) error {
	db := this.db
	if nil == db {
		sqlDB, err := openSqlDB()
		if nil != err {
			return err
		}
		defer closeSqlDB(sqlDB)
		db = sqlDB
	}

	return db.Create(_channel).Error
}

func (this *ChannelDAO) Delete(_notification string) error {
	db := this.db
	if nil == db {
		sqlDB, err := openSqlDB()
		if nil != err {
			return err
		}
		defer closeSqlDB(sqlDB)
		db = sqlDB
	}

	return db.Delete(&Channel{Notification: _notification}).Error
}

func (this *ChannelDAO) List(_offset int64, _count int64) ([]*Channel, error) {
	db := this.db
	if nil == db {
		sqlDB, err := openSqlDB()
		if nil != err {
			return nil, err
		}
		defer closeSqlDB(sqlDB)
		db = sqlDB
	}

	var channels []*Channel
	res := db.Offset(_offset).Limit(_count).Order("created_at desc").Find(&channels)
	return channels, res.Error
}

func (this *ChannelDAO) ListAll() ([]*Channel, error) {
	db := this.db
	if nil == db {
		sqlDB, err := openSqlDB()
		if nil != err {
			return nil, err
		}
		defer closeSqlDB(sqlDB)
		db = sqlDB
	}

	var channels []*Channel
	res := db.Find(&channels)
	return channels, res.Error
}

func (this *ChannelDAO) Count() (int64, error) {
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
	res := db.Model(&Channel{}).Count(&count)
	return count, res.Error
}
