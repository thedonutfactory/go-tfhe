#!/bin/bash

nvcc --ptxas-options=-v --compiler-options '-fPIC' -o libmaxmul2.so --shared maxmul2.cu
go build -o maxmul2 maxmul2.go

echo "run with: LD_LIBRARY_PATH=\${PWD} ./maxmul2"
LD_LIBRARY_PATH=${PWD} ./maxmul2