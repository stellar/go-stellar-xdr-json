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

func TestDecode_empty(t *testing.T) {
	_, err := Decode(SorobanTransactionData, []byte{})
	require.Error(t, err)
}

func TestDecode_executableTag(t *testing.T) {
	/* The base64-encoded string representing the Protocol 28 SCVal
	Created with:
	$ stellar xdr encode --type ScVal << -
	{"executable_tag":"tag"}
	-
	*/
	encodedScVal := "AAAAFgAAAAN0YWcA"

	rawBytes, err := base64.StdEncoding.DecodeString(encodedScVal)
	require.NoError(t, err)

	jsb, err := Decode(ScVal, rawBytes)
	require.NoError(t, err)

	require.JSONEq(t, `{"executable_tag":"tag"}`, string(jsb))
}

// TestDecode_muxedContractAddress covers CAP-84, which adds
// SC_ADDRESS_TYPE_MUXED_CONTRACT = 5 to the SCAddress union along with the
// SEP-23 "W" strkey. Contract event SCVals are decoded through this path, so a
// transfer to a muxed contract reaches downstream consumers as XDR-JSON here.
func TestDecode_muxedContractAddress(t *testing.T) {
	/* An SCAddress on SC_ADDRESS_TYPE_MUXED_CONTRACT. Built by hand, since the
	union arm is nested one level deeper than the stellar CLI's shorthand:

	    00000005                    type = SC_ADDRESS_TYPE_MUXED_CONTRACT (5)
	    00000000499602d2            MuxedContract.id = 1234567890
	    0102...1f20                 MuxedContract.contractId, bytes 0x01..0x20
	*/
	encodedAddress := "AAAABQAAAABJlgLSAQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRobHB0eHyA="

	rawBytes, err := base64.StdEncoding.DecodeString(encodedAddress)
	require.NoError(t, err)

	jsb, err := Decode(ScAddress, rawBytes)
	require.NoError(t, err)

	// A muxed contract renders as a single "W" strkey whose payload is the 32
	// byte contract id followed by the 8 byte big endian mux id.
	var address string
	require.NoError(t, json.Unmarshal(jsb, &address))
	require.Equal(t, "WAAQEAYEAUDAOCAJBIFQYDIOB4IBCEQTCQKRMFYYDENBWHA5DYPSAAAAAAAETFQC2KVDM", address)
}

// TestDecode_stellarValueSignedMs covers CAP-88, which moves ledger close times
// to millisecond resolution and adds STELLAR_VALUE_SIGNED_MS = 3 to the
// StellarValue ext union. Once the previous ledger is at Protocol 30 the whole
// second arms are invalid, so this is the ordinary shape for every ledger.
func TestDecode_stellarValueSignedMs(t *testing.T) {
	/* A StellarValue whose ext is on arm 3:

	   cc * 32                     txSetHash
	   000000006adfe983            closeTime = 1793059203
	   00000000                    upgrades<6>, zero entries
	   00000003                    ext.v = STELLAR_VALUE_SIGNED_MS (3)
	   000001a17aa82833            closeTimeMs = 1793059203123
	   00000000 aa * 32            lcValueSignature.nodeID, ed25519
	   00000040 bb * 64            lcValueSignature.signature
	*/
	encodedValue := "zMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMwAAAAAat/pgwAAAAAAAAADAAABoXqoKDMAAAAAqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAABAu7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7uw=="

	rawBytes, err := base64.StdEncoding.DecodeString(encodedValue)
	require.NoError(t, err)

	jsb, err := Decode(StellarValue, rawBytes)
	require.NoError(t, err)

	var dest map[string]interface{}
	require.NoError(t, json.Unmarshal(jsb, &dest))

	// The top level close time is untouched by CAP-88 and stays whole seconds.
	require.Equal(t, "1793059203", dest["close_time"])

	ext, ok := dest["ext"].(map[string]interface{})
	require.True(t, ok)
	signedMs, ok := ext["signed_ms"].(map[string]interface{})
	require.True(t, ok, "ext should decode onto the signed_ms arm")

	require.Equal(t, "1793059203123", signedMs["close_time_ms"])

	// The validator signature is still carried on the new arm. Consumers that
	// reach for it by the old arm name silently find nothing, so assert it is
	// present and correct rather than just that the decode succeeded.
	sig, ok := signedMs["lc_value_signature"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "GCVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVH7N", sig["node_id"])
	require.NotEmpty(t, sig["signature"])
}

// TestDecode_stellarValueEmptyTxSetMs covers the other CAP-88 arm,
// STELLAR_VALUE_EMPTY_TX_SET_MS = 4, which validators use to propose dropping
// the transaction set. It also carries a signature, one level deeper again.
func TestDecode_stellarValueEmptyTxSetMs(t *testing.T) {
	/* A StellarValue whose ext is on arm 4:

	   cc * 32                     txSetHash
	   000000006adfe983            closeTime = 1793059203
	   00000000                    upgrades<6>, zero entries
	   00000004                    ext.v = STELLAR_VALUE_EMPTY_TX_SET_MS (4)
	   000001a17aa82833            closeTimeMs = 1793059203123
	   cc * 32                     txSetHash
	   dd * 32                     previousLedgerHash
	   0000001e                    previousLedgerVersion = 30
	   00000000 aa * 32            lcValueSignature.nodeID, ed25519
	   00000040 bb * 64            lcValueSignature.signature
	*/
	encodedValue := "zMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMwAAAAAat/pgwAAAAAAAAAEAAABoXqoKDPMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzN3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3d3dAAAAHgAAAACqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqgAAAEC7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7u7"

	rawBytes, err := base64.StdEncoding.DecodeString(encodedValue)
	require.NoError(t, err)

	jsb, err := Decode(StellarValue, rawBytes)
	require.NoError(t, err)

	var dest map[string]interface{}
	require.NoError(t, json.Unmarshal(jsb, &dest))

	ext, ok := dest["ext"].(map[string]interface{})
	require.True(t, ok)
	proposedMs, ok := ext["empty_tx_set_ms"].(map[string]interface{})
	require.True(t, ok, "ext should decode onto the empty_tx_set_ms arm")

	require.Equal(t, "1793059203123", proposedMs["close_time_ms"])
	require.EqualValues(t, 30, proposedMs["previous_ledger_version"])

	sig, ok := proposedMs["lc_value_signature"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "GCVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVKVH7N", sig["node_id"])
	require.NotEmpty(t, sig["signature"])
}
