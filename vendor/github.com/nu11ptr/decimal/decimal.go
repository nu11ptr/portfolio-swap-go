package decimal

import (
	"log"
	"math/big"
)

// Decimal represents a large decimal number
type Decimal big.Rat

// New creates a new decimal from the given string. Ok is true if string is valid otherwise it is false
func New(s string) (*Decimal, bool) {
	d := new(Decimal)
	if _, ok := (*big.Rat)(d).SetString(s); !ok {
		return nil, false
	}
	return d, true
}

// NewMust creates a new decimal from the given string panicking if the string is invalid
func NewMust(s string) *Decimal {
	d, ok := New(s)
	if !ok {
		log.Panicf("Invalid decimal string value (%s)", s)
	}
	return d
}

// NewInt creates a new decimal by setting a single integer as the value with all zeros after the decimal
func NewInt(i int) *Decimal {
	d := new(Decimal)
	(*big.Rat)(d).SetFrac64(int64(i), 1)
	return d
}

// Strings returns a string representation of this decimal with the specified precision
func (d *Decimal) String(prec int) string {
	return (*big.Rat)(d).FloatString(prec)
}

// LT returns true if d is less than d2 otherwise false
func (d *Decimal) LT(d2 *Decimal) bool {
	return (*big.Rat)(d).Cmp((*big.Rat)(d2)) == -1
}

// LTE returns true if d is less than or equal to d2 otherwise false
func (d *Decimal) LTE(d2 *Decimal) bool {
	return (*big.Rat)(d).Cmp((*big.Rat)(d2)) < 1
}

// GT returns true if d is greater than d2 otherwise false
func (d *Decimal) GT(d2 *Decimal) bool {
	return (*big.Rat)(d).Cmp((*big.Rat)(d2)) == 1
}

// GTE returns true if d is greater than or equal to d2 otherwise false
func (d *Decimal) GTE(d2 *Decimal) bool {
	return (*big.Rat)(d).Cmp((*big.Rat)(d2)) > -1
}

// Eq returns true if d is equal to d2 otherwise false
func (d *Decimal) Eq(d2 *Decimal) bool {
	return (*big.Rat)(d).Cmp((*big.Rat)(d2)) == 0
}

// Add adds d2 to d and returns the result
func (d *Decimal) Add(d2 *Decimal) *Decimal {
	d3 := new(Decimal)
	(*big.Rat)(d3).Add((*big.Rat)(d), (*big.Rat)(d2))
	return d3
}

// Sub subtracts d2 from d and returns the result
func (d *Decimal) Sub(d2 *Decimal) *Decimal {
	d3 := new(Decimal)
	(*big.Rat)(d3).Sub((*big.Rat)(d), (*big.Rat)(d2))
	return d3
}

// Mult multiplies d times d2 and returns the result
func (d *Decimal) Mult(d2 *Decimal) *Decimal {
	d3 := new(Decimal)
	(*big.Rat)(d3).Mul((*big.Rat)(d), (*big.Rat)(d2))
	return d3
}

// Div divides d by d2 and returns the result
func (d *Decimal) Div(d2 *Decimal) *Decimal {
	d3 := new(Decimal)
	(*big.Rat)(d3).Quo((*big.Rat)(d), (*big.Rat)(d2))
	return d3
}
