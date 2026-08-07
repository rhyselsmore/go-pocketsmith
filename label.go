package pocketsmith

import (
	"encoding/json"
	"iter"
	"maps"
	"slices"
	"strings"
)

// Labels is a set of transaction labels.
//
// PocketSmith returns transaction labels as a JSON array, but accepts them in
// create and update requests as a comma-separated JSON string.
type Labels map[string]struct{}

// NewLabels returns a set containing values.
func NewLabels(values ...string) Labels {
	labels := make(Labels, len(values))
	for _, value := range values {
		labels.Add(value)
	}
	return labels
}

// Len returns the number of labels in the set.
func (s Labels) Len() int {
	return len(s)
}

// Has reports whether value is in the set.
func (s Labels) Has(value string) bool {
	_, ok := s[value]
	return ok
}

// Add inserts value into the set and reports whether it was newly added.
//
// Add panics if s is nil. Use NewLabels or a Labels literal to initialize it.
func (s Labels) Add(value string) bool {
	if s.Has(value) {
		return false
	}
	s[value] = struct{}{}
	return true
}

// Delete removes value from the set and reports whether it was present.
func (s Labels) Delete(value string) bool {
	if !s.Has(value) {
		return false
	}
	delete(s, value)
	return true
}

// Clear removes all labels from the set.
func (s Labels) Clear() {
	clear(s)
}

// All returns an iterator over the labels. Iteration order is unspecified.
func (s Labels) All() iter.Seq[string] {
	return maps.Keys(s)
}

// String returns the labels as a deterministic comma-separated string.
func (s Labels) String() string {
	return strings.Join(slices.Sorted(s.All()), ",")
}

// MarshalJSON encodes labels in PocketSmith's create and update request
// format: a comma-separated JSON string.
func (s Labels) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON decodes labels from PocketSmith's transaction response format:
// a JSON array of strings. A JSON null is treated as an empty set.
func (s *Labels) UnmarshalJSON(data []byte) error {
	var values []string
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	labels := NewLabels(values...)
	*s = labels
	return nil
}
