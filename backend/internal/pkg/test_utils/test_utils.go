package test_utils

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestTimeout can be edited into a longer timout for local test debug
const TestTimeout = int(600 * time.Second)

func WrapStructWithReader(t *testing.T, instance interface{}) io.Reader {
	rawBody, err := json.Marshal(instance)
	require.NoError(t, err)
	return bytes.NewReader(rawBody)
}
