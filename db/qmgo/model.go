package qmgo

import (
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// DefaultField defines the default fields to handle when operation happens
// import the DefaultField in document struct to make it working
type DefaultField struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	CreateAt time.Time          `bson:"created_at,omitempty"`
	UpdateAt time.Time          `bson:"updated_at,omitempty"`
}

func (df *DefaultField) CustomFields() field.CustomFieldsBuilder {
	return field.NewCustom().SetCreateAt("created_at").SetUpdateAt("updated_at").SetId("_id")
}

// DefaultUpdateAt changes the default updateAt field
func (df *DefaultField) DefaultUpdateAt() {
	df.UpdateAt = time.Now()
}

// DefaultCreateAt changes the default createAt field
func (df *DefaultField) DefaultCreateAt() {
	now := time.Now()
	if df.CreateAt.IsZero() {
		df.CreateAt = now
	}
	if df.UpdateAt.IsZero() {
		df.UpdateAt = now
	}
}

// DefaultId changes the default _id field
func (df *DefaultField) DefaultId() {
	if df.Id.IsZero() {
		df.Id = primitive.NewObjectID()
	}
}
