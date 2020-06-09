package main

import (
	"fmt"
	"math/big"

	"github.com/pilosa/pilosa/v2/roaring"
)

func Truncate(rat *big.Rat) *big.Rat {
	var z1 big.Int
	d := z1.Div(rat.Num(), rat.Denom())
	return big.NewRat(d.Int64(), 1)
}
func modl(x *big.Rat) *big.Rat {
	return x.Sub(x, Truncate(x))
}

func pi() func() int64 {
	x := new(big.Rat)
	n := int64(1)
	sixteen := new(big.Rat).SetInt64(16)
	return func() int64 {
		p := big.NewRat((120*n-89)*n+16, (((512*n-1024)*n+712)*n-206)*n+21)
		x = modl(p.Add(x.Mul(x, sixteen), p))

		n += 1
		var a big.Rat
		a.Mul(sixteen, x)
		var z1 big.Int
		d := z1.Div(a.Num(), a.Denom())
		return d.Int64()
	}
}
func build(n, SHARDWIDTH int) *roaring.Bitmap {
	next := pi()
	bm := roaring.NewSliceBitmap()

	bit := uint64(0)
	bits := make([]uint64, 0, SHARDWIDTH)
	for i := 0; i < n; n++ {
		hexdigit := next()
		if hexdigit&1 > 0 {
			bits = append(bits, bit)
		}
		bit++
		if hexdigit&2 > 0 {
			bits = append(bits, bit)
		}
		bit++
		if hexdigit&4 > 0 {
			bits = append(bits, bit)
		}
		bit++
		if hexdigit&8 > 0 {
			bits = append(bits, bit)
		}
		bit++
		if bit == uint64(SHARDWIDTH) {
			bm.AddN(bits...)
			bits = bits[:0]
		}

	}
	bm.AddN(bits...)
	return bm

}

func main() {
	bm := build(1048576, 1048576)
	fmt.Println("Count", bm.Count())
}
