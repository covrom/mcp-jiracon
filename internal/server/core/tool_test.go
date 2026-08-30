package core

import "testing"

func TestSplitComma(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "simple list",
			in:   "A,B,C",
			want: []string{"A", "B", "C"},
		},
		{
			name: "with spaces",
			in:   "A, B , C",
			want: []string{"A", "B", "C"},
		},
		{
			name: "trailing comma",
			in:   "A,B,",
			want: []string{"A", "B"},
		},
		{
			name: "empty input",
			in:   "",
			want: []string{},
		},
		{
			name: "single value",
			in:   "PROJ-1",
			want: []string{"PROJ-1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitComma(tt.in)
			if len(got) != len(tt.want) {
				t.Errorf("splitComma(%q) = %v (len %d), want %v (len %d)", tt.in, got, len(got), tt.want, len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("splitComma(%q)[%d] = %q, want %q", tt.in, i, got[i], tt.want[i])
				}
			}
		})
	}
}
