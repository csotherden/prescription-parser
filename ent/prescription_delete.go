// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/csotherden/prescription-parser/ent/predicate"
	"github.com/csotherden/prescription-parser/ent/prescription"
)

// PrescriptionDelete is the builder for deleting a Prescription entity.
type PrescriptionDelete struct {
	config
	hooks    []Hook
	mutation *PrescriptionMutation
}

// Where appends a list predicates to the PrescriptionDelete builder.
func (pd *PrescriptionDelete) Where(ps ...predicate.Prescription) *PrescriptionDelete {
	pd.mutation.Where(ps...)
	return pd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pd *PrescriptionDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, pd.sqlExec, pd.mutation, pd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pd *PrescriptionDelete) ExecX(ctx context.Context) int {
	n, err := pd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pd *PrescriptionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(prescription.Table, sqlgraph.NewFieldSpec(prescription.FieldID, field.TypeUUID))
	if ps := pd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pd.mutation.done = true
	return affected, err
}

// PrescriptionDeleteOne is the builder for deleting a single Prescription entity.
type PrescriptionDeleteOne struct {
	pd *PrescriptionDelete
}

// Where appends a list predicates to the PrescriptionDelete builder.
func (pdo *PrescriptionDeleteOne) Where(ps ...predicate.Prescription) *PrescriptionDeleteOne {
	pdo.pd.mutation.Where(ps...)
	return pdo
}

// Exec executes the deletion query.
func (pdo *PrescriptionDeleteOne) Exec(ctx context.Context) error {
	n, err := pdo.pd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{prescription.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pdo *PrescriptionDeleteOne) ExecX(ctx context.Context) {
	if err := pdo.Exec(ctx); err != nil {
		panic(err)
	}
}
