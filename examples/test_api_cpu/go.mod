module github.com/thedonutfactory/examples

go 1.17

require github.com/thedonutfactory/go-tfhe v0.0.0

replace github.com/thedonutfactory/go-tfhe v0.0.0 => ../..

require (
	golang.org/x/exp v0.0.0-20191002040644-a1355ae1e2c3 // indirect
	gonum.org/v1/gonum v0.9.3 // indirect
)
