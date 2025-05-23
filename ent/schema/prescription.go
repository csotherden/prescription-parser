package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/google/uuid"
)

// Prescription holds the schema definition for the Prescription entity.
type Prescription struct {
	ent.Schema
}

// Fields of the Prescription.
func (Prescription) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("file_id"),
		field.String("mime_type"),
		field.JSON("content", models.Prescription{}),
	}
}

// Edges of the Prescription.
func (Prescription) Edges() []ent.Edge {
	return nil
}

// Mixin of the Prescription
func (Prescription) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
