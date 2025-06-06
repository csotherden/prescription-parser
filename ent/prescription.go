// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/csotherden/prescription-parser/ent/prescription"
	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/google/uuid"
)

// Prescription is the model entity for the Prescription schema.
type Prescription struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// FileID holds the value of the "file_id" field.
	FileID string `json:"file_id,omitempty"`
	// MimeType holds the value of the "mime_type" field.
	MimeType string `json:"mime_type,omitempty"`
	// Content holds the value of the "content" field.
	Content      models.Prescription `json:"content,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Prescription) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case prescription.FieldContent:
			values[i] = new([]byte)
		case prescription.FieldFileID, prescription.FieldMimeType:
			values[i] = new(sql.NullString)
		case prescription.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case prescription.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Prescription fields.
func (pr *Prescription) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case prescription.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pr.ID = *value
			}
		case prescription.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pr.CreatedAt = value.Time
			}
		case prescription.FieldFileID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field file_id", values[i])
			} else if value.Valid {
				pr.FileID = value.String
			}
		case prescription.FieldMimeType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field mime_type", values[i])
			} else if value.Valid {
				pr.MimeType = value.String
			}
		case prescription.FieldContent:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field content", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &pr.Content); err != nil {
					return fmt.Errorf("unmarshal field content: %w", err)
				}
			}
		default:
			pr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Prescription.
// This includes values selected through modifiers, order, etc.
func (pr *Prescription) Value(name string) (ent.Value, error) {
	return pr.selectValues.Get(name)
}

// Update returns a builder for updating this Prescription.
// Note that you need to call Prescription.Unwrap() before calling this method if this Prescription
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *Prescription) Update() *PrescriptionUpdateOne {
	return NewPrescriptionClient(pr.config).UpdateOne(pr)
}

// Unwrap unwraps the Prescription entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pr *Prescription) Unwrap() *Prescription {
	_tx, ok := pr.config.driver.(*txDriver)
	if !ok {
		panic("ent: Prescription is not a transactional entity")
	}
	pr.config.driver = _tx.drv
	return pr
}

// String implements the fmt.Stringer.
func (pr *Prescription) String() string {
	var builder strings.Builder
	builder.WriteString("Prescription(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pr.ID))
	builder.WriteString("created_at=")
	builder.WriteString(pr.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("file_id=")
	builder.WriteString(pr.FileID)
	builder.WriteString(", ")
	builder.WriteString("mime_type=")
	builder.WriteString(pr.MimeType)
	builder.WriteString(", ")
	builder.WriteString("content=")
	builder.WriteString(fmt.Sprintf("%v", pr.Content))
	builder.WriteByte(')')
	return builder.String()
}

// Prescriptions is a parsable slice of Prescription.
type Prescriptions []*Prescription
