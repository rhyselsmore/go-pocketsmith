package pocketsmith

import (
	"bytes"
	"encoding/json"
)

// Param represents an optional request parameter.
//
// The zero value is unset and is omitted from a struct field tagged with
// json:",omitzero". Use P to set a parameter, including when its value is the
// zero value for T.
type Param[T any] struct {
	value T
	set   bool
}

// P returns a set Param containing value.
func P[T any](value T) Param[T] {
	return Param[T]{
		value: value,
		set:   true,
	}
}

// Get returns the parameter's value and whether it is set.
func (p Param[T]) Get() (T, bool) {
	return p.value, p.set
}

// IsSet reports whether the parameter is set.
func (p Param[T]) IsSet() bool {
	return p.set
}

// IsZero reports whether the parameter is unset. It allows encoding/json's
// omitzero option to omit unset parameters from containing structs.
func (p Param[T]) IsZero() bool {
	return !p.set
}

// MarshalJSON returns null for an unset parameter and otherwise marshals its
// value. An unset field tagged with json:",omitzero" is omitted before this
// method is called.
func (p Param[T]) MarshalJSON() ([]byte, error) {
	if !p.set {
		return []byte("null"), nil
	}
	return json.Marshal(p.value)
}

// UnmarshalJSON sets the parameter from a JSON value. A JSON null resets the
// parameter to its unset state.
func (p *Param[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		var zero Param[T]
		*p = zero
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	p.value = value
	p.set = true
	return nil
}
