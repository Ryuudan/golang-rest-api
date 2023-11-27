package schema

import (
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
		field.String("first_name").
			StructTag(`json:"first_name" validate:"required,min=1"`),
		field.String("last_name").StructTag(`json:"last_name" validate:"required,min=1"`),
		field.String("middle_name").
			Optional().
			Nillable().
			StructTag(`json:"middle_name" validate:"omitempty,min=1"`),
		field.Time("birthday").
			StructTag(`json:"birthday" validate:"required"`).
			Nillable().
			Optional(),
		field.String("email").
			Unique().
			StructTag(`json:"email" validate:"required,email"`),
		field.String("phone_number").
			Unique().
			Optional().
			Nillable().
			StructTag(`json:"phone_number" validate:"e164"`),
		field.String("password").
			StructTag(`json:"password" validate:"required,min=3"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
