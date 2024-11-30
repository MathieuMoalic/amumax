package cuda

import (
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// multiply: dst[i] = a[i] * b[i]
// a and b must have the same number of components
func Mul(dst, a, b *data.Slice) {
	N := dst.Len()
	nComp := dst.NComp()
	// log.Assert(a.Len() == N && a.NComp() == nComp && b.Len() == N && b.NComp() == nComp)
	cfg := make1DConf(N)
	for c := 0; c < nComp; c++ {
		k_mul_async(dst.DevPtr(c), a.DevPtr(c), b.DevPtr(c), N, cfg)
	}
}

// divide: dst[i] = a[i] / b[i]
// divide-by-zero yields zero.
func Div(dst, a, b *data.Slice) {
	N := dst.Len()
	nComp := dst.NComp()
	log_old.AssertMsg(a.Len() == N && a.NComp() == nComp && b.Len() == N && b.NComp() == nComp, "Length or component mismatch in Mul")
	cfg := make1DConf(N)
	for c := 0; c < nComp; c++ {
		k_pointwise_div_async(dst.DevPtr(c), a.DevPtr(c), b.DevPtr(c), N, cfg)
	}
}

// Add: dst = src1 + src2.
func Add(dst, src1, src2 *data.Slice) {
	Madd2(dst, src1, src2, 1, 1)
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2
func Madd2(dst, src1, src2 *data.Slice, factor1, factor2 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	log_old.AssertMsg(src1.Len() == N && src2.Len() == N, "Length mismatch between src1 and src2 in Madd2")
	log_old.AssertMsg(src1.NComp() == nComp && src2.NComp() == nComp, "Component mismatch between src1 and src2 in Madd2")
	cfg := make1DConf(N)
	for c := 0; c < nComp; c++ {
		k_madd2_async(dst.DevPtr(c), src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2, N, cfg)
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3
func Madd3(dst, src1, src2, src3 *data.Slice, factor1, factor2, factor3 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	log_old.AssertMsg(src1.Len() == N && src2.Len() == N && src3.Len() == N, "Length mismatch between src1, src2, and src3 in Madd3")
	log_old.AssertMsg(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp, "Component mismatch between src1, src2, and src3 in Madd3")
	cfg := make1DConf(N)
	for c := 0; c < nComp; c++ {
		k_madd3_async(dst.DevPtr(c), src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2, src3.DevPtr(c), factor3, N, cfg)
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4
func Madd4(dst, src1, src2, src3, src4 *data.Slice, factor1, factor2, factor3, factor4 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	log_old.AssertMsg(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N, "Length mismatch between src1, src2, src3, and src4 in Madd4")
	log_old.AssertMsg(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp, "Component mismatch between src1, src2, src3, and src4 in Madd4")
	cfg := make1DConf(N)
	for c := 0; c < nComp; c++ {
		k_madd4_async(dst.DevPtr(c),
			src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2,
			src3.DevPtr(c), factor3,
			src4.DevPtr(c), factor4, N, cfg)
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4 + src5[i] * factor5
func Madd5(dst, src1, src2, src3, src4, src5 *data.Slice, factor1, factor2, factor3, factor4, factor5 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	log_old.AssertMsg(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N && src5.Len() == N, "Length mismatch between src1, src2, src3, src4, and src5 in Madd5")
	log_old.AssertMsg(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp && src5.NComp() == nComp, "Component mismatch between src1, src2, src3, src4, and src5 in Madd5")
	cfg := make1DConf(N)
	for c := 0; c < nComp; c++ {
		k_madd5_async(dst.DevPtr(c),
			src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2,
			src3.DevPtr(c), factor3,
			src4.DevPtr(c), factor4,
			src5.DevPtr(c), factor5, N, cfg)
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4 + src5[i] * factor5 + src6[i] * factor6
func Madd6(dst, src1, src2, src3, src4, src5, src6 *data.Slice, factor1, factor2, factor3, factor4, factor5, factor6 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	log_old.AssertMsg(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N && src5.Len() == N && src6.Len() == N, "Length mismatch between src1, src2, src3, src4, src5, and src6 in Madd6")
	log_old.AssertMsg(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp && src5.NComp() == nComp && src6.NComp() == nComp, "Component mismatch between src1, src2, src3, src4, src5, and src6 in Madd6")
	cfg := make1DConf(N)
	for c := 0; c < nComp; c++ {
		k_madd6_async(dst.DevPtr(c),
			src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2,
			src3.DevPtr(c), factor3,
			src4.DevPtr(c), factor4,
			src5.DevPtr(c), factor5,
			src6.DevPtr(c), factor6, N, cfg)
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4 + src5[i] * factor5 + src6[i] * factor6 + src7[i] * factor7
func Madd7(dst, src1, src2, src3, src4, src5, src6, src7 *data.Slice, factor1, factor2, factor3, factor4, factor5, factor6, factor7 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	log_old.AssertMsg(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N && src5.Len() == N && src6.Len() == N && src7.Len() == N, "Length mismatch between src1, src2, src3, src4, src5, src6, and src7 in Madd7")
	log_old.AssertMsg(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp && src5.NComp() == nComp && src6.NComp() == nComp && src7.NComp() == nComp, "Component mismatch between src1, src2, src3, src4, src5, src6, and src7 in Madd7")
	cfg := make1DConf(N)
	for c := 0; c < nComp; c++ {
		k_madd7_async(dst.DevPtr(c),
			src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2,
			src3.DevPtr(c), factor3,
			src4.DevPtr(c), factor4,
			src5.DevPtr(c), factor5,
			src6.DevPtr(c), factor6,
			src7.DevPtr(c), factor7, N, cfg)
	}
}
