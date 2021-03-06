// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebook/ent/dialect/gremlin"
	"github.com/facebook/ent/dialect/gremlin/graph/dsl"
	"github.com/facebook/ent/dialect/gremlin/graph/dsl/__"
	"github.com/facebook/ent/dialect/gremlin/graph/dsl/g"
	"github.com/facebook/ent/dialect/gremlin/graph/dsl/p"
	"github.com/facebook/ent/entc/integration/gremlin/ent/pet"
	"github.com/facebook/ent/entc/integration/gremlin/ent/predicate"
	"github.com/facebook/ent/entc/integration/gremlin/ent/user"
)

// PetUpdate is the builder for updating Pet entities.
type PetUpdate struct {
	config
	hooks      []Hook
	mutation   *PetMutation
	predicates []predicate.Pet
}

// Where adds a new predicate for the builder.
func (pu *PetUpdate) Where(ps ...predicate.Pet) *PetUpdate {
	pu.predicates = append(pu.predicates, ps...)
	return pu
}

// SetName sets the name field.
func (pu *PetUpdate) SetName(s string) *PetUpdate {
	pu.mutation.SetName(s)
	return pu
}

// SetTeamID sets the team edge to User by id.
func (pu *PetUpdate) SetTeamID(id string) *PetUpdate {
	pu.mutation.SetTeamID(id)
	return pu
}

// SetNillableTeamID sets the team edge to User by id if the given value is not nil.
func (pu *PetUpdate) SetNillableTeamID(id *string) *PetUpdate {
	if id != nil {
		pu = pu.SetTeamID(*id)
	}
	return pu
}

// SetTeam sets the team edge to User.
func (pu *PetUpdate) SetTeam(u *User) *PetUpdate {
	return pu.SetTeamID(u.ID)
}

// SetOwnerID sets the owner edge to User by id.
func (pu *PetUpdate) SetOwnerID(id string) *PetUpdate {
	pu.mutation.SetOwnerID(id)
	return pu
}

// SetNillableOwnerID sets the owner edge to User by id if the given value is not nil.
func (pu *PetUpdate) SetNillableOwnerID(id *string) *PetUpdate {
	if id != nil {
		pu = pu.SetOwnerID(*id)
	}
	return pu
}

// SetOwner sets the owner edge to User.
func (pu *PetUpdate) SetOwner(u *User) *PetUpdate {
	return pu.SetOwnerID(u.ID)
}

// Mutation returns the PetMutation object of the builder.
func (pu *PetUpdate) Mutation() *PetMutation {
	return pu.mutation
}

// ClearTeam clears the "team" edge to type User.
func (pu *PetUpdate) ClearTeam() *PetUpdate {
	pu.mutation.ClearTeam()
	return pu
}

// ClearOwner clears the "owner" edge to type User.
func (pu *PetUpdate) ClearOwner() *PetUpdate {
	pu.mutation.ClearOwner()
	return pu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (pu *PetUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(pu.hooks) == 0 {
		affected, err = pu.gremlinSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PetMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pu.mutation = mutation
			affected, err = pu.gremlinSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(pu.hooks) - 1; i >= 0; i-- {
			mut = pu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (pu *PetUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *PetUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *PetUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (pu *PetUpdate) gremlinSave(ctx context.Context) (int, error) {
	res := &gremlin.Response{}
	query, bindings := pu.gremlin().Query()
	if err := pu.driver.Exec(ctx, query, bindings, res); err != nil {
		return 0, err
	}
	if err, ok := isConstantError(res); ok {
		return 0, err
	}
	return res.ReadInt()
}

func (pu *PetUpdate) gremlin() *dsl.Traversal {
	type constraint struct {
		pred *dsl.Traversal // constraint predicate.
		test *dsl.Traversal // test matches and its constant.
	}
	constraints := make([]*constraint, 0, 1)
	v := g.V().HasLabel(pet.Label)
	for _, p := range pu.predicates {
		p(v)
	}
	var (
		rv = v.Clone()
		_  = rv

		trs []*dsl.Traversal
	)
	if value, ok := pu.mutation.Name(); ok {
		v.Property(dsl.Single, pet.FieldName, value)
	}
	if pu.mutation.TeamCleared() {
		tr := rv.Clone().InE(user.TeamLabel).Drop().Iterate()
		trs = append(trs, tr)
	}
	for _, id := range pu.mutation.TeamIDs() {
		v.AddE(user.TeamLabel).From(g.V(id)).InV()
		constraints = append(constraints, &constraint{
			pred: g.E().HasLabel(user.TeamLabel).OutV().HasID(id).Count(),
			test: __.Is(p.NEQ(0)).Constant(NewErrUniqueEdge(pet.Label, user.TeamLabel, id)),
		})
	}
	if pu.mutation.OwnerCleared() {
		tr := rv.Clone().InE(user.PetsLabel).Drop().Iterate()
		trs = append(trs, tr)
	}
	for _, id := range pu.mutation.OwnerIDs() {
		v.AddE(user.PetsLabel).From(g.V(id)).InV()
	}
	v.Count()
	if len(constraints) > 0 {
		constraints = append(constraints, &constraint{
			pred: rv.Count(),
			test: __.Is(p.GT(1)).Constant(&ConstraintError{msg: "update traversal contains more than one vertex"}),
		})
		v = constraints[0].pred.Coalesce(constraints[0].test, v)
		for _, cr := range constraints[1:] {
			v = cr.pred.Coalesce(cr.test, v)
		}
	}
	trs = append(trs, v)
	return dsl.Join(trs...)
}

// PetUpdateOne is the builder for updating a single Pet entity.
type PetUpdateOne struct {
	config
	hooks    []Hook
	mutation *PetMutation
}

// SetName sets the name field.
func (puo *PetUpdateOne) SetName(s string) *PetUpdateOne {
	puo.mutation.SetName(s)
	return puo
}

// SetTeamID sets the team edge to User by id.
func (puo *PetUpdateOne) SetTeamID(id string) *PetUpdateOne {
	puo.mutation.SetTeamID(id)
	return puo
}

// SetNillableTeamID sets the team edge to User by id if the given value is not nil.
func (puo *PetUpdateOne) SetNillableTeamID(id *string) *PetUpdateOne {
	if id != nil {
		puo = puo.SetTeamID(*id)
	}
	return puo
}

// SetTeam sets the team edge to User.
func (puo *PetUpdateOne) SetTeam(u *User) *PetUpdateOne {
	return puo.SetTeamID(u.ID)
}

// SetOwnerID sets the owner edge to User by id.
func (puo *PetUpdateOne) SetOwnerID(id string) *PetUpdateOne {
	puo.mutation.SetOwnerID(id)
	return puo
}

// SetNillableOwnerID sets the owner edge to User by id if the given value is not nil.
func (puo *PetUpdateOne) SetNillableOwnerID(id *string) *PetUpdateOne {
	if id != nil {
		puo = puo.SetOwnerID(*id)
	}
	return puo
}

// SetOwner sets the owner edge to User.
func (puo *PetUpdateOne) SetOwner(u *User) *PetUpdateOne {
	return puo.SetOwnerID(u.ID)
}

// Mutation returns the PetMutation object of the builder.
func (puo *PetUpdateOne) Mutation() *PetMutation {
	return puo.mutation
}

// ClearTeam clears the "team" edge to type User.
func (puo *PetUpdateOne) ClearTeam() *PetUpdateOne {
	puo.mutation.ClearTeam()
	return puo
}

// ClearOwner clears the "owner" edge to type User.
func (puo *PetUpdateOne) ClearOwner() *PetUpdateOne {
	puo.mutation.ClearOwner()
	return puo
}

// Save executes the query and returns the updated entity.
func (puo *PetUpdateOne) Save(ctx context.Context) (*Pet, error) {
	var (
		err  error
		node *Pet
	)
	if len(puo.hooks) == 0 {
		node, err = puo.gremlinSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PetMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			puo.mutation = mutation
			node, err = puo.gremlinSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(puo.hooks) - 1; i >= 0; i-- {
			mut = puo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, puo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (puo *PetUpdateOne) SaveX(ctx context.Context) *Pet {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *PetUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *PetUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (puo *PetUpdateOne) gremlinSave(ctx context.Context) (*Pet, error) {
	res := &gremlin.Response{}
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Pet.ID for update")}
	}
	query, bindings := puo.gremlin(id).Query()
	if err := puo.driver.Exec(ctx, query, bindings, res); err != nil {
		return nil, err
	}
	if err, ok := isConstantError(res); ok {
		return nil, err
	}
	pe := &Pet{config: puo.config}
	if err := pe.FromResponse(res); err != nil {
		return nil, err
	}
	return pe, nil
}

func (puo *PetUpdateOne) gremlin(id string) *dsl.Traversal {
	type constraint struct {
		pred *dsl.Traversal // constraint predicate.
		test *dsl.Traversal // test matches and its constant.
	}
	constraints := make([]*constraint, 0, 1)
	v := g.V(id)
	var (
		rv = v.Clone()
		_  = rv

		trs []*dsl.Traversal
	)
	if value, ok := puo.mutation.Name(); ok {
		v.Property(dsl.Single, pet.FieldName, value)
	}
	if puo.mutation.TeamCleared() {
		tr := rv.Clone().InE(user.TeamLabel).Drop().Iterate()
		trs = append(trs, tr)
	}
	for _, id := range puo.mutation.TeamIDs() {
		v.AddE(user.TeamLabel).From(g.V(id)).InV()
		constraints = append(constraints, &constraint{
			pred: g.E().HasLabel(user.TeamLabel).OutV().HasID(id).Count(),
			test: __.Is(p.NEQ(0)).Constant(NewErrUniqueEdge(pet.Label, user.TeamLabel, id)),
		})
	}
	if puo.mutation.OwnerCleared() {
		tr := rv.Clone().InE(user.PetsLabel).Drop().Iterate()
		trs = append(trs, tr)
	}
	for _, id := range puo.mutation.OwnerIDs() {
		v.AddE(user.PetsLabel).From(g.V(id)).InV()
	}
	v.ValueMap(true)
	if len(constraints) > 0 {
		v = constraints[0].pred.Coalesce(constraints[0].test, v)
		for _, cr := range constraints[1:] {
			v = cr.pred.Coalesce(cr.test, v)
		}
	}
	trs = append(trs, v)
	return dsl.Join(trs...)
}
