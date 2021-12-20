#include <stdio.h>
#include <cuda.h>
 

__global__ void vecmul(int32_t *A, int32_t* B, int32_t *C, int size)
{
    // Row and Column indexes: 
    int row = blockIdx.y*blockDim.y+threadIdx.y;
    int col = blockIdx.x*blockDim.x+threadIdx.x;

    // Are they bellow the maximum?
    if (col < size && row < size) {
       int32_t result = 0;
       for(int ix=0;ix<size;ix++) {
          result += A[row*size+ix]*B[ix*size+col];
       }
       C[row*size+col] = result;
    }
}

extern "C" {

    void maxmul(int32_t *A, int32_t* B, int32_t *C, int size) {

        int total = size*size;

        // Allocate device memory:
        int32_t* gpu_A;
        int32_t* gpu_B;
        int32_t* gpu_C;
        int msize = total * sizeof(int32_t);
        cudaMalloc((void**)&gpu_A, msize);
        cudaMemcpy(gpu_A,A,msize,cudaMemcpyHostToDevice);
        cudaMalloc((void**)&gpu_B, msize);
        cudaMemcpy(gpu_B,B,msize,cudaMemcpyHostToDevice);
        cudaMalloc((void**)&gpu_C,msize);

        // Blocks & grids:
        dim3 blocks(size,size);
        dim3 grid(1,1);

        // Call the kernel:
        vecmul<<<grid,blocks>>>(gpu_A,gpu_B,gpu_C,size);

        // Get the result Matrix:
        cudaMemcpy(C,gpu_C,msize,cudaMemcpyDeviceToHost);

        //Free device matrices
        cudaFree(gpu_A);
        cudaFree(gpu_B);
        cudaFree(gpu_C);
    }

}