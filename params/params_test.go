package params_test

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/params"
)

func TestSecurityLevelSwitching(t *testing.T) {
	// Test 128-bit (default)
	params.CurrentSecurityLevel = params.Security128Bit
	tlwe0 := params.GetTLWELv0()
	if tlwe0.N != 700 {
		t.Errorf("128-bit TLWE Lv0 N = %d, expected 700", tlwe0.N)
	}

	// Test 110-bit
	params.CurrentSecurityLevel = params.Security110Bit
	tlwe0 = params.GetTLWELv0()
	if tlwe0.N != 630 {
		t.Errorf("110-bit TLWE Lv0 N = %d, expected 630", tlwe0.N)
	}

	// Test 80-bit
	params.CurrentSecurityLevel = params.Security80Bit
	tlwe0 = params.GetTLWELv0()
	if tlwe0.N != 550 {
		t.Errorf("80-bit TLWE Lv0 N = %d, expected 550", tlwe0.N)
	}

	// Reset to default
	params.CurrentSecurityLevel = params.Security128Bit
}

func TestParameterConsistency(t *testing.T) {
	params.CurrentSecurityLevel = params.Security128Bit

	tlwe0 := params.GetTLWELv0()
	tlwe1 := params.GetTLWELv1()
	trlwe1 := params.GetTRLWELv1()
	trgsw1 := params.GetTRGSWLv1()

	// Verify basic constraints
	if tlwe0.N <= 0 {
		t.Errorf("TLWE Lv0 N must be positive, got %d", tlwe0.N)
	}
	if tlwe1.N <= 0 {
		t.Errorf("TLWE Lv1 N must be positive, got %d", tlwe1.N)
	}
	if tlwe0.ALPHA <= 0 {
		t.Errorf("TLWE Lv0 ALPHA must be positive, got %f", tlwe0.ALPHA)
	}
	if tlwe1.ALPHA <= 0 {
		t.Errorf("TLWE Lv1 ALPHA must be positive, got %f", tlwe1.ALPHA)
	}

	// Verify TRLWE and TLWE Lv1 have same N
	if trlwe1.N != tlwe1.N {
		t.Errorf("TRLWE Lv1 N (%d) should equal TLWE Lv1 N (%d)", trlwe1.N, tlwe1.N)
	}

	// Verify TRGSW BG matches BGBIT
	expectedBG := uint32(1) << trgsw1.BGBIT
	if trgsw1.BG != expectedBG {
		t.Errorf("TRGSW BG = %d, expected %d (1 << %d)", trgsw1.BG, expectedBG, trgsw1.BGBIT)
	}

	// Verify TRGSW N matches TLWE Lv1 N
	if trgsw1.N != tlwe1.N {
		t.Errorf("TRGSW N (%d) should equal TLWE Lv1 N (%d)", trgsw1.N, tlwe1.N)
	}
}

func TestSecurityInfo(t *testing.T) {
	info := params.SecurityInfo()
	if info == "" {
		t.Error("SecurityInfo returned empty string")
	}
	t.Logf("Security info: %s", info)
}

func TestKSKAndBSKAlpha(t *testing.T) {
	params.CurrentSecurityLevel = params.Security128Bit

	kskAlpha := params.KSKAlpha()
	bskAlpha := params.BSKAlpha()

	if kskAlpha <= 0 {
		t.Errorf("KSKAlpha must be positive, got %f", kskAlpha)
	}
	if bskAlpha <= 0 {
		t.Errorf("BSKAlpha must be positive, got %f", bskAlpha)
	}

	// KSK uses Lv0 alpha, BSK uses Lv1 alpha
	if kskAlpha != params.GetTLWELv0().ALPHA {
		t.Errorf("KSKAlpha (%f) should equal TLWE Lv0 ALPHA (%f)", kskAlpha, params.GetTLWELv0().ALPHA)
	}
	if bskAlpha != params.GetTLWELv1().ALPHA {
		t.Errorf("BSKAlpha (%f) should equal TLWE Lv1 ALPHA (%f)", bskAlpha, params.GetTLWELv1().ALPHA)
	}
}
