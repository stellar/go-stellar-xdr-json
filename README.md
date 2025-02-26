# go-stellar-xdr-json
A Golang library for XDR conversion to and from JSON.
This repository combines Rust code with Go bindings to generate a package that converts XDR to JSON.

## Building the Project
To compile the archive file (.a), run:
```azure
make build-libs
```

To delete all libraries and build files, run:
```azure
make dist-clean
```

## Running Unit Tests
Once you have built the libraries locally, run the unit tests with:
```azure
go test ./...
```


**TODO:**

* [Automate building/updates for xdrjson](https://github.com/stellar/stellar-rpc/issues/350)
