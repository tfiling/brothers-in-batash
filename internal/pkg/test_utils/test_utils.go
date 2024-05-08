package test_utils

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func WrapStructWithReader(t *testing.T, instance interface{}) io.Reader {
	rawBody, err := json.Marshal(instance)
	require.NoError(t, err)
	return bytes.NewReader(rawBody)
}
