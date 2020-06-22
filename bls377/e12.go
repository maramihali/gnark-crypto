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

// Code generated by gurvy/internal/generators DO NOT EDIT

package bls377

import (
	"github.com/consensys/gurvy/bls377/fp"
)

// E12 is a degree-two finite field extension of fp6:
// C0 + C1w where w^3-v is irrep in fp6

// fp2, fp12 are both quadratic field extensions
// template code is duplicated in fp2, fp12
// TODO make an abstract quadratic extension template

type E12 struct {
	C0, C1 E6
}

// Equal returns true if z equals x, fasle otherwise
// TODO can this be deleted?  Should be able to use == operator instead
func (z *E12) Equal(x *E12) bool {
	return z.C0.Equal(&x.C0) && z.C1.Equal(&x.C1)
}

// String puts E12 in string form
func (z *E12) String() string {
	return (z.C0.String() + "+(" + z.C1.String() + ")*w")
}

// SetString sets a E12 from string
func (z *E12) SetString(s0, s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11 string) *E12 {
	z.C0.SetString(s0, s1, s2, s3, s4, s5)
	z.C1.SetString(s6, s7, s8, s9, s10, s11)
	return z
}

// Set copies x into z and returns z
func (z *E12) Set(x *E12) *E12 {
	z.C0 = x.C0
	z.C1 = x.C1
	return z
}

// SetOne sets z to 1 in Montgomery form and returns z
func (z *E12) SetOne() *E12 {
	z.C0.B0.A0.SetOne()
	z.C0.B0.A1.SetZero()
	z.C0.B1.A0.SetZero()
	z.C0.B1.A1.SetZero()
	z.C0.B2.A0.SetZero()
	z.C0.B2.A1.SetZero()
	z.C1.B0.A0.SetZero()
	z.C1.B0.A1.SetZero()
	z.C1.B1.A0.SetZero()
	z.C1.B1.A1.SetZero()
	z.C1.B2.A0.SetZero()
	z.C1.B2.A1.SetZero()
	return z
}

// ToMont converts to Mont form
// TODO can this be deleted?
func (z *E12) ToMont() *E12 {
	z.C0.ToMont()
	z.C1.ToMont()
	return z
}

// FromMont converts from Mont form
// TODO can this be deleted?
func (z *E12) FromMont() *E12 {
	z.C0.FromMont()
	z.C1.FromMont()
	return z
}

// Add set z=x+y in E12 and return z
func (z *E12) Add(x, y *E12) *E12 {
	z.C0.Add(&x.C0, &y.C0)
	z.C1.Add(&x.C1, &y.C1)
	return z
}

// Sub set z=x-y in E12 and return z
func (z *E12) Sub(x, y *E12) *E12 {
	z.C0.Sub(&x.C0, &y.C0)
	z.C1.Sub(&x.C1, &y.C1)
	return z
}

// SetRandom used only in tests
// TODO eliminate this method!
func (z *E12) SetRandom() *E12 {
	z.C0.B0.A0.SetRandom()
	z.C0.B0.A1.SetRandom()
	z.C0.B1.A0.SetRandom()
	z.C0.B1.A1.SetRandom()
	z.C0.B2.A0.SetRandom()
	z.C0.B2.A1.SetRandom()
	z.C1.B0.A0.SetRandom()
	z.C1.B0.A1.SetRandom()
	z.C1.B1.A0.SetRandom()
	z.C1.B1.A1.SetRandom()
	z.C1.B2.A0.SetRandom()
	z.C1.B2.A1.SetRandom()
	return z
}

// Mul set z=x*y in E12 and return z
func (z *E12) Mul(x, y *E12) *E12 {
	// Algorithm 20 from https://eprint.iacr.org/2010/354.pdf

	var t0, t1, xSum, ySum E6

	t0.Mul(&x.C0, &y.C0) // step 1
	t1.Mul(&x.C1, &y.C1) // step 2

	// finish processing input in case z==x or y
	xSum.Add(&x.C0, &x.C1)
	ySum.Add(&y.C0, &y.C1)

	// step 3
	{ // begin inline: set z.C0 to (&t1) * ((0,0),(1,0),(0,0))
		var result E6
		result.B1.Set(&(&t1).B0)
		result.B2.Set(&(&t1).B1)
		{ // begin inline: set result.B0 to (&(&t1).B2) * (0,1)
			buf := (&(&t1).B2).A0
			{ // begin inline: set &(result.B0).A0 to (&(&(&t1).B2).A1) * (5)
				buf := *(&(&(&t1).B2).A1)
				(&(result.B0).A0).Double(&buf).Double(&(result.B0).A0).AddAssign(&buf)
			} // end inline: set &(result.B0).A0 to (&(&(&t1).B2).A1) * (5)
			(result.B0).A1 = buf
		} // end inline: set result.B0 to (&(&t1).B2) * (0,1)
		z.C0.Set(&result)
	} // end inline: set z.C0 to (&t1) * ((0,0),(1,0),(0,0))
	z.C0.Add(&z.C0, &t0)

	// step 4
	z.C1.Mul(&xSum, &ySum).
		Sub(&z.C1, &t0).
		Sub(&z.C1, &t1)

	return z
}

// Square set z=x*x in E12 and return z
func (z *E12) Square(x *E12) *E12 {
	// TODO implement Algorithm 22 from https://eprint.iacr.org/2010/354.pdf
	// or the complex method from fp2
	// for now do it the dumb way
	var b0, b1 E6

	b0.Square(&x.C0)
	b1.Square(&x.C1)
	{ // begin inline: set b1 to (&b1) * ((0,0),(1,0),(0,0))
		var result E6
		result.B1.Set(&(&b1).B0)
		result.B2.Set(&(&b1).B1)
		{ // begin inline: set result.B0 to (&(&b1).B2) * (0,1)
			buf := (&(&b1).B2).A0
			{ // begin inline: set &(result.B0).A0 to (&(&(&b1).B2).A1) * (5)
				buf := *(&(&(&b1).B2).A1)
				(&(result.B0).A0).Double(&buf).Double(&(result.B0).A0).AddAssign(&buf)
			} // end inline: set &(result.B0).A0 to (&(&(&b1).B2).A1) * (5)
			(result.B0).A1 = buf
		} // end inline: set result.B0 to (&(&b1).B2) * (0,1)
		b1.Set(&result)
	} // end inline: set b1 to (&b1) * ((0,0),(1,0),(0,0))
	b1.Add(&b0, &b1)

	z.C1.Mul(&x.C0, &x.C1).Double(&z.C1)
	z.C0 = b1

	return z
}

// Inverse set z to the inverse of x in E12 and return z
func (z *E12) Inverse(x *E12) *E12 {
	// Algorithm 23 from https://eprint.iacr.org/2010/354.pdf

	var t [2]E6

	t[0].Square(&x.C0) // step 1
	t[1].Square(&x.C1) // step 2
	{                  // step 3
		var buf E6
		{ // begin inline: set buf to (&t[1]) * ((0,0),(1,0),(0,0))
			var result E6
			result.B1.Set(&(&t[1]).B0)
			result.B2.Set(&(&t[1]).B1)
			{ // begin inline: set result.B0 to (&(&t[1]).B2) * (0,1)
				buf := (&(&t[1]).B2).A0
				{ // begin inline: set &(result.B0).A0 to (&(&(&t[1]).B2).A1) * (5)
					buf := *(&(&(&t[1]).B2).A1)
					(&(result.B0).A0).Double(&buf).Double(&(result.B0).A0).AddAssign(&buf)
				} // end inline: set &(result.B0).A0 to (&(&(&t[1]).B2).A1) * (5)
				(result.B0).A1 = buf
			} // end inline: set result.B0 to (&(&t[1]).B2) * (0,1)
			buf.Set(&result)
		} // end inline: set buf to (&t[1]) * ((0,0),(1,0),(0,0))
		t[0].Sub(&t[0], &buf)
	}
	t[1].Inverse(&t[0])               // step 4
	z.C0.Mul(&x.C0, &t[1])            // step 5
	z.C1.Mul(&x.C1, &t[1]).Neg(&z.C1) // step 6

	return z
}

// InverseUnitary inverse a unitary element
// TODO deprecate in favour of Conjugate
func (z *E12) InverseUnitary(x *E12) *E12 {
	return z.Conjugate(x)
}

// Conjugate set z to (x.C0, -x.C1) and return z
func (z *E12) Conjugate(x *E12) *E12 {
	z.Set(x)
	z.C1.Neg(&z.C1)
	return z
}

// MulByVW set z to x*(y*v*w) and return z
// here y*v*w means the E12 element with C1.B1=y and all other components 0
func (z *E12) MulByVW(x *E12, y *E2) *E12 {
	var result E12
	var yNR E2

	{ // begin inline: set yNR to (y) * (0,1)
		buf := (y).A0
		{ // begin inline: set &(yNR).A0 to (&(y).A1) * (5)
			buf := *(&(y).A1)
			(&(yNR).A0).Double(&buf).Double(&(yNR).A0).AddAssign(&buf)
		} // end inline: set &(yNR).A0 to (&(y).A1) * (5)
		(yNR).A1 = buf
	} // end inline: set yNR to (y) * (0,1)
	result.C0.B0.Mul(&x.C1.B1, &yNR)
	result.C0.B1.Mul(&x.C1.B2, &yNR)
	result.C0.B2.Mul(&x.C1.B0, y)
	result.C1.B0.Mul(&x.C0.B2, &yNR)
	result.C1.B1.Mul(&x.C0.B0, y)
	result.C1.B2.Mul(&x.C0.B1, y)
	z.Set(&result)
	return z
}

// MulByV set z to x*(y*v) and return z
// here y*v means the E12 element with C0.B1=y and all other components 0
func (z *E12) MulByV(x *E12, y *E2) *E12 {
	var result E12
	var yNR E2

	{ // begin inline: set yNR to (y) * (0,1)
		buf := (y).A0
		{ // begin inline: set &(yNR).A0 to (&(y).A1) * (5)
			buf := *(&(y).A1)
			(&(yNR).A0).Double(&buf).Double(&(yNR).A0).AddAssign(&buf)
		} // end inline: set &(yNR).A0 to (&(y).A1) * (5)
		(yNR).A1 = buf
	} // end inline: set yNR to (y) * (0,1)
	result.C0.B0.Mul(&x.C0.B2, &yNR)
	result.C0.B1.Mul(&x.C0.B0, y)
	result.C0.B2.Mul(&x.C0.B1, y)
	result.C1.B0.Mul(&x.C1.B2, &yNR)
	result.C1.B1.Mul(&x.C1.B0, y)
	result.C1.B2.Mul(&x.C1.B1, y)
	z.Set(&result)
	return z
}

// MulByV2W set z to x*(y*v^2*w) and return z
// here y*v^2*w means the E12 element with C1.B2=y and all other components 0
func (z *E12) MulByV2W(x *E12, y *E2) *E12 {
	var result E12
	var yNR E2

	{ // begin inline: set yNR to (y) * (0,1)
		buf := (y).A0
		{ // begin inline: set &(yNR).A0 to (&(y).A1) * (5)
			buf := *(&(y).A1)
			(&(yNR).A0).Double(&buf).Double(&(yNR).A0).AddAssign(&buf)
		} // end inline: set &(yNR).A0 to (&(y).A1) * (5)
		(yNR).A1 = buf
	} // end inline: set yNR to (y) * (0,1)
	result.C0.B0.Mul(&x.C1.B0, &yNR)
	result.C0.B1.Mul(&x.C1.B1, &yNR)
	result.C0.B2.Mul(&x.C1.B2, &yNR)
	result.C1.B0.Mul(&x.C0.B1, &yNR)
	result.C1.B1.Mul(&x.C0.B2, &yNR)
	result.C1.B2.Mul(&x.C0.B0, y)
	z.Set(&result)
	return z
}

// MulByV2NRInv set z to x*(y*v^2*(0,1)^{-1}) and return z
// here y*v^2 means the E12 element with C0.B2=y and all other components 0
func (z *E12) MulByV2NRInv(x *E12, y *E2) *E12 {
	var result E12
	var yNRInv E2

	{ // begin inline: set yNRInv to (y) * (0,1)^{-1}
		buf := (y).A1
		{ // begin inline: set &(yNRInv).A1 to (&(y).A0) * (5)^{-1}
			nrinv := fp.Element{
				330620507644336508,
				9878087358076053079,
				11461392860540703536,
				6973035786057818995,
				8846909097162646007,
				104838758629667239,
			}
			(&(yNRInv).A1).Mul(&(y).A0, &nrinv)
		} // end inline: set &(yNRInv).A1 to (&(y).A0) * (5)^{-1}
		(yNRInv).A0 = buf
	} // end inline: set yNRInv to (y) * (0,1)^{-1}

	result.C0.B0.Mul(&x.C0.B1, y)
	result.C0.B1.Mul(&x.C0.B2, y)
	result.C0.B2.Mul(&x.C0.B0, &yNRInv)

	result.C1.B0.Mul(&x.C1.B1, y)
	result.C1.B1.Mul(&x.C1.B2, y)
	result.C1.B2.Mul(&x.C1.B0, &yNRInv)

	z.Set(&result)
	return z
}

// MulByVWNRInv set z to x*(y*v*w*(0,1)^{-1}) and return z
// here y*v*w means the E12 element with C1.B1=y and all other components 0
func (z *E12) MulByVWNRInv(x *E12, y *E2) *E12 {
	var result E12
	var yNRInv E2

	{ // begin inline: set yNRInv to (y) * (0,1)^{-1}
		buf := (y).A1
		{ // begin inline: set &(yNRInv).A1 to (&(y).A0) * (5)^{-1}
			nrinv := fp.Element{
				330620507644336508,
				9878087358076053079,
				11461392860540703536,
				6973035786057818995,
				8846909097162646007,
				104838758629667239,
			}
			(&(yNRInv).A1).Mul(&(y).A0, &nrinv)
		} // end inline: set &(yNRInv).A1 to (&(y).A0) * (5)^{-1}
		(yNRInv).A0 = buf
	} // end inline: set yNRInv to (y) * (0,1)^{-1}

	result.C0.B0.Mul(&x.C1.B1, y)
	result.C0.B1.Mul(&x.C1.B2, y)
	result.C0.B2.Mul(&x.C1.B0, &yNRInv)

	result.C1.B0.Mul(&x.C0.B2, y)
	result.C1.B1.Mul(&x.C0.B0, &yNRInv)
	result.C1.B2.Mul(&x.C0.B1, &yNRInv)

	z.Set(&result)
	return z
}

// MulByWNRInv set z to x*(y*w*(0,1)^{-1}) and return z
// here y*w means the E12 element with C1.B0=y and all other components 0
func (z *E12) MulByWNRInv(x *E12, y *E2) *E12 {
	var result E12
	var yNRInv E2

	{ // begin inline: set yNRInv to (y) * (0,1)^{-1}
		buf := (y).A1
		{ // begin inline: set &(yNRInv).A1 to (&(y).A0) * (5)^{-1}
			nrinv := fp.Element{
				330620507644336508,
				9878087358076053079,
				11461392860540703536,
				6973035786057818995,
				8846909097162646007,
				104838758629667239,
			}
			(&(yNRInv).A1).Mul(&(y).A0, &nrinv)
		} // end inline: set &(yNRInv).A1 to (&(y).A0) * (5)^{-1}
		(yNRInv).A0 = buf
	} // end inline: set yNRInv to (y) * (0,1)^{-1}

	result.C0.B0.Mul(&x.C1.B2, y)
	result.C0.B1.Mul(&x.C1.B0, &yNRInv)
	result.C0.B2.Mul(&x.C1.B1, &yNRInv)

	result.C1.B0.Mul(&x.C0.B0, &yNRInv)
	result.C1.B1.Mul(&x.C0.B1, &yNRInv)
	result.C1.B2.Mul(&x.C0.B2, &yNRInv)

	z.Set(&result)
	return z
}

// MulByNonResidue multiplies a E6 by ((0,0),(1,0),(0,0))
// TODO delete this method once you have another way of testing the inlined code
func (z *E6) MulByNonResidue(x *E6) *E6 {
	{ // begin inline: set z to (x) * ((0,0),(1,0),(0,0))
		var result E6
		result.B1.Set(&(x).B0)
		result.B2.Set(&(x).B1)
		{ // begin inline: set result.B0 to (&(x).B2) * (0,1)
			buf := (&(x).B2).A0
			{ // begin inline: set &(result.B0).A0 to (&(&(x).B2).A1) * (5)
				buf := *(&(&(x).B2).A1)
				(&(result.B0).A0).Double(&buf).Double(&(result.B0).A0).AddAssign(&buf)
			} // end inline: set &(result.B0).A0 to (&(&(x).B2).A1) * (5)
			(result.B0).A1 = buf
		} // end inline: set result.B0 to (&(x).B2) * (0,1)
		z.Set(&result)
	} // end inline: set z to (x) * ((0,0),(1,0),(0,0))
	return z
}
