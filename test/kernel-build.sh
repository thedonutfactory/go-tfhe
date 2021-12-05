#!/bin/bash

nvcc kernel.cu -o kernel-test

nvcc --ptxas-options=-v --compiler-options '-fPIC' -o libmultkernel.so --shared kernel.cu
go build -o kernel kernel.go

echo "run with: LD_LIBRARY_PATH=\${PWD} ./kernel"
LD_LIBRARY_PATH=${PWD} ./kernel
