package crypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandom(t *testing.T) {
	r1 := NewPlaintextPrivateSecret()
	r2 := NewPlaintextPrivateSecret()
	require.NotEqual(t, r1, r2)
}
