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
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go github.com/stellar/go-stellar-xdr-json/cmddecode <asset_base64_string>\n")
		os.Exit(1)
	}

	xdrBase64 := os.Args[1]

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

	fmt.Println(string(json))
}
