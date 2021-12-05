#!/bin/bash

nvcc --ptxas-options=-v --compiler-options '-fPIC' -o libmaxmul.so --shared maxmul.cu
go build -o maxmul maxmul.go

echo "run with: LD_LIBRARY_PATH=\${PWD} ./maxmul"
LD_LIBRARY_PATH=${PWD} ./maxmul