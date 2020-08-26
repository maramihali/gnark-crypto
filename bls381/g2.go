// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gurvy DO NOT EDIT

package bls381

import (
	"math"
	"math/big"
	"runtime"

	"github.com/consensys/gurvy/bls381/fr"
	"github.com/consensys/gurvy/utils"
	"github.com/consensys/gurvy/utils/parallel"
)

// G2Jac is a point with E2 coordinates
type G2Jac struct {
	X, Y, Z E2
}

// G2Proj point in projective coordinates
type G2Proj struct {
	X, Y, Z E2
}

// G2Affine point in affine coordinates
type G2Affine struct {
	X, Y E2
}

//  g2JacExtended parameterized jacobian coordinates (x=X/ZZ, y=Y/ZZZ, ZZ**3=ZZZ**2)
type g2JacExtended struct {
	X, Y, ZZ, ZZZ E2
}

// SetInfinity sets p to O
func (p *g2JacExtended) SetInfinity() *g2JacExtended {
	p.X.SetOne()
	p.Y.SetOne()
	p.ZZ.SetZero()
	p.ZZZ.SetZero()
	return p
}

// ToAffine sets p in affine coords
func (p *g2JacExtended) ToAffine(Q *G2Affine) *G2Affine {
	var zero E2
	if p.ZZ.Equal(&zero) {
		Q.X.Set(&zero)
		Q.Y.Set(&zero)
		return Q
	}
	Q.X.Inverse(&p.ZZ).Mul(&Q.X, &p.X)
	Q.Y.Inverse(&p.ZZZ).Mul(&Q.Y, &p.Y)
	return Q
}

// ToJac sets p in affine coords
func (p *g2JacExtended) ToJac(Q *G2Jac) *G2Jac {
	var zero E2
	if p.ZZ.Equal(&zero) {
		Q.Set(&g2Infinity)
		return Q
	}
	Q.X.Mul(&p.ZZ, &p.X).Mul(&Q.X, &p.ZZ)
	Q.Y.Mul(&p.ZZZ, &p.Y).Mul(&Q.Y, &p.ZZZ)
	Q.Z.Set(&p.ZZZ)
	return Q
}

// unsafeToJac sets p in affine coords, but don't check for infinity
func (p *g2JacExtended) unsafeToJac(Q *G2Jac) *G2Jac {
	Q.X.Mul(&p.ZZ, &p.X).Mul(&Q.X, &p.ZZ)
	Q.Y.Mul(&p.ZZZ, &p.Y).Mul(&Q.Y, &p.ZZZ)
	Q.Z.Set(&p.ZZZ)
	return Q
}

// mSub
// http://www.hyperelliptic.org/EFD/ g2p/auto-shortw-xyzz.html#addition-madd-2008-s
func (p *g2JacExtended) mSub(a *G2Affine) *g2JacExtended {

	//if a is infinity return p
	if a.X.IsZero() && a.Y.IsZero() {
		return p
	}
	// p is infinity, return a
	if p.ZZ.IsZero() {
		p.X = a.X
		p.Y = a.Y
		p.Y.Neg(&p.Y)
		p.ZZ.SetOne()
		p.ZZZ.SetOne()
		return p
	}

	var U2, S2, P, R, PP, PPP, Q, Q2, RR, X3, Y3 E2

	// p2: a, p1: p
	U2.Mul(&a.X, &p.ZZ)
	S2.Mul(&a.Y, &p.ZZZ)
	S2.Neg(&S2)

	P.Sub(&U2, &p.X)
	R.Sub(&S2, &p.Y)

	pIsZero := P.IsZero()
	rIsZero := R.IsZero()

	if pIsZero && rIsZero {
		return p.doubleNeg(a)
	} else if pIsZero {
		p.ZZ.SetZero()
		p.ZZZ.SetZero()
		return p
	}

	PP.Square(&P)
	PPP.Mul(&P, &PP)
	Q.Mul(&p.X, &PP)
	RR.Square(&R)
	X3.Sub(&RR, &PPP)
	Q2.Double(&Q)
	p.X.Sub(&X3, &Q2)
	Y3.Sub(&Q, &p.X).Mul(&Y3, &R)
	R.Mul(&p.Y, &PPP)
	p.Y.Sub(&Y3, &R)
	p.ZZ.Mul(&p.ZZ, &PP)
	p.ZZZ.Mul(&p.ZZZ, &PPP)

	return p
}

// mAdd
// http://www.hyperelliptic.org/EFD/ g2p/auto-shortw-xyzz.html#addition-madd-2008-s
func (p *g2JacExtended) mAdd(a *G2Affine) *g2JacExtended {

	//if a is infinity return p
	if a.X.IsZero() && a.Y.IsZero() {
		return p
	}
	// p is infinity, return a
	if p.ZZ.IsZero() {
		p.X = a.X
		p.Y = a.Y
		p.ZZ.SetOne()
		p.ZZZ.SetOne()
		return p
	}

	var U2, S2, P, R, PP, PPP, Q, Q2, RR, X3, Y3 E2

	// p2: a, p1: p
	U2.Mul(&a.X, &p.ZZ)
	S2.Mul(&a.Y, &p.ZZZ)

	P.Sub(&U2, &p.X)
	R.Sub(&S2, &p.Y)

	pIsZero := P.IsZero()
	rIsZero := R.IsZero()

	if pIsZero && rIsZero {
		return p.double(a)
	} else if pIsZero {
		p.ZZ.SetZero()
		p.ZZZ.SetZero()
		return p
	}

	PP.Square(&P)
	PPP.Mul(&P, &PP)
	Q.Mul(&p.X, &PP)
	RR.Square(&R)
	X3.Sub(&RR, &PPP)
	Q2.Double(&Q)
	p.X.Sub(&X3, &Q2)
	Y3.Sub(&Q, &p.X).Mul(&Y3, &R)
	R.Mul(&p.Y, &PPP)
	p.Y.Sub(&Y3, &R)
	p.ZZ.Mul(&p.ZZ, &PP)
	p.ZZZ.Mul(&p.ZZZ, &PPP)

	return p
}

func (p *g2JacExtended) doubleNeg(q *G2Affine) *g2JacExtended {

	var U, S, M, _M, Y3 E2

	U.Double(&q.Y)
	U.Neg(&U)
	p.ZZ.Square(&U)
	p.ZZZ.Mul(&U, &p.ZZ)
	S.Mul(&q.X, &p.ZZ)
	_M.Square(&q.X)
	M.Double(&_M).
		Add(&M, &_M) // -> + a, but a=0 here
	p.X.Square(&M).
		Sub(&p.X, &S).
		Sub(&p.X, &S)
	Y3.Sub(&S, &p.X).Mul(&Y3, &M)
	U.Mul(&p.ZZZ, &q.Y)
	U.Neg(&U)
	p.Y.Sub(&Y3, &U)

	return p
}

// double point in ZZ coords
// http://www.hyperelliptic.org/EFD/ g2p/auto-shortw-xyzz.html#doubling-dbl-2008-s-1
func (p *g2JacExtended) double(q *G2Affine) *g2JacExtended {

	var U, S, M, _M, Y3 E2

	U.Double(&q.Y)
	p.ZZ.Square(&U)
	p.ZZZ.Mul(&U, &p.ZZ)
	S.Mul(&q.X, &p.ZZ)
	_M.Square(&q.X)
	M.Double(&_M).
		Add(&M, &_M) // -> + a, but a=0 here
	p.X.Square(&M).
		Sub(&p.X, &S).
		Sub(&p.X, &S)
	Y3.Sub(&S, &p.X).Mul(&Y3, &M)
	U.Mul(&p.ZZZ, &q.Y)
	p.Y.Sub(&Y3, &U)

	return p
}

// Set set p to the provided point
func (p *G2Jac) Set(a *G2Jac) *G2Jac {
	p.X.Set(&a.X)
	p.Y.Set(&a.Y)
	p.Z.Set(&a.Z)
	return p
}

// Equal tests if two points (in Jacobian coordinates) are equal
func (p *G2Jac) Equal(a *G2Jac) bool {

	if p.Z.IsZero() && a.Z.IsZero() {
		return true
	}
	_p := G2Affine{}
	_p.FromJacobian(p)

	_a := G2Affine{}
	_a.FromJacobian(a)

	return _p.X.Equal(&_a.X) && _p.Y.Equal(&_a.Y)
}

// Equal tests if two points (in Affine coordinates) are equal
func (p *G2Affine) Equal(a *G2Affine) bool {
	return p.X.Equal(&a.X) && p.Y.Equal(&a.Y)
}

// Neg computes -G
func (p *G2Jac) Neg(a *G2Jac) *G2Jac {
	p.Set(a)
	p.Y.Neg(&a.Y)
	return p
}

// Neg computes -G
func (p *G2Affine) Neg(a *G2Affine) *G2Affine {
	p.X.Set(&a.X)
	p.Y.Neg(&a.Y)
	return p
}

// SubAssign substracts two points on the curve
func (p *G2Jac) SubAssign(a G2Jac) *G2Jac {
	a.Y.Neg(&a.Y)
	p.AddAssign(&a)
	return p
}

// FromJacobian rescale a point in Jacobian coord in z=1 plane
func (p *G2Affine) FromJacobian(p1 *G2Jac) *G2Affine {

	var a, b E2

	if p1.Z.IsZero() {
		p.X.SetZero()
		p.Y.SetZero()
		return p
	}

	a.Inverse(&p1.Z)
	b.Square(&a)
	p.X.Mul(&p1.X, &b)
	p.Y.Mul(&p1.Y, &b).Mul(&p.Y, &a)

	return p
}

// FromJacobian converts a point from Jacobian to projective coordinates
func (p *G2Proj) FromJacobian(Q *G2Jac) *G2Proj {
	// memalloc
	var buf E2
	buf.Square(&Q.Z)

	p.X.Mul(&Q.X, &Q.Z)
	p.Y.Set(&Q.Y)
	p.Z.Mul(&Q.Z, &buf)

	return p
}

func (p *G2Jac) String() string {
	if p.Z.IsZero() {
		return "O"
	}
	_p := G2Affine{}
	_p.FromJacobian(p)
	return "E([" + _p.X.String() + "," + _p.Y.String() + "]),"
}

// FromAffine sets p = Q, p in Jacboian, Q in affine
func (p *G2Jac) FromAffine(Q *G2Affine) *G2Jac {
	if Q.X.IsZero() && Q.Y.IsZero() {
		p.Z.SetZero()
		p.X.SetOne()
		p.Y.SetOne()
		return p
	}
	p.Z.SetOne()
	p.X.Set(&Q.X)
	p.Y.Set(&Q.Y)
	return p
}

func (p *G2Affine) String() string {
	var x, y E2
	x.Set(&p.X)
	y.Set(&p.Y)
	return "E([" + x.String() + "," + y.String() + "]),"
}

// IsInfinity checks if the point is infinity (in affine, it's encoded as (0,0))
func (p *G2Affine) IsInfinity() bool {
	return p.X.IsZero() && p.Y.IsZero()
}

// IsOnCurve returns true if p in on the curve
func (p *G2Proj) IsOnCurve() bool {
	var left, right, tmp E2
	left.Square(&p.Y).
		Mul(&left, &p.Z)
	right.Square(&p.X).
		Mul(&right, &p.X)
	tmp.Square(&p.Z).
		Mul(&tmp, &p.Z).
		Mul(&tmp, &Btwist)
	right.Add(&right, &tmp)
	return left.Equal(&right)
}

// IsOnCurve returns true if p in on the curve
func (p *G2Jac) IsOnCurve() bool {
	var left, right, tmp E2
	left.Square(&p.Y)
	right.Square(&p.X).Mul(&right, &p.X)
	tmp.Square(&p.Z).
		Square(&tmp).
		Mul(&tmp, &p.Z).
		Mul(&tmp, &p.Z).
		Mul(&tmp, &Btwist)
	right.Add(&right, &tmp)
	return left.Equal(&right)
}

// IsOnCurve returns true if p in on the curve
func (p *G2Affine) IsOnCurve() bool {
	var point G2Jac
	point.FromAffine(p)
	return point.IsOnCurve() // call this function to handle infinity point
}

// AddAssign point addition in montgomery form
// https://hyperelliptic.org/EFD/g2p/auto-shortw-jacobian-3.html#addition-add-2007-bl
func (p *G2Jac) AddAssign(a *G2Jac) *G2Jac {

	// p is infinity, return a
	if p.Z.IsZero() {
		p.Set(a)
		return p
	}

	// a is infinity, return p
	if a.Z.IsZero() {
		return p
	}

	var Z1Z1, Z2Z2, U1, U2, S1, S2, H, I, J, r, V E2
	Z1Z1.Square(&a.Z)
	Z2Z2.Square(&p.Z)
	U1.Mul(&a.X, &Z2Z2)
	U2.Mul(&p.X, &Z1Z1)
	S1.Mul(&a.Y, &p.Z).
		Mul(&S1, &Z2Z2)
	S2.Mul(&p.Y, &a.Z).
		Mul(&S2, &Z1Z1)

	// if p == a, we double instead
	if U1.Equal(&U2) && S1.Equal(&S2) {
		return p.DoubleAssign()
	}

	H.Sub(&U2, &U1)
	I.Double(&H).
		Square(&I)
	J.Mul(&H, &I)
	r.Sub(&S2, &S1).Double(&r)
	V.Mul(&U1, &I)
	p.X.Square(&r).
		Sub(&p.X, &J).
		Sub(&p.X, &V).
		Sub(&p.X, &V)
	p.Y.Sub(&V, &p.X).
		Mul(&p.Y, &r)
	S1.Mul(&S1, &J).Double(&S1)
	p.Y.Sub(&p.Y, &S1)
	p.Z.Add(&p.Z, &a.Z)
	p.Z.Square(&p.Z).
		Sub(&p.Z, &Z1Z1).
		Sub(&p.Z, &Z2Z2).
		Mul(&p.Z, &H)

	return p
}

// AddMixed point addition
// http://www.hyperelliptic.org/EFD/g2p/auto-shortw-jacobian-0.html#addition-madd-2007-bl
func (p *G2Jac) AddMixed(a *G2Affine) *G2Jac {

	//if a is infinity return p
	if a.X.IsZero() && a.Y.IsZero() {
		return p
	}
	// p is infinity, return a
	if p.Z.IsZero() {
		p.X = a.X
		p.Y = a.Y
		p.Z.SetOne()
		return p
	}

	// get some Element from our pool
	var Z1Z1, U2, S2, H, HH, I, J, r, V E2
	Z1Z1.Square(&p.Z)
	U2.Mul(&a.X, &Z1Z1)
	S2.Mul(&a.Y, &p.Z).
		Mul(&S2, &Z1Z1)

	// if p == a, we double instead
	if U2.Equal(&p.X) && S2.Equal(&p.Y) {
		return p.DoubleAssign()
	}

	H.Sub(&U2, &p.X)
	HH.Square(&H)
	I.Double(&HH).Double(&I)
	J.Mul(&H, &I)
	r.Sub(&S2, &p.Y).Double(&r)
	V.Mul(&p.X, &I)
	p.X.Square(&r).
		Sub(&p.X, &J).
		Sub(&p.X, &V).
		Sub(&p.X, &V)
	J.Mul(&J, &p.Y).Double(&J)
	p.Y.Sub(&V, &p.X).
		Mul(&p.Y, &r)
	p.Y.Sub(&p.Y, &J)
	p.Z.Add(&p.Z, &H)
	p.Z.Square(&p.Z).
		Sub(&p.Z, &Z1Z1).
		Sub(&p.Z, &HH)

	return p
}

// Double doubles a point in Jacobian coordinates
// https://hyperelliptic.org/EFD/g2p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G2Jac) Double(q *G2Jac) *G2Jac {
	p.Set(q)
	p.DoubleAssign()
	return p
}

// DoubleAssign doubles a point in Jacobian coordinates
// https://hyperelliptic.org/EFD/g2p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G2Jac) DoubleAssign() *G2Jac {

	// get some Element from our pool
	var XX, YY, YYYY, ZZ, S, M, T E2

	XX.Square(&p.X)
	YY.Square(&p.Y)
	YYYY.Square(&YY)
	ZZ.Square(&p.Z)
	S.Add(&p.X, &YY)
	S.Square(&S).
		Sub(&S, &XX).
		Sub(&S, &YYYY).
		Double(&S)
	M.Double(&XX).Add(&M, &XX)
	p.Z.Add(&p.Z, &p.Y).
		Square(&p.Z).
		Sub(&p.Z, &YY).
		Sub(&p.Z, &ZZ)
	T.Square(&M)
	p.X = T
	T.Double(&S)
	p.X.Sub(&p.X, &T)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M)
	YYYY.Double(&YYYY).Double(&YYYY).Double(&YYYY)
	p.Y.Sub(&p.Y, &YYYY)

	return p
}

// ScalarMulByGen multiplies given scalar by generator
func (p *G2Jac) ScalarMulByGen(s *big.Int) *G2Jac {
	return p.ScalarMulGLV(&g2GenAff, s)
}

// ScalarMultiplication 2-bits windowed exponentiation
func (p *G2Jac) ScalarMultiplication(a *G2Affine, s *big.Int) *G2Jac {

	var res, tmp G2Jac
	var ops [3]G2Affine

	res.Set(&g2Infinity)
	ops[0] = *a
	tmp.FromAffine(a).DoubleAssign()
	ops[1].FromJacobian(&tmp)
	tmp.AddMixed(a)
	ops[2].FromJacobian(&tmp)

	b := s.Bytes()
	for i := range b {
		w := b[i]
		mask := byte(0xc0)
		for j := 0; j < 4; j++ {
			res.DoubleAssign().DoubleAssign()
			c := (w & mask) >> (6 - 2*j)
			if c != 0 {
				res.AddMixed(&ops[c-1])
			}
			mask = mask >> 2
		}
	}
	p.Set(&res)

	return p

}

// phi assigns p to phi(a) where phi: (x,y)->(ux,y), and returns p
func (p *G2Jac) phi(a *G2Affine) *G2Jac {
	p.FromAffine(a)

	p.X.MulByElement(&p.X, &thirdRootOneG2)

	return p
}

// ScalarMulGLV performs scalar multiplication using GLV
func (p *G2Jac) ScalarMulGLV(a *G2Affine, s *big.Int) *G2Jac {

	var table [3]G2Jac
	var zero big.Int
	var res G2Jac
	var k1, k2 fr.Element

	res.Set(&g2Infinity)

	// table stores [+-a, +-phi(a), +-a+-phi(a)]
	table[0].FromAffine(a)
	table[1].phi(a)

	// split the scalar, modifies +-a, phi(a) accordingly
	k := utils.SplitScalar(s, &glvBasis)

	if k[0].Cmp(&zero) == -1 {
		k[0].Neg(&k[0])
		table[0].Neg(&table[0])
	}
	if k[1].Cmp(&zero) == -1 {
		k[1].Neg(&k[1])
		table[1].Neg(&table[1])
	}
	table[2].Set(&table[0]).AddAssign(&table[1])

	// bounds on the lattice base vectors guarantee that k1, k2 are len(r)/2 bits long max
	k1.SetBigInt(&k[0]).FromMont()
	k2.SetBigInt(&k[1]).FromMont()

	// loop starts from len(k1)/2 due to the bounds
	for i := len(k1)/2 - 1; i >= 0; i-- {
		mask := uint64(1) << 63
		for j := 0; j < 64; j++ {
			res.Double(&res)
			b1 := (k1[i] & mask) >> (63 - j)
			b2 := (k2[i] & mask) >> (63 - j)
			if b1|b2 != 0 {
				s := (b2<<1 | b1)
				res.AddAssign(&table[s-1])
			}
			mask = mask >> 1
		}
	}

	p.Set(&res)
	return p
}

// MultiExp implements section 4 of https://eprint.iacr.org/2012/549.pdf
func (p *G2Jac) MultiExp(points []G2Affine, scalars []fr.Element) *G2Jac {
	// note:
	// each of the multiExpcX method is the same, except for the c constant it declares
	// duplicating (through template generation) these methods allows to declare the buckets on the stack
	// the choice of c needs to be improved:
	// there is a theoritical value that gives optimal asymptotics
	// but in practice, other factors come into play, including:
	// * if c doesn't divide 64, the word size, then we're bound to select bits over 2 words of our scalars, instead of 1
	// * number of CPUs
	// * cache friendliness (which depends on the host, G1 or G2... )
	//	--> for example, on BN256, a G1 point fits into one cache line of 64bytes, but a G2 point don't.

	// for each multiExpcX
	// step 1
	// we compute, for each scalars over c-bit wide windows, nbChunk digits
	// if the digit is larger than 2^{c-1}, then, we borrow 2^c from the next window and substract
	// 2^{c} to the current digit, making it negative.
	// negative digits will be processed in the next step as adding -G into the bucket instead of G
	// (computing -G is cheap, and this saves us half of the buckets)
	// step 2
	// buckets are declared on the stack
	// notice that we have 2^{c-1} buckets instead of 2^{c} (see step1)
	// we use jacobian extended formulas here as they are faster than mixed addition
	// bucketAccumulate places points into buckets base on their selector and return the weighted bucket sum in given channel
	// step 3
	// reduce the buckets weigthed sums into our result (chunkReduce)

	// approximate cost (in group operations)
	// cost = bits/c * (nbPoints + 2^{c-1})
	// this needs to be verified empirically.
	// for example, on a MBP 2016, for G2 MultiExp > 8M points, hand picking c gives better results
	implementedCs := []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	nbPoints := len(points)
	min := math.MaxFloat64
	bestC := 0
	for _, c := range implementedCs {
		cc := fr.Limbs * 64 * (nbPoints + (1 << (c - 1)))
		cost := float64(cc) / float64(c)
		if cost < min {
			min = cost
			bestC = c
		}
	}

	// semaphore to limit number of cpus
	numCpus := runtime.NumCPU()
	chCpus := make(chan struct{}, numCpus)
	for i := 0; i < numCpus; i++ {
		chCpus <- struct{}{}
	}

	switch bestC {

	case 4:
		return p.multiExpc4(points, scalars, chCpus)

	case 5:
		return p.multiExpc5(points, scalars, chCpus)

	case 6:
		return p.multiExpc6(points, scalars, chCpus)

	case 7:
		return p.multiExpc7(points, scalars, chCpus)

	case 8:
		return p.multiExpc8(points, scalars, chCpus)

	case 9:
		return p.multiExpc9(points, scalars, chCpus)

	case 10:
		return p.multiExpc10(points, scalars, chCpus)

	case 11:
		return p.multiExpc11(points, scalars, chCpus)

	case 12:
		return p.multiExpc12(points, scalars, chCpus)

	case 13:
		return p.multiExpc13(points, scalars, chCpus)

	case 14:
		return p.multiExpc14(points, scalars, chCpus)

	case 15:
		return p.multiExpc15(points, scalars, chCpus)

	case 16:
		return p.multiExpc16(points, scalars, chCpus)

	case 17:
		return p.multiExpc17(points, scalars, chCpus)

	case 18:
		return p.multiExpc18(points, scalars, chCpus)

	case 19:
		return p.multiExpc19(points, scalars, chCpus)

	case 20:
		return p.multiExpc20(points, scalars, chCpus)

	default:
		panic("unimplemented")
	}
}

// chunkReduceG2 reduces the weighted sum of the buckets into the result of the multiExp
func chunkReduceG2(p *G2Jac, c int, chTotals []chan G2Jac) *G2Jac {
	totalj := <-chTotals[len(chTotals)-1]
	p.Set(&totalj)
	for j := len(chTotals) - 2; j >= 0; j-- {
		for l := 0; l < c; l++ {
			p.DoubleAssign()
		}
		totalj := <-chTotals[j]
		p.AddAssign(&totalj)
	}
	return p
}

func bucketAccumulateG2(chunk uint64,
	chRes chan<- G2Jac,
	chCpus chan struct{},
	buckets []g2JacExtended,
	c uint64,
	points []G2Affine,
	scalars []fr.Element) {

	<-chCpus // wait and decrement avaiable CPUs on the semaphore

	mask := uint64((1 << c) - 1) // low c bits are 1
	msbWindow := uint64(1 << (c - 1))

	for i := 0; i < len(buckets); i++ {
		buckets[i].SetInfinity()
	}

	jc := uint64(chunk * c)
	s := selector{}
	s.index = jc / 64
	s.shift = jc - (s.index * 64)
	s.mask = mask << s.shift
	s.multiWordSelect = (64%c) != 0 && s.shift > (64-c) && s.index < (fr.Limbs-1)
	if s.multiWordSelect {
		nbBitsHigh := s.shift - uint64(64-c)
		s.maskHigh = (1 << nbBitsHigh) - 1
		s.shiftHigh = (c - nbBitsHigh)
	}

	// for each scalars, get the digit corresponding to the chunk we're processing.
	for i := 0; i < len(scalars); i++ {
		bits := (scalars[i][s.index] & s.mask) >> s.shift
		if s.multiWordSelect {
			bits += (scalars[i][s.index+1] & s.maskHigh) << s.shiftHigh
		}

		if bits == 0 {
			continue
		}

		// if msbWindow bit is set, we need to substract
		if bits&msbWindow == 0 {
			// add
			buckets[bits-1].mAdd(&points[i])
		} else {
			// sub
			buckets[bits & ^msbWindow].mSub(&points[i])
		}
	}

	// reduce buckets into total
	// total =  bucket[0] + 2*bucket[1] + 3*bucket[2] ... + n*bucket[n-1]

	var runningSum, tj, total G2Jac
	runningSum.Set(&g2Infinity)
	total.Set(&g2Infinity)
	for k := len(buckets) - 1; k >= 0; k-- {
		if !buckets[k].ZZ.IsZero() {
			runningSum.AddAssign(buckets[k].unsafeToJac(&tj))
		}
		total.AddAssign(&runningSum)
	}

	chRes <- total
	close(chRes)
	chCpus <- struct{}{} // increment avaiable CPUs into the semaphore
}

func (p *G2Jac) multiExpc4(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 4                          // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc5(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 5                              // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc6(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 6                              // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc7(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 7                              // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc8(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 8                          // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc9(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 9                              // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc10(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 10                             // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc11(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 11                             // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc12(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 12                             // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc13(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 13                             // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc14(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 14                             // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc15(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 15                             // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc16(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 16                         // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc17(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 17                             // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc18(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 18                             // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc19(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 19                             // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

func (p *G2Jac) multiExpc20(points []G2Affine, scalars []fr.Element, chCpus chan struct{}) *G2Jac {

	const c = 20                             // scalars partitioned into c-bit radixes
	const nbChunks = (fr.Limbs * 64 / c) + 1 // number of c-bit radixes in a scalar

	// 1 channel per chunk, which will contain the weighted sum of the its buckets
	var chTotals [nbChunks]chan G2Jac
	for i := 0; i < nbChunks; i++ {
		chTotals[i] = make(chan G2Jac, 1)
	}

	newScalars := ScalarsToDigits(scalars, c)

	for chunk := nbChunks - 1; chunk >= 0; chunk-- {
		go func(j uint64) {
			var buckets [1 << (c - 1)]g2JacExtended
			bucketAccumulateG2(j, chTotals[j], chCpus, buckets[:], c, points, newScalars)
		}(uint64(chunk))
	}

	return chunkReduceG2(p, c, chTotals[:])
}

// BatchScalarMultiplicationG2 multiplies the same base (generator) by all scalars
// and return resulting points in affine coordinates
// currently uses a simple windowed-NAF like exponentiation algorithm, and use fixed windowed size (16 bits)
// TODO : implement variable window size depending on input size
func BatchScalarMultiplicationG2(base *G2Affine, scalars []fr.Element) []G2Affine {
	const c = 16 // window size
	const nbChunks = fr.Limbs * 64 / c
	const mask uint64 = (1 << c) - 1 // low c bits are 1
	const msbWindow uint64 = (1 << (c - 1))

	// precompute all powers of base for our window
	var baseTable [(1 << (c - 1))]G2Jac
	baseTable[0].Set(&g2Infinity)
	baseTable[0].AddMixed(base)
	for i := 1; i < len(baseTable); i++ {
		baseTable[i] = baseTable[i-1]
		baseTable[i].AddMixed(base)
	}

	newScalars := ScalarsToDigits(scalars, c)

	// compute offset and word selector / shift to select the right bits of our windows
	selectors := make([]selector, nbChunks)
	for chunk := uint64(0); chunk < nbChunks; chunk++ {
		jc := uint64(chunk * c)
		d := selector{}
		d.index = jc / 64
		d.shift = jc - (d.index * 64)
		d.mask = mask << d.shift
		d.multiWordSelect = (64%c) != 0 && d.shift > (64-c) && d.index < (fr.Limbs-1)
		if d.multiWordSelect {
			nbBitsHigh := d.shift - uint64(64-c)
			d.maskHigh = (1 << nbBitsHigh) - 1
			d.shiftHigh = (c - nbBitsHigh)
		}
		selectors[chunk] = d
	}

	toReturn := make([]G2Affine, len(scalars))

	// for each digit, take value in the base table, double it c time, voila.
	parallel.Execute(len(newScalars), func(start, end int) {
		var p G2Jac
		for i := start; i < end; i++ {
			p.Set(&g2Infinity)

			for chunk := nbChunks - 1; chunk >= 0; chunk-- {
				s := selectors[chunk]
				if chunk != nbChunks-1 {
					for j := 0; j < c; j++ {
						p.DoubleAssign()
					}
				}

				bits := (newScalars[i][s.index] & s.mask) >> s.shift
				if s.multiWordSelect {
					bits += (newScalars[i][s.index+1] & s.maskHigh) << s.shiftHigh
				}

				if bits == 0 {
					continue
				}

				// if msbWindow bit is set, we need to substract
				if bits&msbWindow == 0 {
					// add
					p.AddAssign(&baseTable[bits-1])
				} else {
					// sub
					t := baseTable[bits & ^msbWindow]
					t.Neg(&t)
					p.AddAssign(&t)
				}
			}

			// set our result point

			toReturn[i].FromJacobian(&p)

		}
	})

	return toReturn

}
