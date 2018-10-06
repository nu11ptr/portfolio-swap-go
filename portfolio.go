package portfolio

import (
	"errors"
	"math/big"
	"sync"
)

// SecType represents a security type
type SecType int

const (
	maxPos = 16

	// CashSym represents a cash position
	CashSym = "*CASH*"

	// Stock represents a stock
	Stock SecType = iota
	// Fund represents a mutual fund
	Fund
	// Cash represents a cash position
	Cash
)

var (
	zero       = new(big.Rat)
	oneHundred = big.NewRat(100, 1)

	ErrBadSym       = errors.New("Symbol must be set to a valid value")
	ErrBadSecType   = errors.New("Invalid security type for the given symbol")
	ErrBadNumShares = errors.New("Actual positions require shares to be set")
	ErrBadPct       = errors.New("Percent must be set for a desired position")
	ErrDupSym       = errors.New("Duplicate symbol")
	ErrPctOverflow  = errors.New("Total position percentage cannot exceed 100")
	ErrPctUnderflow = errors.New("Total position percentage must add up to 100")
)

// Position represents an actual or desired position in an account
type Position struct {
	Sym                string
	SecType            SecType
	Shares, Price, Pct big.Rat
}

func (p *Position) validate(actual bool) error {
	if p.Sym == "" || (p.SecType == Cash && p.Sym != CashSym) {
		return ErrBadSym
	}
	if p.Sym == CashSym && p.SecType != Cash {
		return ErrBadSecType
	}
	if actual {
		if p.Shares.Cmp(zero) < 1 {
			return ErrBadNumShares
		}
	} else {
		if p.Pct.Cmp(zero) < 1 || p.Pct.Cmp(oneHundred) > 0 {
			return ErrBadPct
		}
	}
	return nil
}

// Account represents a brokerage account
type Account struct {
	actual, desired map[string]Position
	mut             sync.Mutex

	Margin, NonTaxable bool
}

// NewAccount creates a new account
func NewAccount(margin, nonTaxable bool) *Account {
	return &Account{
		actual: make(map[string]Position, maxPos), desired: make(map[string]Position, maxPos),
		Margin: margin, NonTaxable: nonTaxable,
	}
}

func setPositions(m map[string]Position, p []Position, actual bool) error {
	totalPct := new(big.Rat)

	for _, pos := range p {
		if err := pos.validate(actual); err != nil {
			return err
		}
		if _, ok := m[pos.Sym]; ok {
			return ErrDupSym
		}
		if !actual {
			totalPct.Add(totalPct, &pos.Pct)
			if totalPct.Cmp(oneHundred) == 1 {
				return ErrPctOverflow
			}
		}
		m[pos.Sym] = pos
	}
	if !actual && totalPct.Cmp(oneHundred) == -1 {
		return ErrPctUnderflow
	}
	return nil
}

// SetActual sets the actual set of positions for the account
func (a *Account) SetActual(p []Position) error {
	a.mut.Lock()
	defer a.mut.Unlock()

	a.actual = make(map[string]Position, maxPos)
	return setPositions(a.actual, p, true)
}

// SetDesired sets the desired sets of positions for the account
func (a *Account) SetDesired(p []Position) error {
	a.mut.Lock()
	defer a.mut.Unlock()

	a.desired = make(map[string]Position, maxPos)
	return setPositions(a.desired, p, false)
}
