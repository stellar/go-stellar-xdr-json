package xdrjson

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	/* The base64-encoded string representing the asset
	Created with:
	$ stellar xdr encode --type Asset << -
	{"credit_alphanum4":{"asset_code":"ABCD","issuer":"GBRXG22ZWIZ4JZ3QNPQGCWCS2I4DJ4UH35Y4FHE32YCPGTIIYO6GCL64"}}
	-
	*/
	encodedAsset := "AAAAAUFCQ0QAAAAAY3NrWbIzxOdwa+BhWFLSODTyh99xwpyb1gTzTQjDvGE="

	rawBytes, err := base64.StdEncoding.DecodeString(encodedAsset)
	require.NoError(t, err)

	jsb, err := Decode(Asset, rawBytes)
	require.NoError(t, err)

	var dest map[string]interface{}
	require.NoError(t, json.Unmarshal(jsb, &dest))

	// Ensure the asset has the correct fields
	require.Contains(t, dest, "credit_alphanum4")
	require.Contains(t, dest["credit_alphanum4"], "asset_code")
	require.Contains(t, dest["credit_alphanum4"], "issuer")
	require.IsType(t, map[string]interface{}{}, dest["credit_alphanum4"])

	// Check the issuer address and asset code
	if converted, ok := dest["credit_alphanum4"].(map[string]interface{}); assert.True(t, ok) {
		require.Equal(t, "GBRXG22ZWIZ4JZ3QNPQGCWCS2I4DJ4UH35Y4FHE32YCPGTIIYO6GCL64", converted["issuer"])
		require.Equal(t, "ABCD", converted["asset_code"])
	}
}

func TestEmptyConversion(t *testing.T) {
	js, err := Decode(SorobanTransactionData, []byte{})
	require.NoError(t, err)
	require.Equal(t, "", string(js))
}
