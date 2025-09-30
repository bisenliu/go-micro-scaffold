package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CasbinRule holds the schema definition for the CasbinRule entity.
type CasbinRule struct {
	ent.Schema
}

// Annotations of the CasbinRule.
func (CasbinRule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
	}
}

// Fields of the CasbinRule.
func (CasbinRule) Fields() []ent.Field {
	return []ent.Field{
		field.String("ptype").
			MaxLen(100).
			NotEmpty().
			Comment("策略类型"),
		field.String("v0").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("策略字段0"),
		field.String("v1").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("策略字段1"),
		field.String("v2").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("策略字段2"),
		field.String("v3").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("策略字段3"),
		field.String("v4").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("策略字段4"),
		field.String("v5").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("策略字段5"),
		field.Time("created_at").
			Default(time.Now).
			Comment("创建时间"),
	}
}

// Edges of the CasbinRule.
func (CasbinRule) Edges() []ent.Edge {
	return nil
}

// Indexes of the CasbinRule.
func (CasbinRule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ptype"),
		index.Fields("ptype", "v0", "v1", "v2", "v3", "v4", "v5"),
	}
}
