// Rust FFT Bridge for Go-TFHE
// Exposes C-compatible FFT functions that Go can call via CGO
// Uses the same ExtendedFFT algorithm as rs-tfhe

use rustfft::num_complex::Complex;
use rustfft::Fft;
use std::f64::consts::PI;
use std::sync::Arc;

// Opaque handle for the FFT processor
pub struct FFTProcessor {
  // Pre-computed twisting factors (2N-th roots of unity)
  twisties_re: Vec<f64>,
  twisties_im: Vec<f64>,
  // rustfft's optimized N/2-point FFT (512 for N=1024)
  fft_n2_fwd: Arc<dyn Fft<f64>>,
  fft_n2_inv: Arc<dyn Fft<f64>>,
  // Pre-allocated buffers
  fourier_buffer: Vec<Complex<f64>>,
  scratch_fwd: Vec<Complex<f64>>,
  scratch_inv: Vec<Complex<f64>>,
}

/// Create a new FFT processor for 1024-point transforms
/// Returns an opaque pointer to be used in subsequent calls
#[no_mangle]
pub extern "C" fn fft_processor_new() -> *mut FFTProcessor {
  const N: usize = 1024;
  const N2: usize = N / 2; // 512

  // Compute twisting factors: exp(i*π*k/N) for k=0..N/2-1
  let mut twisties_re = Vec::with_capacity(N2);
  let mut twisties_im = Vec::with_capacity(N2);
  let twist_unit = PI / (N as f64);
  for i in 0..N2 {
    let angle = i as f64 * twist_unit;
    let (im, re) = angle.sin_cos();
    twisties_re.push(re);
    twisties_im.push(im);
  }

  // Use rustfft's planner - auto-detects NEON (ARM), AVX (x86), or scalar
  use rustfft::FftPlanner;
  let mut planner = FftPlanner::new();
  let fft_n2_fwd = planner.plan_fft_forward(N2);
  let fft_n2_inv = planner.plan_fft_inverse(N2);

  // Pre-allocate scratch buffers
  let scratch_fwd_len = fft_n2_fwd.get_inplace_scratch_len();
  let scratch_inv_len = fft_n2_inv.get_inplace_scratch_len();

  let processor = Box::new(FFTProcessor {
    twisties_re,
    twisties_im,
    fft_n2_fwd,
    fft_n2_inv,
    fourier_buffer: vec![Complex::new(0.0, 0.0); N2],
    scratch_fwd: vec![Complex::new(0.0, 0.0); scratch_fwd_len],
    scratch_inv: vec![Complex::new(0.0, 0.0); scratch_inv_len],
  });

  Box::into_raw(processor)
}

/// Free the FFT processor
#[no_mangle]
pub extern "C" fn fft_processor_free(processor: *mut FFTProcessor) {
  if !processor.is_null() {
    unsafe {
      let _ = Box::from_raw(processor);
    }
  }
}

/// IFFT: Convert Torus32 to frequency domain (f64)
/// Matches rs-tfhe: ifft_1024(&[u32; 1024]) -> [f64; 1024]
/// Input: torus_in[1024] - Torus32 polynomial
/// Output: freq_out[1024] - frequency domain (split: re[0..512], im[0..512])
#[no_mangle]
pub extern "C" fn ifft_1024_negacyclic(
  processor: *mut FFTProcessor,
  torus_in: *const u32,
  freq_out: *mut f64,
) {
  if processor.is_null() || torus_in.is_null() || freq_out.is_null() {
    return;
  }

  unsafe {
    let proc = &mut *processor;
    let input = std::slice::from_raw_parts(torus_in, 1024);
    let output = std::slice::from_raw_parts_mut(freq_out, 1024);

    const N: usize = 1024;
    const N2: usize = N / 2; // 512

    let (input_re, input_im) = input.split_at(N2);

    // Apply twisting factors and convert (same as rs-tfhe)
    for i in 0..N2 {
      let in_re = input_re[i] as i32 as f64;
      let in_im = input_im[i] as i32 as f64;
      let w_re = proc.twisties_re[i];
      let w_im = proc.twisties_im[i];
      proc.fourier_buffer[i] =
        Complex::new(in_re * w_re - in_im * w_im, in_re * w_im + in_im * w_re);
    }

    // 512-point FORWARD FFT with scratch buffer
    proc
      .fft_n2_fwd
      .process_with_scratch(&mut proc.fourier_buffer, &mut proc.scratch_fwd);

    // Scale by 2 and convert to output (same as rs-tfhe)
    for i in 0..N2 {
      output[i] = proc.fourier_buffer[i].re * 2.0;
      output[i + N2] = proc.fourier_buffer[i].im * 2.0;
    }
  }
}

/// FFT: Convert frequency domain (f64) to Torus32
/// Matches rs-tfhe: fft_1024(&[f64; 1024]) -> [u32; 1024]
/// Input: freq_in[1024] - frequency domain (split: re[0..512], im[0..512])
/// Output: torus_out[1024] - Torus32 polynomial
#[no_mangle]
pub extern "C" fn fft_1024_negacyclic(
  processor: *mut FFTProcessor,
  freq_in: *const f64,
  torus_out: *mut u32,
) {
  if processor.is_null() || freq_in.is_null() || torus_out.is_null() {
    return;
  }

  unsafe {
    let proc = &mut *processor;
    let input = std::slice::from_raw_parts(freq_in, 1024);
    let output = std::slice::from_raw_parts_mut(torus_out, 1024);

    const N: usize = 1024;
    const N2: usize = N / 2; // 512

    // Convert to complex and scale (same as rs-tfhe)
    let (input_re, input_im) = input.split_at(N2);
    for i in 0..N2 {
      proc.fourier_buffer[i] = Complex::new(input_re[i] * 0.5, input_im[i] * 0.5);
    }

    // 512-point INVERSE FFT with scratch buffer
    proc
      .fft_n2_inv
      .process_with_scratch(&mut proc.fourier_buffer, &mut proc.scratch_inv);

    // Apply inverse twisting and convert to u32 (same as rs-tfhe)
    let normalization = 1.0 / (N2 as f64);
    for i in 0..N2 {
      let w_re = proc.twisties_re[i];
      let w_im = proc.twisties_im[i];
      let f_re = proc.fourier_buffer[i].re;
      let f_im = proc.fourier_buffer[i].im;
      let tmp_re = (f_re * w_re + f_im * w_im) * normalization;
      let tmp_im = (f_im * w_re - f_re * w_im) * normalization;
      output[i] = tmp_re.round() as i64 as u32;
      output[i + N2] = tmp_im.round() as i64 as u32;
    }
  }
}

/// Batch IFFT for multiple polynomials (used in blind rotation)
/// Input: torus_in (count * 1024 u32 values)
/// Output: freq_out (count * 1024 f64 values)
#[no_mangle]
pub extern "C" fn batch_ifft_1024_negacyclic(
  processor: *mut FFTProcessor,
  torus_in: *const u32,
  freq_out: *mut f64,
  count: usize,
) {
  if processor.is_null() || torus_in.is_null() || freq_out.is_null() {
    return;
  }

  unsafe {
    let input = std::slice::from_raw_parts(torus_in, count * 1024);
    let output = std::slice::from_raw_parts_mut(freq_out, count * 1024);

    // Process each polynomial
    for i in 0..count {
      let offset = i * 1024;
      ifft_1024_negacyclic(
        processor,
        input[offset..].as_ptr(),
        output[offset..].as_mut_ptr(),
      );
    }
  }
}

#[cfg(test)]
mod tests {
  use super::*;

  #[test]
  fn test_fft_roundtrip() {
    let processor = fft_processor_new();

    let mut input = [0u32; 1024];
    input[0] = 1000;
    input[100] = 500;
    input[512] = 300; // imaginary part

    let mut freq = [0.0f64; 1024];
    let mut output = [0u32; 1024];

    unsafe {
      // First cast u32 array to f64 for FFT input
      let input_f64: [f64; 1024] = std::array::from_fn(|i| input[i] as i32 as f64);
      fft_1024_negacyclic(processor, input_f64.as_ptr(), freq.as_mut_ptr() as *mut u32);
      ifft_1024_negacyclic(processor, freq.as_ptr(), output.as_mut_ptr());
    }

    // Check roundtrip (allowing for rounding errors)
    let mut max_error = 0i64;
    for i in 0..1024 {
      let diff = (input[i] as i64 - output[i] as i64).abs();
      max_error = max_error.max(diff);
    }

    assert!(
      max_error <= 2,
      "Max roundtrip error: {} (should be ≤ 2)",
      max_error
    );

    fft_processor_free(processor);
  }

  #[test]
  fn test_processor_lifecycle() {
    // Test that we can create and free multiple processors
    for _ in 0..10 {
      let proc = fft_processor_new();
      assert!(!proc.is_null());
      fft_processor_free(proc);
    }
  }
}
