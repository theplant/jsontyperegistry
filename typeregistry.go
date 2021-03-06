// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package typeregistry implements a simple type registry.
package jsontyperegistry

import (
	"errors"
	"reflect"
	"sort"
	"sync"
)

var (
	// ErrTypeRegistry is the base typeregistry error.
	ErrTypeRegistry = errors.New("typeregistry")
	// ErrNotFound is returned when a type was not found in the registry.
	ErrNotFound = errors.New("entry not found")
	// ErrDuplicateEntry is returned when registering a type that is already
	// registered.
	ErrDuplicateEntry = errors.New("entry already exists")
	// ErrInvalidParam is returned when an invalid parameter was passed to
	// a registration method.
	ErrInvalidParam = errors.New("invalid parameter")
)

// Registry is a simple Go type registry. Safe for concurrent use, no panic.
//
// Types can be registered under auto-generated or custom type names then
// retrieved under those names as a reflect Type, Value or an Interface.
type Registry struct {
	mu      sync.Mutex
	entries map[string]reflect.Type
}

// New returns a new *Registry instance.
func New() *Registry {
	p := &Registry{
		mu:      sync.Mutex{},
		entries: make(map[string]reflect.Type),
	}
	return p
}

// GetLongTypeName generates a long Type name from a Go value contained in i
// thet is constructed of pointer dereferene token for each pointer level of i,
// package path if the type was not predeclared or is an alias and the type
// name as returned by reflect.Value.Type().String(). For example:
// ***Registry => "***github.com/vedranvuk/typeregistry/typeregistry.Registry".
//
// It is used to generate a type name for values registered by this package.
func GetLongTypeName(i interface{}) (r string) {
	if i == nil {
		return
	}
	v := reflect.ValueOf(i)
	for v.Kind() == reflect.Ptr && !v.IsZero() {
		r += "*"
		v = v.Elem()
	}
	if s := v.Type().PkgPath(); s != "" {
		r += s + "/"
	}
	r += v.Type().String()
	return
}

// Register registers a reflect.Type of value specified by v under a name
// generated by GetLongTypeName or returns an error.
func (r *Registry) Register(v interface{}) error {
	return r.RegisterNamed(GetLongTypeName(v), v)
}

// RegisterNamed registers a reflect.Type of value specified by v under
// specified name, which cannot be empty, or returns an error.
func (r *Registry) RegisterNamed(name string, v interface{}) error {
	if name == "" || v == nil {
		return ErrInvalidParam
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.entries[name]; exists {
		return ErrDuplicateEntry
	}
	r.entries[name] = reflect.TypeOf(v)
	return nil
}

// Unregister unregisters a reflect.Type registered under specified name or
// returns an error.
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.entries[name]; !exists {
		return ErrNotFound
	}
	delete(r.entries, name)
	return nil
}

// GetType returns a registered reflect.Type specified by name or an error.
func (r *Registry) GetType(name string) (reflect.Type, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.entries[name]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

// GetValue returns a new reflect.Value of reflect.Type registered under
// specified name or an error.
func (r *Registry) GetValue(name string) (reflect.Value, error) {
	t, err := r.GetType(name)
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.New(t).Elem(), nil
}

// GetInterface returns an interface to a new reflect.Value of reflect.Type
// registered under specified name or an error.
func (r *Registry) GetInterface(name string) (interface{}, error) {
	t, err := r.GetType(name)
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.New(t).Elem().Interface(), nil
}

// RegisteredNames returns a slice of registered type names.
func (r *Registry) RegisteredNames() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	names := make([]string, 0, len(r.entries))
	for key := range r.entries {
		names = append(names, key)
	}
	sort.Strings(names)
	return names
}
