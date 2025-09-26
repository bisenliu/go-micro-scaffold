package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// BaseSchema 定义所有实体的公共字段
type BaseSchema struct {
	ent.Schema
}

// Fields of the BaseSchema.
func (BaseSchema) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Unique().
			Immutable(),
	}
}

// Edges of the BaseSchema.
func (BaseSchema) Edges() []ent.Edge {
	return nil
}

// Indexes of the BaseSchema.
func (BaseSchema) Indexes() []ent.Index {
	return []ent.Index{}
}
