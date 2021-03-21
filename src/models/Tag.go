package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TTag 标签表
const TTag = "t_tag"

// Tag 标签schema
type Tag struct {
	ID       *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`            // id
	UID      *primitive.ObjectID `json:"uid,omitempty" bson:"uid,omitempty"`           // uid
	Name     *string             `json:"name,omitempty" bson:"name,omitempty"`         // 标签名
	Color    *string             `json:"color,omitempty" bson:"color,omitempty"`       // 颜色
	CreateAt *time.Time          `json:"createAt,omitempty" bson:"createAt,omitempty"` // 创建时间
	UpdateAt *time.Time          `json:"updateAt,omitempty" bson:"updateAt,omitempty"` // 更新时间
}
