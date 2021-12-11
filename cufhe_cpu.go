package tfhe

func NewCtxt(isAlias bool) *Ctxt {
	param := GetDefaultParam()
	return &Ctxt{
		lwe_sample_: NewLWESample(param.lwe_n_),
	}
}

func NewCtxtArray(n int) (arr []*Ctxt) {
	arr = make([]*Ctxt, n)
	for i := 0; i < n; i++ {
		arr[i] = NewCtxt(false)
	}
	return
}
