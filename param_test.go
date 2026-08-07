package pocketsmith

import (
	"encoding/json"
	"testing"
)

func TestP(t *testing.T) {
	p := P("")

	if !p.IsSet() {
		t.Fatal("P returned an unset Param")
	}

	value, ok := p.Get()
	if !ok {
		t.Fatal("Get reported that Param was unset")
	}
	if value != "" {
		t.Fatalf("Get returned %q, want an empty string", value)
	}
}

func TestParamZeroValueIsUnset(t *testing.T) {
	var p Param[string]

	if p.IsSet() {
		t.Fatal("zero Param is set")
	}

	value, ok := p.Get()
	if ok {
		t.Fatal("Get reported that zero Param was set")
	}
	if value != "" {
		t.Fatalf("Get returned %q, want an empty string", value)
	}
}

func TestParamMarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		param any
		want  string
	}{
		{
			name:  "unset",
			param: Param[string]{},
			want:  "null",
		},
		{
			name:  "string",
			param: P("memo"),
			want:  `"memo"`,
		},
		{
			name:  "empty string",
			param: P(""),
			want:  `""`,
		},
		{
			name:  "false",
			param: P(false),
			want:  "false",
		},
		{
			name:  "zero",
			param: P(0),
			want:  "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.param)
			if err != nil {
				t.Fatalf("Marshal returned an error: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("Marshal returned %s, want %s", got, tt.want)
			}
		})
	}
}

func TestParamOmitZero(t *testing.T) {
	type payload struct {
		Memo       Param[string]  `json:"memo,omitzero"`
		IsTransfer Param[bool]    `json:"is_transfer,omitzero"`
		Amount     Param[float64] `json:"amount,omitzero"`
	}

	tests := []struct {
		name    string
		payload payload
		want    string
	}{
		{
			name: "unset fields",
			want: `{}`,
		},
		{
			name: "set zero values",
			payload: payload{
				Memo:       P(""),
				IsTransfer: P(false),
				Amount:     P(0.0),
			},
			want: `{"memo":"","is_transfer":false,"amount":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("Marshal returned an error: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("Marshal returned %s, want %s", got, tt.want)
			}
		})
	}
}

func TestParamUnmarshalJSON(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		var p Param[bool]
		if err := json.Unmarshal([]byte("false"), &p); err != nil {
			t.Fatalf("Unmarshal returned an error: %v", err)
		}

		value, ok := p.Get()
		if !ok {
			t.Fatal("Unmarshal left Param unset")
		}
		if value {
			t.Fatal("Unmarshal returned true, want false")
		}
	})

	t.Run("null resets parameter", func(t *testing.T) {
		p := P("memo")
		if err := json.Unmarshal([]byte("null"), &p); err != nil {
			t.Fatalf("Unmarshal returned an error: %v", err)
		}

		if p.IsSet() {
			t.Fatal("Unmarshal of null left Param set")
		}
	})

	t.Run("invalid value preserves parameter", func(t *testing.T) {
		p := P(42)
		if err := json.Unmarshal([]byte(`"not a number"`), &p); err == nil {
			t.Fatal("Unmarshal returned nil error for an invalid value")
		}

		value, ok := p.Get()
		if !ok {
			t.Fatal("failed Unmarshal left Param unset")
		}
		if value != 42 {
			t.Fatalf("failed Unmarshal changed value to %d, want 42", value)
		}
	})
}

func TestUpdateTransactionParamsJSON(t *testing.T) {
	params := UpdateTransactionParams{
		Memo:         P(""),
		Amount:       P(0.0),
		IsTransfer:   P(false),
		CategoryID:   P[int64](123),
		NeedsReview:  P(true),
		ChequeNumber: P("456"),
	}

	got, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Marshal returned an error: %v", err)
	}

	const want = `{"memo":"","cheque_number":"456","amount":0,"is_transfer":false,"category_id":123,"needs_review":true}`
	if string(got) != want {
		t.Fatalf("Marshal returned %s, want %s", got, want)
	}
}
