package pocketsmith

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestLabels(t *testing.T) {
	labels := NewLabels("groceries", "weekly", "groceries")

	if got, want := labels.Len(), 2; got != want {
		t.Fatalf("Len returned %d, want %d", got, want)
	}
	if !labels.Has("groceries") {
		t.Fatal("Has reported that groceries was absent")
	}
	if labels.Has("missing") {
		t.Fatal("Has reported that missing was present")
	}

	if !labels.Add("budget") {
		t.Fatal("Add reported that a new label was already present")
	}
	if labels.Add("budget") {
		t.Fatal("Add reported that an existing label was newly added")
	}

	if !labels.Delete("weekly") {
		t.Fatal("Delete reported that an existing label was absent")
	}
	if labels.Delete("weekly") {
		t.Fatal("Delete reported that an absent label was present")
	}

	got := slices.Sorted(labels.All())
	want := []string{"budget", "groceries"}
	if !slices.Equal(got, want) {
		t.Fatalf("All returned %v, want %v", got, want)
	}

	labels.Clear()
	if got := labels.Len(); got != 0 {
		t.Fatalf("Clear left %d labels, want 0", got)
	}
}

func TestLabelsString(t *testing.T) {
	labels := NewLabels("weekly", "budget", "groceries")

	if got, want := labels.String(), "budget,groceries,weekly"; got != want {
		t.Fatalf("String returned %q, want %q", got, want)
	}
}

func TestLabelsMarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		labels Labels
		want   string
	}{
		{
			name:   "labels",
			labels: NewLabels("weekly", "budget", "groceries"),
			want:   `"budget,groceries,weekly"`,
		},
		{
			name:   "empty set",
			labels: NewLabels(),
			want:   `""`,
		},
		{
			name: "nil set",
			want: `""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.labels)
			if err != nil {
				t.Fatalf("Marshal returned an error: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("Marshal returned %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLabelsUnmarshalJSON(t *testing.T) {
	t.Run("array", func(t *testing.T) {
		var labels Labels
		if err := json.Unmarshal([]byte(`["weekly","budget","weekly"]`), &labels); err != nil {
			t.Fatalf("Unmarshal returned an error: %v", err)
		}

		got := slices.Sorted(labels.All())
		want := []string{"budget", "weekly"}
		if !slices.Equal(got, want) {
			t.Fatalf("Unmarshal produced %v, want %v", got, want)
		}
	})

	t.Run("null", func(t *testing.T) {
		labels := NewLabels("existing")
		if err := json.Unmarshal([]byte("null"), &labels); err != nil {
			t.Fatalf("Unmarshal returned an error: %v", err)
		}

		if labels == nil {
			t.Fatal("Unmarshal of null produced a nil set")
		}
		if got := labels.Len(); got != 0 {
			t.Fatalf("Unmarshal of null left %d labels, want 0", got)
		}
	})

	t.Run("invalid input preserves labels", func(t *testing.T) {
		labels := NewLabels("existing")
		if err := json.Unmarshal([]byte(`"not an array"`), &labels); err == nil {
			t.Fatal("Unmarshal returned nil error for invalid input")
		}

		if !labels.Has("existing") || labels.Len() != 1 {
			t.Fatalf("failed Unmarshal changed labels to %v", slices.Collect(labels.All()))
		}
	})
}

func TestUpdateTransactionParamsLabelsJSON(t *testing.T) {
	tests := []struct {
		name   string
		labels Param[Labels]
		want   string
	}{
		{
			name: "unset leaves labels unchanged",
			want: `{}`,
		},
		{
			name:   "empty set removes all labels",
			labels: P(NewLabels()),
			want:   `{"labels":""}`,
		},
		{
			name:   "set replaces all labels",
			labels: P(NewLabels("weekly", "budget", "groceries")),
			want:   `{"labels":"budget,groceries,weekly"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := UpdateTransactionParams{Labels: tt.labels}
			got, err := json.Marshal(params)
			if err != nil {
				t.Fatalf("Marshal returned an error: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("Marshal returned %s, want %s", got, tt.want)
			}
		})
	}
}
