#!/bin/bash

nvcc -I/usr/local/cuda/include/ -L/usr/local/cuda/lib64/ -lcufft fft.cu -o fft
./fft