package schema

import (
	"user-management-servie/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.String("username"),
		field.String("email").Unique(),
	}
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{}, // Include the Time mixin
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
