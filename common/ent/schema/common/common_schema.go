package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// CommonSchema holds the schema definition for the CommonSchema entity.
type CommonSchema struct {
	ent.Schema
}

// Fields of the CommonSchema.
func (CommonSchema) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Unique().
			Immutable(),
	}
}

// Edges of the CommonSchema.
func (CommonSchema) Edges() []ent.Edge {
	return nil
}

// Indexes of the CommonSchema.
func (CommonSchema) Indexes() []ent.Index {
	return []ent.Index{}
}
