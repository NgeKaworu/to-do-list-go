package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TRecord 记录表
const TRecord = "t_record"

// Record 记录schema
type Record struct {
	ID       *primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`            // id
	UID      *primitive.ObjectID   `json:"uid,omitempty" bson:"uid,omitempty"`           // uid
	TID      *[]primitive.ObjectID `json:"tid,omitempty" bson:"tid,omitempty"`           // tid
	Event    *string               `json:"event,omitempty" bson:"event,omitempty"`       // 事件
	CreateAt *time.Time            `json:"createAt,omitempty" bson:"createAt,omitempty"` // 创建时间
	UpdateAt *time.Time            `json:"updateAt,omitempty" bson:"updateAt,omitempty"` // 更新时间
	Deration *time.Duration        `json:"deration,omitempty" bson:"deration,omitempty"` // 持续时间
}
