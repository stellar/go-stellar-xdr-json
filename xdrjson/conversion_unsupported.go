//go:build !(cgo && ((windows && amd64) || (darwin && amd64) || (darwin && arm64) || (linux && amd64) || (linux && arm64)))

package xdrjson

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Decode is the unsupported-platform stub. The real cgo-backed implementation
// lives in conversion.go and is only compiled when cgo is enabled on a
// supported platform. On all other configurations this stub is compiled
// instead so that the package still builds, and Decode returns an error.
func Decode(xdrTypeName XdrType, xdrBinary []byte) (json.RawMessage, error) {
	return nil, errors.New("xdrjson.Decode requires cgo on a supported platform " +
		"(windows/amd64, darwin/amd64, darwin/arm64, linux/amd64, or linux/arm64)")
}
