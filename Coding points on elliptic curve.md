Bitcoin rely heavily on a match object called elliptic curve, without this math structure bitcoin will like a castle on beach, it will collapse in any time.What is 
ellipitic curve, its a equation like this: y^2 = x^3 + ax +b, and its shape just like following:

![image](https://github.com/wycl16514/golang-bitcoin-elliptic-curve/assets/7506958/cf4158f0-a3d0-45e2-9423-20d4f41af422)

For bitcoin, its elliptic curve has a name: secp256k1 and its equation is y^2 = x ^ 3 + 7, We don't care too much about the elliptic curve function, we care about
certain set of points on the curve,let's have some code for points on the curve. First we add a new file named point.go under the folder of elliptic-curve, and add
the following code:
```g
package elliptic_curve

import (
	"fmt"
	"math/big"
)

type OP_TPYE int

const (
	ADD OP_TPYE = iota
	SUB
	MUL
	DIV
	EXP
)

type Point struct {
	//for coefficients for elliptic curve
	a *big.Int
	b *big.Int
	//the value of x, y may be ver huge
	x *big.Int
	y *big.Int
}

func OpOnBig(x *big.Int, y *big.Int, opType OP_TPYE) *big.Int {
/*
		why we need to bring operation on big.Int into one function? try following
		var opAdd big.Int
		res := opAdd.Add(big.NewInt(int64(1)), big.NewInt(int64(2)))
		opAdd.Add(big.NewInt(int64(3)), big.NewInt(int64(4)))
		//res is 3 or 7?
		fmt.Printf("val of res is :%d\n", res.String())
	*/
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

	panic("should not come here")
}

func NewEllipticPoint(x *big.Int, y *big.Int, a *big.Int, b *big.Int) *Point {
	//first check (x,y) on the curve defined by a, b
	left := OpOnBig(y, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(x, big.NewInt(int64(3)), EXP)
	ax := OpOnBig(a, x, MUL)
	right := OpOnBig(OpOnBig(x3, ax, ADD), b, ADD)
	if left.Cmp(right) != 0 {
		err := fmt.Sprintf("point:(%v, %v) is not on the curve with a: %v, b:%v\n", x, y, a, b)
		panic(err)
	}

	return &Point{
		a: a,
		b: b,
		x: x,
		y: y,
	}
}

func (p *Point) Equal(other *Point) bool {
	if p.a.Cmp(other.a) == 0 && p.b.Cmp(other.b) == 0 &&
		p.x.Cmp(other.x) == 0 && p.y.Cmp(other.y) == 0 {
		return true
	}

	return false
}

func (p *Point) NoEqual(other *Point) bool {
	if p.a.Cmp(other.a) != 0 || p.b.Cmp(other.b) != 0 ||
		p.x.Cmp(other.x) != 0 || p.y.Cmp(other.y) != 0 {
		return true
	}

	return false
}


```
In the code, we init the elliptic curve point with coefficients of a and b, and check wether the point (x,y) on the curve by computing the value on both side of 
the equation, if they are not equal we throw out an panic. When we check two points are equal or not, we need to compare its four components which are a,b,x,y.

Let's try to new two elliptic points from main function:
```g
func main() {
	/*
		check pint(-1, -1) on y^2 = x^3 + 5x + 7 or not
	*/
	ecc.NewEllipticPoint(big.NewInt(int64(-1)), big.NewInt(int64(-1)),
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	fmt.Println("point(-1, -1) is on curve y^2=x^3+5x+7")
	/*
		check pint(-1, -2) on y^2 = x^3 + 5x + 7 or not
	*/
	ecc.NewEllipticPoint(big.NewInt(int64(-1)), big.NewInt(int64(-2)),
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	fmt.Println("point(-1, -2) is on curve y^2=x^3+5x+7")
}
```
We construct the point struct by using its creator function NewEllipticPoint, the first two params are the coordinate of the point(x,y), if the point is not on the
curve, there would be a panic, otherwise we can print out the message following the creator function, let's run the code for a check by go run main.go and get the
following result:
```g
point(-1, -1) is on curve y^2=x^3+5x+7
panic: point:(-1, -2) is not on the curve with a: 5, b:7
```
we can see point (-1,-1) is on the curve but point (-1, -2) is not on the curve, now we come to practise, please check points (2,4),  (18,77), (5,7) on the curve 
or not.
