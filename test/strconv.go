package test

import "strconv"

type StrconvWrapper struct{}

func (s *StrconvWrapper) FormatInt(i int64, base int) string {
	return strconv.FormatInt(i, base)
}

func (s *StrconvWrapper) FormatBool(b bool) string {
	return strconv.FormatBool(b)
}

func (s *StrconvWrapper) FormatFloat(f float64, fmt byte, prec, bitSize int) string {
	return strconv.FormatFloat(f, fmt, prec, bitSize)
}

func (s *StrconvWrapper) FormatUint(i uint64, base int) string {
	return strconv.FormatUint(i, base)
}
