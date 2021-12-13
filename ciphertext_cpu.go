package tfhe

func NewCtxt() *Ctxt {
	param := GetDefaultParam()
	return &Ctxt{
		lwe_sample: NewLWESample(param.lwe_n),
	}
}

func NewCtxtArray(n int) (arr []*Ctxt) {
	arr = make([]*Ctxt, n)
	for i := 0; i < n; i++ {
		arr[i] = NewCtxt()
	}
	return
}
