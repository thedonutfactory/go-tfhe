package tfhe

import (
	"fmt"
	"math/rand"
	"testing"
)

var (
	pub_key *PubKey
	pri_key *PriKey
)

func keys() (*PubKey, *PriKey) {
	if pri_key == nil {
		pub_key, pri_key = KeyGen()
	}
	return pub_key, pri_key
}

func setAllBit(num uint32) uint32 {
	var n uint32 = num
	n = n | n>>1
	n = n | n>>2
	n = n | n>>4
	n = n | n>>8
	n = n | n>>16
	n = n | n>>31
	return n
}

func xnor(x, y uint32) uint32 {
	var result uint32 = 0
	if x > y {
		result = (setAllBit(x) ^ x) ^ y
	} else {
		result = (setAllBit(y) ^ y) ^ x
	}
	return result
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func TestDecrypt1(t *testing.T) {
	_, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = 1
	Encrypt(ct[0], pt[0], pri_key)
	Decrypt(pt[1], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestDecrypt0(t *testing.T) {
	_, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = 0
	Encrypt(ct[0], pt[0], pri_key)
	Decrypt(pt[1], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestNand(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = in0.Message &^ in1.Message
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	Nand(ct[0], ct[0], ct[1], pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
	fmt.Println("done")
}

func TestOr(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = in0.Message | in1.Message
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	Or(ct[0], ct[0], ct[1], pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestMultiOr(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = in0.Message | (in0.Message | in1.Message)
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = 1
	pt[1].Message = 0
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	r1, r2 := NewCtxt(), NewCtxt()

	Or(r1, ct[0], ct[1], pub_key)
	Or(r2, ct[0], r1, pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], r2, pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestMultiOrAnd(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = in0.Message & (in0.Message | in1.Message)
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	r1, r2 := NewCtxt(), NewCtxt()
	And(r1, ct[0], ct[1], pub_key)
	Or(r2, ct[0], r1, pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], r2, pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestMultiAndOr(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = in0.Message | (in0.Message & in1.Message)
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = 1 //uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = 0 //uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	r1, r2 := NewCtxt(), NewCtxt()
	And(r1, ct[0], ct[1], pub_key)
	Or(r2, ct[0], r1, pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], r2, pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestAnd(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = in0.Message & in1.Message
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = 0 //uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = 1 //uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	And(ct[0], ct[0], ct[1], pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestMultiAnd(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = in0.Message & (in0.Message & in1.Message)
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = 0 //uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = 1 //uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	r1, r2 := NewCtxt(), NewCtxt()
	And(r1, ct[0], ct[1], pub_key)
	And(r2, ct[0], r1, pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], r2, pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestNor(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = 2 + ^(in0.Message | in1.Message)
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	Nor(ct[0], ct[0], ct[1], pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestXor(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = in0.Message ^ in1.Message
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = 1 //uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = 0 //uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	Xor(ct[0], ct[0], ct[1], pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestMultiXor(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = in0.Message ^ (in0.Message ^ in1.Message)
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = 1 //uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = 0 //uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	r1, r2 := NewCtxt(), NewCtxt()
	Xor(r1, ct[0], ct[1], pub_key)
	Xor(r2, ct[0], r1, pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], r2, pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestXnor(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = xnor(in0.Message, in1.Message)
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = 1 //uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = 0 //uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	Xnor(ct[0], ct[0], ct[1], pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestNot(t *testing.T) {
	Check := func(out, in0 *Ptxt) {
		out.Message = ^in0.Message
	}
	_, pri_key := keys()
	ct := NewCtxt()
	out, pt := NewPtxt(), NewPtxt()
	pt.Message = 1 //uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct, pt, pri_key)
	Not(ct, ct)
	Check(out, pt)
	Decrypt(pt, ct, pri_key)
	assertEqual(t, pt.Message, pt.Message, "plaintext messages not equal")
}

func TestCopy(t *testing.T) {
	Check := func(out, in0 *Ptxt) {
		out.Message = in0.Message
	}
	_, pri_key := keys()
	ct := NewCtxt()
	out, pt := NewPtxt(), NewPtxt()
	pt.Message = 1 //uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct, pt, pri_key)
	Copy(ct, ct)
	Check(out, pt)
	Decrypt(pt, ct, pri_key)
	assertEqual(t, pt.Message, pt.Message, "plaintext messages not equal")
}

func TestConstant0(t *testing.T) {
	Check := func(out *Ptxt) {
		out.Message = 0
	}
	_, pri_key := keys()
	ct := NewCtxt()
	out, pt := NewPtxt(), NewPtxt()
	pt.Message = 0 //uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct, pt, pri_key)
	Constant(ct, 0)
	Check(out)
	Decrypt(pt, ct, pri_key)
	assertEqual(t, pt.Message, pt.Message, "plaintext messages not equal")
}

func TestConstant1(t *testing.T) {
	Check := func(out *Ptxt) {
		out.Message = 1
	}
	_, pri_key := keys()
	ct := NewCtxt()
	out, pt := NewPtxt(), NewPtxt()
	pt.Message = 1 //uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct, pt, pri_key)
	Constant(ct, 1)
	Check(out)
	Decrypt(pt, ct, pri_key)
	assertEqual(t, pt.Message, pt.Message, "plaintext messages not equal")
}

func TestAndNY(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = ^in0.Message & in1.Message
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	AndNY(ct[0], ct[0], ct[1], pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestAndYN(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = in0.Message &^ in1.Message
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	AndYN(ct[0], ct[0], ct[1], pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestOrNY(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = 2 + (^in0.Message | in1.Message)
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	OrNY(ct[0], ct[0], ct[1], pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}

func TestOrYN(t *testing.T) {
	Check := func(out, in0, in1 *Ptxt) {
		out.Message = 2 + (in0.Message | ^in1.Message)
	}
	pub_key, pri_key := keys()
	pt := NewPtxtArray(2)
	ct := NewCtxtArray(2)
	pt[0].Message = uint32(rand.Int31() % KPtxtSpace)
	pt[1].Message = uint32(rand.Int31() % KPtxtSpace)
	Encrypt(ct[0], pt[0], pri_key)
	Encrypt(ct[1], pt[1], pri_key)
	OrYN(ct[0], ct[0], ct[1], pub_key)
	Check(pt[1], pt[0], pt[1])
	Decrypt(pt[0], ct[0], pri_key)
	assertEqual(t, pt[1].Message, pt[0].Message, "plaintext messages not equal")
}
