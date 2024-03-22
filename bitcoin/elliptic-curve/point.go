package elliptic_curve

import (
	"fmt"
	"math/big"
)

type OP_TYPE int

const (
	ADD OP_TYPE = iota
	SUB
	MUL
	DIV
	EXP
)

type Point struct {
	//coefficients of curve
	a *big.Int
	b *big.Int
	//x, y should be the point on the curve
	x *big.Int
	y *big.Int
}

func OpOnBig(x *big.Int, y *big.Int, opType OP_TYPE) *big.Int {
	var op big.Int
	switch opType {
	case ADD:
		return op.Add(x, y)
	case SUB:
		return op.Sub(x, y)
	case MUL:
		return op.Mul(x, y)
	case DIV:
		return op.Div(x, y)
	case EXP:
		return op.Exp(x, y, nil)
	}

	panic("should not come to here")
}

func NewEllipticCurvePoint(x *big.Int, y *big.Int, a *big.Int, b *big.Int) *Point {
	if x == nil && y == nil {
		return &Point{
			x: x,
			y: y,
			a: a,
			b: b,
		}
	}

	left := OpOnBig(y, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(x, big.NewInt(int64(3)), EXP)
	ax := OpOnBig(a, x, MUL)
	right := OpOnBig(OpOnBig(x3, ax, ADD), b, ADD)
	if left.Cmp(right) != 0 {
		err := fmt.Sprintf("Point(%v, %v) is not on the curve with a:%v, b:%v\n", x, y, a, b)
		panic(err)
	}

	return &Point{
		x: x,
		y: y,
		a: a,
		b: b,
	}
}

func (p *Point) Add(other *Point) *Point {
	//check two points are on the same curve
	if p.a.Cmp(other.a) != 0 || p.b.Cmp(other.b) != 0 {
		panic("given two points are not on the same curve")
	}

	if p.x == nil {
		return other
	}

	if other.x == nil {
		return p
	}

	//points are on the vertical A(x,y) B(x,-y),
	if p.x.Cmp(other.x) == 0 &&
		OpOnBig(p.y, other.y, ADD).Cmp(big.NewInt(int64(0))) == 0 {
		return &Point{
			x: nil,
			y: nil,
			a: p.a,
			b: p.b,
		}
	}

	//find slope of line AB
	//x1 -> p.x, y1->p.y, x2 ->other.x, y2->other.y
	var numerator *big.Int
	var denominator *big.Int
	if p.x.Cmp(other.x) == 0 && p.y.Cmp(other.y) == 0 {
		//slope = (3*x^2+a)/2y
		xSqrt := OpOnBig(p.x, big.NewInt(int64(2)), EXP)
		threeXSqrt := OpOnBig(xSqrt, big.NewInt(int64(3)), MUL)
		numerator = OpOnBig(threeXSqrt, p.a, ADD)
		//denominator: 2y
		denominator = OpOnBig(p.y, big.NewInt(int64(2)), MUL)
	} else {
		numerator = OpOnBig(other.y, p.y, SUB)   //(y2-y1)
		denominator = OpOnBig(other.x, p.x, SUB) //(x2-x1)
	}

	//s=(y2-y1) / (x2-x1)
	slope := OpOnBig(numerator, denominator, DIV)
	//s^2
	slopeSqrt := OpOnBig(slope, big.NewInt(int64(2)), EXP)
	//x3 = s^2 - x1 - x2
	x3 := OpOnBig(OpOnBig(slopeSqrt, p.x, SUB), other.x, SUB)
	// x3-x1
	x3Minusx1 := OpOnBig(x3, p.x, SUB)
	//y3 = s(x3-x1)+y1
	y3 := OpOnBig(OpOnBig(slope, x3Minusx1, MUL), p.y, ADD)
	//-y3
	minusY3 := OpOnBig(y3, big.NewInt(int64(-1)), MUL)

	return &Point{
		x: x3,
		y: minusY3,
		a: p.a,
		b: p.b,
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("(x:%s, y:%s, a:%s, b:%s)", p.x.String(), p.y.String(),
		p.a.String(), p.b.String())
}

func (p *Point) Equal(other *Point) bool {
	if p.a.Cmp(other.a) == 0 && p.b.Cmp(other.b) == 0 && p.x.Cmp(other.x) == 0 &&
		p.y.Cmp(other.y) == 0 {
		return true
	}

	return false
}

func (p *Point) NotEqual(other *Point) bool {
	if p.a.Cmp(other.a) != 0 || p.b.Cmp(other.b) != 0 || p.x.Cmp(other.x) != 0 ||
		p.y.Cmp(other.y) != 0 {
		return true
	}

	return false
}
