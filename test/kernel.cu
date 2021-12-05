
#include "cuda_runtime.h"
#include "device_launch_parameters.h"

#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>

//void multiplyWithCuda(long *c, const long *a, const long *b, unsigned int size);
//void multiplyWithCudaKaratsuba(long *c, const long *a, const long *b, unsigned int size);

__device__ void multiplyKernelKaratsubaRec(long *z, const long *x, const long *y, unsigned int size)
{
	const long *a, *b, *c, *d;
	long *ab, *ac;
	long *bd, *cd;
	long *adbc;


	if (size <= 1)
	{
		z[0] = x[0] * y[0];
	}
	else
	{
		int half = (int)size / 2;

		ab = (long*)malloc(half * sizeof(long));
		ac = (long*)malloc(half * sizeof(long));
		cd = (long*)malloc(half * sizeof(long));
		bd = (long*)malloc(half * sizeof(long));
		adbc = (long*)malloc(half * sizeof(long));

		a = x;
		b = x + half;

		c = y;
		d = y + half;

		multiplyKernelKaratsubaRec(ac, a, c, half);
		multiplyKernelKaratsubaRec(bd, b, d, size - half);

		int i = 0;
		for (i = 0; i < half; i++)
		{
			ab[i] = a[i] + b[i];
			cd[i] = c[i] + d[i];
		}

		multiplyKernelKaratsubaRec(adbc, ab, cd, half);

		for (i = 0; i < half; i++)
		{
			z[i] = adbc[i] - ac[i] - bd[i];
		}
	}
}

__global__ void multiplyKernelKaratsuba(long *z, const long *x, const long *y, unsigned int size)
{
	multiplyKernelKaratsubaRec(z, x, y, size);
}

__global__ void multiplyKernel(long *c, const long *a, const long *b, unsigned int size)
{
    int i = threadIdx.x;
	c[i] = 0;
	for (auto x = 0; x < size; x++)
    {
	    for (auto y = 0; y < size; y++)
	    {
		    if (x + y == i)
		    {
				c[i] += a[x] * b[y];
		    }
	    }
    }
}

extern "C" {
	void multiplyWithCudaKaratsuba(long *c, const long *a, const long *b, unsigned int size)
	{
		long *dev_a = nullptr;
		long *dev_b = nullptr;
		long *dev_c = nullptr;

		cudaSetDevice(0);

		cudaMalloc(&dev_c, 2 * size * sizeof(long));
		cudaMalloc(&dev_a, size * sizeof(long));
		cudaMalloc(&dev_b, size * sizeof(long));

		cudaMemcpy(dev_a, a, size * sizeof(long), cudaMemcpyHostToDevice);
		cudaMemcpy(dev_b, b, size * sizeof(long), cudaMemcpyHostToDevice);

		int thread_num = 2 * size;
		multiplyKernelKaratsuba <<<1, thread_num >>> (dev_c, dev_a, dev_b, size);

		cudaDeviceSynchronize();

		cudaMemcpy(c, dev_c, 2 * size * sizeof(long), cudaMemcpyDeviceToHost);

		cudaFree(dev_c);
		cudaFree(dev_a);
		cudaFree(dev_b);
	}

	void multiplyWithCuda(long *c, const long *a, const long *b, unsigned int size)
	{
		long *dev_a = nullptr;
		long *dev_b = nullptr;
		long *dev_c = nullptr;

		cudaSetDevice(0);

		cudaMalloc(&dev_c, 2 * size * sizeof(long));
		cudaMalloc(&dev_a, size * sizeof(long));
		cudaMalloc(&dev_b, size * sizeof(long));
		
		cudaMemcpy(dev_a, a, size * sizeof(long), cudaMemcpyHostToDevice);
		cudaMemcpy(dev_b, b, size * sizeof(long), cudaMemcpyHostToDevice);

		int thread_num = 2 * size;
		multiplyKernel<<<1, thread_num>>>(dev_c, dev_a, dev_b, size);
		
		cudaDeviceSynchronize();

		cudaMemcpy(c, dev_c, 2 * size * sizeof(long), cudaMemcpyDeviceToHost);

		cudaFree(dev_c);
		cudaFree(dev_a);
		cudaFree(dev_b);

		printf("first val: %ld\n", c[0]);
	}

}

int main()
{
	const auto arraySize = 9;
	//long a[arraySize];
	//long b[arraySize];
	long c[2 * arraySize];

	long a[] = { 0, 9, 9, 9, 0, 9, 9, 9 };
	long b[] = { 0, 0, 0, 0, 0, 9, 9, 9 };

	for (auto i = 0; i < arraySize; i++) {
		//a[i] = rand() % 100;
		//b[i] = rand() % 100;
		c[i] = c[arraySize + i] = 0;
	}

    // Multiply polynomials in parallel.
	multiplyWithCuda(c, a, b, arraySize);
	
	for (auto i = 0; i < arraySize; i++) {
		printf("%ld, ", c[i]);
	}
	printf("\n");
}

int main2()
{
	srand(time(nullptr));
	
	const auto arraySize = 1024;
	long a[arraySize];
	long b[arraySize];
	long c[2 * arraySize];

	for (auto i = 0; i < arraySize; i++)
	{
		a[i] = rand() % 100;
		b[i] = rand() % 100;
		c[i] = c[arraySize + i] = 0;
	}

    // Multiply polynomials in parallel.
	time_t timeStart;
	time_t timeEnd;
	time(&timeStart);
	for (auto i = 0; i < 10000; i++)
		multiplyWithCuda(c, a, b, arraySize);
	time(&timeEnd);

	printf("time taken (normal) : %ld (%ld : %ld) \n", timeEnd - timeStart, timeStart, timeEnd);

	time(&timeStart);
	for (auto i = 0; i < 10000; i++)
		multiplyWithCudaKaratsuba(c, a, b, arraySize);
	time(&timeEnd);

	printf("time taken (karatsuba) : %ld (%ld : %ld) \n", timeEnd - timeStart, timeStart, timeEnd);


//	for (auto i = 0; i < arraySize; i++)
//	{
//		printf("%d ", a[i]);
//	}
//	printf("\n");
//
//	for (auto i = 0; i < arraySize; i++)
//	{
//		printf("%d ", b[i]);
//	}
//	printf("\n");
//
//    for (auto i = 0; i < 2 * arraySize; i++)
//    {
//		printf("%d ", c[i]);
//    }
//	printf("\n");

    cudaDeviceReset();

    return 0;
}
