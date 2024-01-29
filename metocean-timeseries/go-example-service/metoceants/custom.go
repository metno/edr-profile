package metoceants

//go:generate /Users/vegardb/go/bin/oapi-codegen -o stub.go -generate server,types -package stub ../../openapi/edr.yaml

// N0 is needed since oapi-codegen creates code that references it, but does not create a type for it.
type N0 [4]float64

// N1 is needed since oapi-codegen creates code that references it, but does not create a type for it.
type N1 [6]float64
