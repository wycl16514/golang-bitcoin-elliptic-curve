package elliptic_curve

import (
	"fmt"
	"math/big"
)

type FieldElement struct {
	order *big.Int //field order
	num   *big.Int //value of the given element in the field
}

// overlow 64bits integer
//huge number ,  +, *, ^ => overflow 64bits, we use large or big number
//

func NewFieldElement(order *big.Int, num *big.Int) *FieldElement {
	/*
		init function for FieldElement
	*/
	if order.Cmp(num) == -1 {
		err := fmt.Sprintf("Num not in the range of 0 to %d", order)
		panic(err)
	}

	return &FieldElement{
		order: order,
		num:   num,
	}
}

func (f *FieldElement) String() string {
	//__repr for python
	return fmt.Sprintf("FieldElement{order: %s, num: %s}\n", f.order.String(),
		f.num.String())
}

func (f *FieldElement) EqualTo(other *FieldElement) bool {
	return f.order.Cmp(other.order) == 0 && f.num.Cmp(other.num) == 0
}

func (f *FieldElement) checkOrder(other *FieldElement) {
	if f.order.Cmp(other.order) != 0 {
		panic("add need to do on the field element with the same order")
	}
}

func (f *FieldElement) Add(other *FieldElement) *FieldElement {
	f.checkOrder(other)

	//remember the modulur
	//operator overloading for +, __add__ python
	var op big.Int
	return NewFieldElement(f.order, op.Mod(op.Add(f.num, other.num), f.order))
}

/*
a ,b (a + b) % order = 0, b is called negate of a, b = -a
*/
func (f *FieldElement) Negate() *FieldElement {
	/*
		b, (a+b) % order = 0, b = order - a, (a+b) % order => (a + order - a) => order => order%order = 0
	*/
	var op big.Int
	return NewFieldElement(f.order, op.Sub(f.order, f.num))
}

func (f *FieldElement) Subtract(other *FieldElement) *FieldElement {
	/*
		a, b element of the finite set,  c = a - b, given b how can we find c,
		(b + c) % order = a, a - b => (a + (-b) ) % order
	*/

	return f.Add(other.Negate())
}

func (f *FieldElement) Multiply(other *FieldElement) *FieldElement {
	f.checkOrder(other)

	//Arithmetic multiplie over modulur of the order
	var op big.Int
	mul := op.Mul(f.num, other.num)
	return NewFieldElement(f.order, op.Mod(mul, f.order))
}

func (f *FieldElement) Power(power *big.Int) *FieldElement {
	//Arithmetic power over modulur of the order
	// k ^(p -1) % p = 1, power > p-1 => power %(p-1)
	var op big.Int
	t := op.Mod(power, op.Sub(f.order, big.NewInt(int64(1))))
	powerRes := op.Exp(f.num, t, nil)
	modRes := op.Mod(powerRes, f.order)
	return NewFieldElement(f.order, modRes)
}

func (f *FieldElement) ScalarMul(val *big.Int) *FieldElement {
	var op big.Int
	res := op.Mul(f.num, val)
	res = op.Mod(res, f.order)
	return NewFieldElement(f.order, res)
}

func (f *FieldElement) Divide(other *FieldElement) *FieldElement {
	f.checkOrder(other)
	//a / b => a . b^(p-2)
	var op big.Int
	otherReverse := other.Power(op.Sub(f.order, big.NewInt(int64(2))))
	return f.Multiply(otherReverse)
}
