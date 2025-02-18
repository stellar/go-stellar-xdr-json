package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/stellar/go-stellar-xdr-json/xdr2json"
)

type Asset struct{}

// A command-line tool to test that the xdr2json package compiles and runs on platforms:
//
//	go run github.com/stellar/go-stellar-xdr-json/xdr2json/tests/cmddecode AAAAAUFCQ0QAAAAAY3NrWbIzxOdwa+BhWFLSODTyh99xwpyb1gTzTQjDvGE=
func main() {
	xdrBase64 := "AAAAAUFCQ0QAAAAAY3NrWbIzxOdwa+BhWFLSODTyh99xwpyb1gTzTQjDvGE="
	wantJson := `{"credit_alphanum4":{"asset_code":"ABCD","issuer":"GBRXG22ZWIZ4JZ3QNPQGCWCS2I4DJ4UH35Y4FHE32YCPGTIIYO6GCL64"}}`

	fmt.Fprintf(os.Stderr, "xdr : %s\n", xdrBase64)
	fmt.Fprintf(os.Stderr, "want: %s\n", wantJson)

	rawBytes, err := base64.StdEncoding.DecodeString(xdrBase64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: decoding base64: %v\n", err)
		os.Exit(1)
	}

	json, err := xdr2json.ConvertBytes(Asset{}, rawBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: converting binary to json: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "got : %s\n", string(json))

	if string(json) != wantJson {
		fmt.Fprintf(os.Stderr, "FAIL\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "PASS\n")
}
