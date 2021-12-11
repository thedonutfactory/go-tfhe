#include <device_launch_parameters.h>
#include <cuda_runtime.h>
#include <cufft.h>
#include <vector>
#include <cmath>

#include <stdio.h>

const auto BATCH = 1;

__global__ void ComplexPointwiseMulAndScale(cufftComplex *a, cufftComplex *b, int size)
{
    const int numThreads = blockDim.x * gridDim.x;
    const int threadID = blockIdx.x * blockDim.x + threadIdx.x;
    float scale = 1.0f / (float)size;
    cufftComplex c;
    for (int i = threadID; i < size; i += numThreads)
    {
        c = cuCmulf(a[i], b[i]);
        b[i] = make_cuFloatComplex(scale*cuCrealf(c), scale*cuCimagf(c));
    }
}

__global__ void ConvertToInt(cufftReal *a, int size)
{
    const int numThreads = blockDim.x * gridDim.x;
    const int threadID = blockIdx.x * blockDim.x + threadIdx.x;
    auto b = (int*)a;
    for (int i = threadID; i < size; i += numThreads)
        b[i] = static_cast<int>(round(a[i]));
}

std::vector<int> multiply(const std::vector<float> &a, const std::vector<float> &b)
{
    const auto NX = a.size();
    cufftHandle plan_a, plan_b, plan_c;
    cufftComplex *data_a, *data_b;
    std::vector<int> c(a.size() + 1);
    c[0] = 0;

    //Allocate graphics card memory and initialize, assuming sizeof(int)==sizeof(float), sizeof(cufftComplex)==2*sizeof(float)
    cudaMalloc((void**)&data_a, sizeof(cufftComplex) * (NX / 2 + 1) * BATCH);
    cudaMalloc((void**)&data_b, sizeof(cufftComplex) * (NX / 2 + 1) * BATCH);
    cudaMemcpy(data_a, a.data(), sizeof(float) * a.size(), cudaMemcpyHostToDevice);
    cudaMemcpy(data_b, b.data(), sizeof(float) * b.size(), cudaMemcpyHostToDevice);
    if (cudaGetLastError() != cudaSuccess) { fprintf(stderr, "Cuda error: Failed to allocate\n"); return c; }

    if (cufftPlan1d(&plan_a, NX, CUFFT_R2C, BATCH) != CUFFT_SUCCESS) { fprintf(stderr, "CUFFT error: Plan creation failed"); return c; }
    if (cufftPlan1d(&plan_b, NX, CUFFT_R2C, BATCH) != CUFFT_SUCCESS) { fprintf(stderr, "CUFFT error: Plan creation failed"); return c; }
    if (cufftPlan1d(&plan_c, NX, CUFFT_C2R, BATCH) != CUFFT_SUCCESS) { fprintf(stderr, "CUFFT error: Plan creation failed"); return c; }

    //Converting A(x) to Frequency Domain
    if (cufftExecR2C(plan_a, (cufftReal*)data_a, data_a) != CUFFT_SUCCESS)
    {
        fprintf(stderr, "CUFFT error: ExecR2C Forward failed");
        return c;
    }

    //Converting B(x) to Frequency Domain
    if (cufftExecR2C(plan_b, (cufftReal*)data_b, data_b) != CUFFT_SUCCESS)
    {
        fprintf(stderr, "CUFFT error: ExecR2C Forward failed");
        return c;
    }

    //Point multiplication
    ComplexPointwiseMulAndScale<<<NX / 256 + 1, 256>>>(data_a, data_b, NX);

    //Converting C(x) back to time domain
    if (cufftExecC2R(plan_c, data_b, (cufftReal*)data_b) != CUFFT_SUCCESS)
    {
        fprintf(stderr, "CUFFT error: ExecC2R Forward failed");
        return c;
    }

    //Converting the results of floating-point numbers to integers
    ConvertToInt<<<NX / 256 + 1, 256>>>((cufftReal*)data_b, NX);

    if (cudaDeviceSynchronize() != cudaSuccess) 
    {
        fprintf(stderr, "Cuda error: Failed to synchronize\n");
        return c;
    }

    cudaMemcpy(&c[1], data_b, sizeof(float) * b.size(), cudaMemcpyDeviceToHost);

    cufftDestroy(plan_a);
    cufftDestroy(plan_b);
    cufftDestroy(plan_c);
    cudaFree(data_a);
    cudaFree(data_b);
    return c;
}


int main(int argc, char **argv) 
{
    //Set base
    const auto base = 10;

    //999 * 9
    std::vector<float> a{ 0, 9, 9, 9 }; 
    std::vector<float> b{ 0, 0, 0, 9 };

    auto c = multiply(a, b);

    for (auto i : c)
        printf("%d ", i);
    printf("\n");

    //Processing carry
    for (int i = c.size() - 1; i > 0; i--)
    {
        if (c[i] >= base)
        {
            c[i - 1] += c[i] / base;
            c[i] %= base;
        }
    }

    //Remove excess zeros
    c.pop_back();
    auto i = 0;
    if (c[0] == 0)
        i++;

    //To output the final result, we need to change the mode of output, such as the decimal system is "% 2d" and the decimal system is "% 3d".
    for (; i < c.size(); i++)
        printf("%d", c[i]);
    printf("\n");

    return 0;
}