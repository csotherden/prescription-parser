// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/csotherden/prescription-parser/ent/embedding"
	"github.com/csotherden/prescription-parser/ent/predicate"
	"github.com/csotherden/prescription-parser/ent/prescription"
	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeEmbedding    = "Embedding"
	TypePrescription = "Prescription"
)

// EmbeddingMutation represents an operation that mutates the Embedding nodes in the graph.
type EmbeddingMutation struct {
	config
	op                  Op
	typ                 string
	id                  *uuid.UUID
	embedding           *pgvector.Vector
	clearedFields       map[string]struct{}
	prescription        *uuid.UUID
	clearedprescription bool
	done                bool
	oldValue            func(context.Context) (*Embedding, error)
	predicates          []predicate.Embedding
}

var _ ent.Mutation = (*EmbeddingMutation)(nil)

// embeddingOption allows management of the mutation configuration using functional options.
type embeddingOption func(*EmbeddingMutation)

// newEmbeddingMutation creates new mutation for the Embedding entity.
func newEmbeddingMutation(c config, op Op, opts ...embeddingOption) *EmbeddingMutation {
	m := &EmbeddingMutation{
		config:        c,
		op:            op,
		typ:           TypeEmbedding,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withEmbeddingID sets the ID field of the mutation.
func withEmbeddingID(id uuid.UUID) embeddingOption {
	return func(m *EmbeddingMutation) {
		var (
			err   error
			once  sync.Once
			value *Embedding
		)
		m.oldValue = func(ctx context.Context) (*Embedding, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Embedding.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withEmbedding sets the old Embedding of the mutation.
func withEmbedding(node *Embedding) embeddingOption {
	return func(m *EmbeddingMutation) {
		m.oldValue = func(context.Context) (*Embedding, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EmbeddingMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EmbeddingMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Embedding entities.
func (m *EmbeddingMutation) SetID(id uuid.UUID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *EmbeddingMutation) ID() (id uuid.UUID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *EmbeddingMutation) IDs(ctx context.Context) ([]uuid.UUID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uuid.UUID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Embedding.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetEmbedding sets the "embedding" field.
func (m *EmbeddingMutation) SetEmbedding(pg pgvector.Vector) {
	m.embedding = &pg
}

// Embedding returns the value of the "embedding" field in the mutation.
func (m *EmbeddingMutation) Embedding() (r pgvector.Vector, exists bool) {
	v := m.embedding
	if v == nil {
		return
	}
	return *v, true
}

// OldEmbedding returns the old "embedding" field's value of the Embedding entity.
// If the Embedding object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EmbeddingMutation) OldEmbedding(ctx context.Context) (v pgvector.Vector, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEmbedding is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEmbedding requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEmbedding: %w", err)
	}
	return oldValue.Embedding, nil
}

// ResetEmbedding resets all changes to the "embedding" field.
func (m *EmbeddingMutation) ResetEmbedding() {
	m.embedding = nil
}

// SetPrescriptionID sets the "prescription" edge to the Prescription entity by id.
func (m *EmbeddingMutation) SetPrescriptionID(id uuid.UUID) {
	m.prescription = &id
}

// ClearPrescription clears the "prescription" edge to the Prescription entity.
func (m *EmbeddingMutation) ClearPrescription() {
	m.clearedprescription = true
}

// PrescriptionCleared reports if the "prescription" edge to the Prescription entity was cleared.
func (m *EmbeddingMutation) PrescriptionCleared() bool {
	return m.clearedprescription
}

// PrescriptionID returns the "prescription" edge ID in the mutation.
func (m *EmbeddingMutation) PrescriptionID() (id uuid.UUID, exists bool) {
	if m.prescription != nil {
		return *m.prescription, true
	}
	return
}

// PrescriptionIDs returns the "prescription" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// PrescriptionID instead. It exists only for internal usage by the builders.
func (m *EmbeddingMutation) PrescriptionIDs() (ids []uuid.UUID) {
	if id := m.prescription; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetPrescription resets all changes to the "prescription" edge.
func (m *EmbeddingMutation) ResetPrescription() {
	m.prescription = nil
	m.clearedprescription = false
}

// Where appends a list predicates to the EmbeddingMutation builder.
func (m *EmbeddingMutation) Where(ps ...predicate.Embedding) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the EmbeddingMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *EmbeddingMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Embedding, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *EmbeddingMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *EmbeddingMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Embedding).
func (m *EmbeddingMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *EmbeddingMutation) Fields() []string {
	fields := make([]string, 0, 1)
	if m.embedding != nil {
		fields = append(fields, embedding.FieldEmbedding)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *EmbeddingMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case embedding.FieldEmbedding:
		return m.Embedding()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *EmbeddingMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case embedding.FieldEmbedding:
		return m.OldEmbedding(ctx)
	}
	return nil, fmt.Errorf("unknown Embedding field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *EmbeddingMutation) SetField(name string, value ent.Value) error {
	switch name {
	case embedding.FieldEmbedding:
		v, ok := value.(pgvector.Vector)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEmbedding(v)
		return nil
	}
	return fmt.Errorf("unknown Embedding field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *EmbeddingMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *EmbeddingMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *EmbeddingMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Embedding numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *EmbeddingMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *EmbeddingMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *EmbeddingMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Embedding nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *EmbeddingMutation) ResetField(name string) error {
	switch name {
	case embedding.FieldEmbedding:
		m.ResetEmbedding()
		return nil
	}
	return fmt.Errorf("unknown Embedding field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *EmbeddingMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.prescription != nil {
		edges = append(edges, embedding.EdgePrescription)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *EmbeddingMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case embedding.EdgePrescription:
		if id := m.prescription; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *EmbeddingMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *EmbeddingMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *EmbeddingMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedprescription {
		edges = append(edges, embedding.EdgePrescription)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *EmbeddingMutation) EdgeCleared(name string) bool {
	switch name {
	case embedding.EdgePrescription:
		return m.clearedprescription
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *EmbeddingMutation) ClearEdge(name string) error {
	switch name {
	case embedding.EdgePrescription:
		m.ClearPrescription()
		return nil
	}
	return fmt.Errorf("unknown Embedding unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *EmbeddingMutation) ResetEdge(name string) error {
	switch name {
	case embedding.EdgePrescription:
		m.ResetPrescription()
		return nil
	}
	return fmt.Errorf("unknown Embedding edge %s", name)
}

// PrescriptionMutation represents an operation that mutates the Prescription nodes in the graph.
type PrescriptionMutation struct {
	config
	op            Op
	typ           string
	id            *uuid.UUID
	created_at    *time.Time
	file_id       *string
	mime_type     *string
	content       *models.Prescription
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Prescription, error)
	predicates    []predicate.Prescription
}

var _ ent.Mutation = (*PrescriptionMutation)(nil)

// prescriptionOption allows management of the mutation configuration using functional options.
type prescriptionOption func(*PrescriptionMutation)

// newPrescriptionMutation creates new mutation for the Prescription entity.
func newPrescriptionMutation(c config, op Op, opts ...prescriptionOption) *PrescriptionMutation {
	m := &PrescriptionMutation{
		config:        c,
		op:            op,
		typ:           TypePrescription,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withPrescriptionID sets the ID field of the mutation.
func withPrescriptionID(id uuid.UUID) prescriptionOption {
	return func(m *PrescriptionMutation) {
		var (
			err   error
			once  sync.Once
			value *Prescription
		)
		m.oldValue = func(ctx context.Context) (*Prescription, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Prescription.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withPrescription sets the old Prescription of the mutation.
func withPrescription(node *Prescription) prescriptionOption {
	return func(m *PrescriptionMutation) {
		m.oldValue = func(context.Context) (*Prescription, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m PrescriptionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m PrescriptionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Prescription entities.
func (m *PrescriptionMutation) SetID(id uuid.UUID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *PrescriptionMutation) ID() (id uuid.UUID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *PrescriptionMutation) IDs(ctx context.Context) ([]uuid.UUID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uuid.UUID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Prescription.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreatedAt sets the "created_at" field.
func (m *PrescriptionMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *PrescriptionMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the Prescription entity.
// If the Prescription object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PrescriptionMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *PrescriptionMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetFileID sets the "file_id" field.
func (m *PrescriptionMutation) SetFileID(s string) {
	m.file_id = &s
}

// FileID returns the value of the "file_id" field in the mutation.
func (m *PrescriptionMutation) FileID() (r string, exists bool) {
	v := m.file_id
	if v == nil {
		return
	}
	return *v, true
}

// OldFileID returns the old "file_id" field's value of the Prescription entity.
// If the Prescription object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PrescriptionMutation) OldFileID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldFileID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldFileID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldFileID: %w", err)
	}
	return oldValue.FileID, nil
}

// ResetFileID resets all changes to the "file_id" field.
func (m *PrescriptionMutation) ResetFileID() {
	m.file_id = nil
}

// SetMimeType sets the "mime_type" field.
func (m *PrescriptionMutation) SetMimeType(s string) {
	m.mime_type = &s
}

// MimeType returns the value of the "mime_type" field in the mutation.
func (m *PrescriptionMutation) MimeType() (r string, exists bool) {
	v := m.mime_type
	if v == nil {
		return
	}
	return *v, true
}

// OldMimeType returns the old "mime_type" field's value of the Prescription entity.
// If the Prescription object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PrescriptionMutation) OldMimeType(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldMimeType is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldMimeType requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldMimeType: %w", err)
	}
	return oldValue.MimeType, nil
}

// ResetMimeType resets all changes to the "mime_type" field.
func (m *PrescriptionMutation) ResetMimeType() {
	m.mime_type = nil
}

// SetContent sets the "content" field.
func (m *PrescriptionMutation) SetContent(value models.Prescription) {
	m.content = &value
}

// Content returns the value of the "content" field in the mutation.
func (m *PrescriptionMutation) Content() (r models.Prescription, exists bool) {
	v := m.content
	if v == nil {
		return
	}
	return *v, true
}

// OldContent returns the old "content" field's value of the Prescription entity.
// If the Prescription object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PrescriptionMutation) OldContent(ctx context.Context) (v models.Prescription, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldContent is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldContent requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldContent: %w", err)
	}
	return oldValue.Content, nil
}

// ResetContent resets all changes to the "content" field.
func (m *PrescriptionMutation) ResetContent() {
	m.content = nil
}

// Where appends a list predicates to the PrescriptionMutation builder.
func (m *PrescriptionMutation) Where(ps ...predicate.Prescription) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the PrescriptionMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *PrescriptionMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Prescription, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *PrescriptionMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *PrescriptionMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Prescription).
func (m *PrescriptionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *PrescriptionMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.created_at != nil {
		fields = append(fields, prescription.FieldCreatedAt)
	}
	if m.file_id != nil {
		fields = append(fields, prescription.FieldFileID)
	}
	if m.mime_type != nil {
		fields = append(fields, prescription.FieldMimeType)
	}
	if m.content != nil {
		fields = append(fields, prescription.FieldContent)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *PrescriptionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case prescription.FieldCreatedAt:
		return m.CreatedAt()
	case prescription.FieldFileID:
		return m.FileID()
	case prescription.FieldMimeType:
		return m.MimeType()
	case prescription.FieldContent:
		return m.Content()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *PrescriptionMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case prescription.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case prescription.FieldFileID:
		return m.OldFileID(ctx)
	case prescription.FieldMimeType:
		return m.OldMimeType(ctx)
	case prescription.FieldContent:
		return m.OldContent(ctx)
	}
	return nil, fmt.Errorf("unknown Prescription field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *PrescriptionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case prescription.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case prescription.FieldFileID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFileID(v)
		return nil
	case prescription.FieldMimeType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMimeType(v)
		return nil
	case prescription.FieldContent:
		v, ok := value.(models.Prescription)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetContent(v)
		return nil
	}
	return fmt.Errorf("unknown Prescription field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *PrescriptionMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *PrescriptionMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *PrescriptionMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Prescription numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *PrescriptionMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *PrescriptionMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *PrescriptionMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Prescription nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *PrescriptionMutation) ResetField(name string) error {
	switch name {
	case prescription.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case prescription.FieldFileID:
		m.ResetFileID()
		return nil
	case prescription.FieldMimeType:
		m.ResetMimeType()
		return nil
	case prescription.FieldContent:
		m.ResetContent()
		return nil
	}
	return fmt.Errorf("unknown Prescription field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *PrescriptionMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *PrescriptionMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *PrescriptionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *PrescriptionMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *PrescriptionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *PrescriptionMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *PrescriptionMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Prescription unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *PrescriptionMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Prescription edge %s", name)
}
