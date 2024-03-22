package main

import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {
	//(-1, -1) + (-1, 1)
	p := ecc.NewEllipticCurvePoint(big.NewInt(int64(-1)), big.NewInt(int64(-1)),
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	p2 := ecc.NewEllipticCurvePoint(big.NewInt(int64(-1)), big.NewInt(int64(1)),
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	res := p.Add(p2)
	fmt.Printf("result of adding points on vertical line: %s\n", res)

	//C = A(2,5) + B(-1,-1)
	A := ecc.NewEllipticCurvePoint(big.NewInt(int64(2)), big.NewInt(int64(5)),
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	B := ecc.NewEllipticCurvePoint(big.NewInt(int64(-1)), big.NewInt(int64(-1)),
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	C := A.Add(B)
	fmt.Printf("A(2,5)+B(-1,-1) = %s\n", C)

	//C=B+B
	C = B.Add(B)
	fmt.Printf("B(-1,-1) + B(-1,-1)=%s\n", C)
}
