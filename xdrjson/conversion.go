//nolint:lll
package xdrjson

/*
// See preflight.go for add'l explanations:
// Note: no blank lines allowed.
#include <stdlib.h>
#include "../lib/xdrjson.h"
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/libs/x86_64-pc-windows-gnu -lxdrjson -lntdll -static -lws2_32 -lbcrypt -luserenv
#cgo darwin,amd64  LDFLAGS: -L${SRCDIR}/libs/x86_64-apple-darwin -lxdrjson -ldl -lm
#cgo darwin,arm64  LDFLAGS: -L${SRCDIR}/libs/aarch64-apple-darwin -lxdrjson -ldl -lm
#cgo linux,amd64   LDFLAGS: -L${SRCDIR}/libs/x86_64-unknown-linux-gnu -lxdrjson -ldl -lm
#cgo linux,arm64   LDFLAGS: -L${SRCDIR}/libs/aarch64-unknown-linux-gnu -lxdrjson -ldl -lm

*/
import "C"

import (
	"encoding/json"
	"unsafe"

	"github.com/pkg/errors"
)

// Decode XDR binary into JSON
//
// Takes in XDR binary, decodes it as the XdrType, returning XDR-JSON.
//
// If the XDR binary passed in is zero length, decoding will error.
//
// Returns the JSON message if decoding successful, otherwise an error.
func Decode(xdrTypeName XdrType, xdrBinary []byte) (json.RawMessage, error) {
	var jsonStr, errStr string
	// scope just added to show matching alloc/frees
	{
		goRawXdr := CXDR(xdrBinary)
		b := C.CString(string(xdrTypeName))

		result := C.xdr_to_json(b, goRawXdr)
		C.free(unsafe.Pointer(b))

		jsonStr = C.GoString(result.json)
		errStr = C.GoString(result.error)

		C.free_conversion_result(result)
	}

	if errStr != "" {
		return json.RawMessage(jsonStr), errors.New(errStr)
	}

	return json.RawMessage(jsonStr), nil
}

// CXDR is ripped directly from preflight.go to avoid a dependency.
func CXDR(xdr []byte) C.xdr_t {
	return C.xdr_t{
		xdr: (*C.uchar)(C.CBytes(xdr)),
		len: C.size_t(len(xdr)),
	}
}
