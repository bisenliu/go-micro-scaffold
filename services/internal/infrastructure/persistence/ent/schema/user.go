package schema

import (
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"

	uservo "services/internal/domain/user/valueobject"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Comment("用户ID"),
		field.String("name").
			MaxLen(50).
			NotEmpty().
			Comment("用户名"),
		field.String("open_id").
			Comment("open_id"),
		field.String("password").
			MaxLen(100).
			Sensitive().
			Comment("密码"),
		field.String("phone_number").
			Default("").
			Comment("手机号"),
		field.Int("gender").
			Validate(func(i int) error {
				switch i {
				case int(uservo.GenderMale), int(uservo.GenderFemale), int(uservo.GenderOther):
					return nil
				default:
					return errors.New("invalid gender value")
				}
			}).Comment("性别"),
		field.Time("created_at").
			Default(time.Now).
			Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("open_id").Unique(),
		index.Fields("phone_number").Unique(),
		index.Fields("created_at"),
	}
}
