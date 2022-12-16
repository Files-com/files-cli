package version

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Version
		wantErr bool
	}{
		{
			name: "1.2.1",
			args: "1.2.1",
			want: Version{Major: 1, Minor: 2, Patch: 1},
		},
		{
			name: "v1.2.1",
			args: "v1.2.1",
			want: Version{Major: 1, Minor: 2, Patch: 1},
		},
		{
			name: " v3.2.1",
			args: " v3.2.1\n",
			want: Version{Major: 3, Minor: 2, Patch: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_Greater(t *testing.T) {
	tests := []struct {
		name string
		Version
		args Version
		want bool
	}{
		{
			name:    "equal",
			Version: Version{Major: 3, Minor: 2, Patch: 1},
			args:    Version{Major: 3, Minor: 2, Patch: 1},
			want:    false,
		},

		{
			name:    "greater major",
			Version: Version{Major: 3, Minor: 2, Patch: 1},
			args:    Version{Major: 4, Minor: 2, Patch: 1},
			want:    true,
		},
		{
			name:    "greater minor",
			Version: Version{Major: 3, Minor: 2, Patch: 1},
			args:    Version{Major: 3, Minor: 3, Patch: 1},
			want:    true,
		},
		{
			name:    "greater patch",
			Version: Version{Major: 3, Minor: 2, Patch: 1},
			args:    Version{Major: 3, Minor: 2, Patch: 2},
			want:    true,
		},

		{
			name:    "lesser major",
			Version: Version{Major: 3, Minor: 2, Patch: 1},
			args:    Version{Major: 2, Minor: 2, Patch: 1},
			want:    false,
		},
		{
			name:    "lesser minor",
			Version: Version{Major: 3, Minor: 2, Patch: 1},
			args:    Version{Major: 3, Minor: 1, Patch: 1},
			want:    false,
		},
		{
			name:    "lesser patch",
			Version: Version{Major: 3, Minor: 2, Patch: 1},
			args:    Version{Major: 3, Minor: 2, Patch: 0},
			want:    false,
		},

		{
			name:    "greater patch high number",
			Version: Version{Major: 3, Minor: 2, Patch: 1},
			args:    Version{Major: 3, Minor: 2, Patch: 225},
			want:    true,
		},

		{
			name:    "lessor minor greater patch",
			Version: Version{Major: 2, Minor: 2, Patch: 8},
			args:    Version{Major: 2, Minor: 1, Patch: 38},
			want:    false,
		},
		{
			name:    "lessor major greater minor",
			Version: Version{Major: 2, Minor: 2, Patch: 8},
			args:    Version{Major: 1, Minor: 21, Patch: 8},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.Version.Greater(tt.args); got != tt.want {
				t.Errorf("Greater() = %v, want %v", got, tt.want)
			}
		})
	}
}
