package dbx

import "reflect"

type RelationKind int

const (
	RelationBelongsTo RelationKind = iota
	RelationHasOne
	RelationHasMany
	RelationManyToMany
)

type RelationMeta struct {
	Name                string
	FieldName           string
	Kind                RelationKind
	SourceTable         string
	SourceAlias         string
	TargetTable         string
	LocalColumn         string
	TargetColumn        string
	ThroughTable        string
	ThroughLocalColumn  string
	ThroughTargetColumn string
	TargetType          reflect.Type
}

type relationBinder interface {
	bindRelation(binding relationBinding) any
	relationKind() RelationKind
	targetType() reflect.Type
}

type relationBinding struct {
	meta RelationMeta
}

type relationAccessor interface {
	relationRef() RelationMeta
}

type BelongsTo[E any, T any] struct {
	meta RelationMeta
}

type HasOne[E any, T any] struct {
	meta RelationMeta
}

type HasMany[E any, T any] struct {
	meta RelationMeta
}

type ManyToMany[E any, T any] struct {
	meta RelationMeta
}

func (r BelongsTo[E, T]) bindRelation(binding relationBinding) any {
	r.meta = binding.meta
	return r
}

func (r HasOne[E, T]) bindRelation(binding relationBinding) any {
	r.meta = binding.meta
	return r
}

func (r HasMany[E, T]) bindRelation(binding relationBinding) any {
	r.meta = binding.meta
	return r
}

func (r ManyToMany[E, T]) bindRelation(binding relationBinding) any {
	r.meta = binding.meta
	return r
}

func (BelongsTo[E, T]) relationKind() RelationKind  { return RelationBelongsTo }
func (HasOne[E, T]) relationKind() RelationKind     { return RelationHasOne }
func (HasMany[E, T]) relationKind() RelationKind    { return RelationHasMany }
func (ManyToMany[E, T]) relationKind() RelationKind { return RelationManyToMany }

func (BelongsTo[E, T]) targetType() reflect.Type  { return reflect.TypeFor[T]() }
func (HasOne[E, T]) targetType() reflect.Type     { return reflect.TypeFor[T]() }
func (HasMany[E, T]) targetType() reflect.Type    { return reflect.TypeFor[T]() }
func (ManyToMany[E, T]) targetType() reflect.Type { return reflect.TypeFor[T]() }

func (r BelongsTo[E, T]) Name() string  { return r.meta.Name }
func (r HasOne[E, T]) Name() string     { return r.meta.Name }
func (r HasMany[E, T]) Name() string    { return r.meta.Name }
func (r ManyToMany[E, T]) Name() string { return r.meta.Name }

func (r BelongsTo[E, T]) Meta() RelationMeta  { return r.meta }
func (r HasOne[E, T]) Meta() RelationMeta     { return r.meta }
func (r HasMany[E, T]) Meta() RelationMeta    { return r.meta }
func (r ManyToMany[E, T]) Meta() RelationMeta { return r.meta }

func (r BelongsTo[E, T]) refNode()  {}
func (r HasOne[E, T]) refNode()     {}
func (r HasMany[E, T]) refNode()    {}
func (r ManyToMany[E, T]) refNode() {}

func (r BelongsTo[E, T]) relationRef() RelationMeta  { return r.meta }
func (r HasOne[E, T]) relationRef() RelationMeta     { return r.meta }
func (r HasMany[E, T]) relationRef() RelationMeta    { return r.meta }
func (r ManyToMany[E, T]) relationRef() RelationMeta { return r.meta }
