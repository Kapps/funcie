package internal

import "testing"

func TestMarshalArray(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		want string
	}{
		{
			name: "empty",
			s:    []string{},
			want: "[]",
		},
		{
			name: "single",
			s:    []string{"foo"},
			want: "[\"foo\"]",
		},
		{
			name: "multiple",
			s:    []string{"foo", "bar", "baz"},
			want: "[\"foo\", \"bar\", \"baz\"]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MarshalArray(tt.s); got != tt.want {
				t.Errorf("MarshalArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
