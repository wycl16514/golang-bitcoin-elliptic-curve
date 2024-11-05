<img width="640" alt="image" src="https://github.com/wycl16514/golang-bitcoin-elliptic-curve/assets/7506958/6594cc23-53e2-4314-9ca2-c88578374976">

Bitcoin rely heavily on a match object called elliptic curve, without this math structure bitcoin will like a castle on beach, it will collapse in any time.What is 
ellipitic curve, its a equation like this: y^2 = x^3 + ax +b, and its shape just like following:

![image](https://github.com/wycl16514/golang-bitcoin-elliptic-curve/assets/7506958/cf4158f0-a3d0-45e2-9423-20d4f41af422)

For bitcoin, its elliptic curve has a name: secp256k1 and its equation is y^2 = x ^ 3 + 7, We don't care too much about the elliptic curve function, we care about
certain set of points on the curve,let's have some code for points on the curve. First we add a new file named point.go under the folder of elliptic-curve, and add 
the following code:
```go
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
```go
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
```go
point(-1, -1) is on curve y^2=x^3+5x+7
panic: point:(-1, -2) is not on the curve with a: 5, b:7
```
we can see point (-1,-1) is on the curve but point (-1, -2) is not on the curve, now we come to practise, please check points (2,4),  (18,77), (5,7) on the curve or not.

Now we come to the key point, that is given to points A(x1,y1), B(x2,y2) on a given elliptic curve, how we can 
define the addition of them. We use a line to connect the two points, and extend the line, if the extended line can
interset with the curve on a third point C like following:

<img width="651" alt="截屏2024-03-22 15 33 45" src="https://github.com/wycl16514/golang-bitcoin-elliptic-curve/assets/7506958/3926ee42-bdc9-49fb-80af-559552515206">

Then the point that is symetric with c over x-axis is defined as A+B. The same apply to A+C, when we use a line to 
connect A and C, the third point that intersect with the curve is B, then we find the point that is symetric to B over
x-axis would be the result of A+C, the same goes to B+C.

The definition of point addition here have following properties:

1, commutativity, that is A+B = B+A, this is obvious.

2, associativity, that is (A+B) + C = A + (B+C)

The following image shows (A+B)+C:
<img width="645" alt="截屏2024-03-22 15 44 07" src="https://github.com/wycl16514/golang-bitcoin-elliptic-curve/assets/7506958/c6104f2b-7b03-4385-8f8f-580f3bd7b104">

The following image shows A + (B+C):
<img width="737" alt="截屏2024-03-22 15 47 24" src="https://github.com/wycl16514/golang-bitcoin-elliptic-curve/assets/7506958/a77d791b-e089-4526-bb55-fbb7112a4a70">

There is a special case that A and B are on the same vertical line, and no matter how we extend this line, there is 
impossible a third piont that can interset with the curve, but we can give this no existent third point a name called
identitiy marked as I, and define that any point P on the curve, if it add with the identitity point the result is 
itself, that is P + I = P:

<img width="436" alt="截屏2024-03-22 15 52 59" src="https://github.com/wycl16514/golang-bitcoin-elliptic-curve/assets/7506958/ec74254f-9f33-4319-b1b2-151ce7c45109">


How about A and B are the same point on the curve? We defer this case to later time and now let's add some code for 
point addition. First we handle the simple case that is at least one point in the addition is identity point, and 
identity point is with its x and y set to nil, we have code like following:
```go
func NewEllipticPoint(x *big.Int, y *big.Int, a *big.Int, b *big.Int) *Point {
       if x == nil && y == nil {
		return &Point{
			a: a,
			b: b,
			x: x,
			y: y,
		}
	}

	//first check (x,y) on the curve defined by a, b
	left := OpOnBig(y, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(x, big.NewInt(int64(3)), EXP)
	ax := OpOnBig(a, x, MUL)
	right := OpOnBig(OpOnBig(x3, ax, ADD), b, ADD)
	//if x and y are nil, then its identity point and
	//we don't need to check it on curve
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

func (p *Point) Add(other *Point) *Point {
	//check points are on the same curve
	if p.a.Cmp(other.a) != 0 || p.b.Cmp(other.b) != 0 {
		panic("given two points are not on the same curve")
	}

	if p.x == nil {
		//current point is identity point
		return other
	}

	if other.x == nil {
		//the other point is identity
		return p
	}

/*
		another simple case, two points on the same vertical line, that is
		they have the same x but inverse y, the addition of them should be
		identity
	*/
	if p.x.Cmp(other.x) == 0 &&
		OpOnBig(p.y, other.y, ADD).Cmp(big.NewInt(int64(0))) == 0 {
		return &Point{
			x: nil,
			y: nil,
			a: p.a,
			b: p.b,
		}
	}

	//TODO
	return nil
}

func (p *Point) String() string {
	return fmt.Sprintf("x: %s, y: %s, a: %s, b: %s\n", p.x.String(),
		p.y.String(), p.a.String(), p.b.String())
}
```

Let's test the code by adding one point with an identity point and the result should be the point itself:

```go
func main() {
	p := ecc.NewEllipticPoint(big.NewInt(int64(-1)), big.NewInt(int64(-1)),
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	identity := ecc.NewEllipticPoint(nil, nil,
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	fmt.Printf("p is :%s\n", p)

	res := p.Add(identity)
	fmt.Printf("result of point p add to identity is: %s\n", res)

        p2 := ecc.NewEllipticPoint(big.NewInt(int64(-1)), big.NewInt(int64(1)),
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	res = p.Add(p2)
	fmt.Printf("result of adding points on vertical line: %s", res)
}
```

running the code above will have the following result:

```go
p is :(x: -1, y: -1, a: 5, b: 7)

result of point p add to identity is: (x: -1, y: -1, a: 5, b: 7)

result of adding points on vertical line: (x: <nil>, y: <nil>, a: 5, b: 7)
```
If no point is identity, then we need some mathmatical derivation to compute the addition.Given A(x1,y1), B(x2,y2),
we need to know C(x3,y3), first we find the function for line AB.

The slope of line AB is s=(y2-y1)/(x2-x1), and function of line AB is y=s(x-x1)+y1

replace y in y^2=x^3+ax+b by y=s(x-x1)+y1 we have: [s(x-x1)+y1]^2=x^3+ax+b, extend the left side and move it to right 
side we get:
<img width="640" alt="image" src="https://github.com/wycl16514/golang-bitcoin-elliptic-curve/assets/7506958/971f7bb2-6ce5-4cd0-8a3a-4cd6a6a489d8">

we know certainly that x1, x2, x3 are solution for eqution above and below:

<img width="642" alt="截屏2024-03-22 22 31 56" src="https://github.com/wycl16514/golang-bitcoin-elliptic-curve/assets/7506958/8f369b26-4b5b-4ca3-862e-2ee772940845">

By Vieta's foluma, coefficient for term with the same order should equal, then we have:

s^2 = x1+x2+x3 => x3 = s^2 - x1 - x2

This give us the value of x3, we can put it to the line function and get y3:

y=s(x-x1)+y1=> y3=s(x3-x1)+y1

Now we have C, and remember we need to reflect it over x-axis, and we get A+B is (-s^2-x1-x2, -[s(x3-x1)+y1])

Let's have code for this:
```go
func (p *Point) Add(other *Point) *Point {
 ...
 //find slope of line AB
	//x1 = p.x, y1 = p.y, x2 = other.x, y2 = other.y
	numerator := OpOnBig(other.y, p.y, SUB)
	denominator := OpOnBig(other.x, p.x, SUB)
	//s = (y2-y2)/(x2-x1)
	slope := OpOnBig(numerator, denominator, DIV)

	//-s^2
	slopeSqrt := OpOnBig(slope, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(OpOnBig(slopeSqrt, p.x, SUB), other.x, SUB)
	//x3-x1
	x3Minusx1 := OpOnBig(x3, p.x, SUB)
	//y3=s(x3-x1)+y1
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
```
Let's test the code above:
```go
func main() {
	//C = A(2,5) + B(-1, -1)
	A := ecc.NewEllipticPoint(big.NewInt(int64(2)), big.NewInt(int64(5)),
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	B := ecc.NewEllipticPoint(big.NewInt(int64(-1)), big.NewInt(int64(-1)),
		big.NewInt(int64(5)), big.NewInt(int64(7)))
	C := A.Add(B)
	fmt.Printf("A(2,5) + B(-1,-1) = %s\n", C)
}
```
The running result for the above code is:
```g
A(2,5) + B(-1,-1) = (x: 3, y: -7, a: 5, b: 7)
```
Please use a pen and paper to check the result if you are hesitate, I can assure you the result is correct. There still one special case we need to handle,that 
is when A=B, in this case, the line AB turns into a tangent line of the elliptic curve :

<img width="653" alt="截屏2024-03-23 03 44 04" src="https://github.com/wycl16514/golang-bitcoin-elliptic-curve/assets/7506958/8666bc6b-174e-4520-973f-614c248bedce">

This time we can't get the slope for the line easily, we need the help from calculus. To find the slope of tangent line on a curve, we need to compute the 
Derivative at the point of the curve, for the function of the curve y^2 = x^3 + ax+b, we take drivative for point x on both side and we get:
d(y^2)/dx = d(x^3+ax+b)/dx -> 2y * dy/dx = 3x^2 + a -> dy/dx = (3x^2+a)/2y, therefore we only need to change the computation of slope and other steps remain 
the same, here is the code:
```go
func (p *Point) Add(other *Point) *Point {
...
//find slope of line AB
	//x1 = p.x, y1 = p.y, x2 = other.x, y2 = other.y
	var numerator *big.Int
	var denominator *big.Int
	if p.x.Cmp(other.x) == 0 && p.y.Cmp(other.y) == 0 {
		//two points are the same and compute the slope of tangent line
		//numerator is (3x^2+a)
		xSqrt := OpOnBig(p.x, big.NewInt(int64(2)), EXP)
		threeXSqrt := OpOnBig(xSqrt, big.NewInt(int64(3)), MUL)
		numerator = OpOnBig(threeXSqrt, p.a, ADD)
		//demoninator is 2y
		denominator = OpOnBig(p.y, big.NewInt(int64(2)), MUL)
	} else {
		//s = (y2-y2)/(x2-x1)
		numerator = OpOnBig(other.y, p.y, SUB)
		denominator = OpOnBig(other.x, p.x, SUB)
	}
...
}
```
Let's have a test for this case:
```go
func main() {
...
        //B=(-1,-1) C=B+B
	C = B.Add(B)
	fmt.Printf("B(-1,-1) + B(-1,-1) = %s\n", C)
}
```
The result for the above code is :
```go
B(-1,-1) + B(-1,-1) = (x: 18, y: 77, a: 5, b: 7)
```
