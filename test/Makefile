.PHONY: run
run: kernel
	./kernel

install:
	cp libmultkernel.so /usr/local/lib/
	@echo you can export /usr/local/lib to your LD_LIBRARY_PATH: 
	@echo \\t "export LD_LIBRARY_PATH=\$$LD_LIBRARY_PATH:/usr/local/lib"
	@echo or you can call your program like this:
	@echo \\t "LD_LIBRARY_PATH=/usr/local/lib your-program"
	@echo \\t "LD_LIBRARY_PATH=/path/to/lib your-program"

uninstall:
	rm -f /usr/local/lib/libmultkernel.so

clean:
	@-rm kernel kernel-test maxmul maxmul2 fft *.so

kernel:
	nvcc kernel.cu -o kernel-test
	nvcc --ptxas-options=-v --compiler-options '-fPIC' -o libmultkernel.so --shared kernel.cu
	go build -o kernel kernel.go

maxmul:
	nvcc --ptxas-options=-v --compiler-options '-fPIC' -o libmaxmul.so --shared maxmul.cu
	go build -o maxmul maxmul.go
	echo "run with: LD_LIBRARY_PATH=\${PWD} ./maxmul"
	LD_LIBRARY_PATH=${PWD} ./maxmul

maxmul2:
	nvcc --ptxas-options=-v --compiler-options '-fPIC' -o libmaxmul2.so --shared maxmul2.cu
	go build -o maxmul2 maxmul2.go
	echo "run with: LD_LIBRARY_PATH=\${PWD} ./maxmul2"
	LD_LIBRARY_PATH=${PWD} ./maxmul2

fft:
	nvcc -I/usr/local/cuda/include/ -L/usr/local/cuda/lib64/ -lcufft fft.cu -o fft
	./fft

test:
	go test -timeout 30s -run ^TestTorusPolynomialMultCuda$ github.com/thedonutfactory/go-tfhe
