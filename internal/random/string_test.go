package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantLen int
	}{
		{
			name:    "random_string_size_6",
			size:    6,
			wantLen: 6,
		},
		{
			name:    "random_string_size_9",
			size:    9,
			wantLen: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := NewRandomString(tt.size)
			assert.Equal(t, tt.wantLen, len(str))
		})
	}
}
