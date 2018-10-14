package portfolio

import (
	"errors"
	"sync"

	"github.com/nu11ptr/decimal"
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
	zero       = decimal.NewInt(0)
	one        = decimal.NewInt(1)
	oneHundred = decimal.NewInt(100)

	ErrBadSym      = errors.New("Symbol must be set to a valid value")
	ErrDupSym      = errors.New("Duplicate symbol")
	ErrSymNotFound = errors.New("The specified symbol could not be found")

	ErrBadPrice     = errors.New("Price must be greater than zero")
	ErrBadSecType   = errors.New("Invalid security type for the given symbol")
	ErrBadNumShares = errors.New("Actual positions require shares to be set")
	ErrBadPct       = errors.New("Percent must be set for a desired position")

	ErrPctOverflow  = errors.New("Total position percentage cannot exceed 100")
	ErrPctUnderflow = errors.New("Total position percentage must add up to 100")
)

// Position represents an actual or desired position in an account
type Position struct {
	Sym                string
	SecType            SecType
	Shares, Price, Pct decimal.Decimal
}

func (p *Position) validate(actual bool) error {
	if p.Sym == "" || (p.SecType == Cash && p.Sym != CashSym) {
		return ErrBadSym
	}
	if p.Sym == CashSym && p.SecType != Cash {
		return ErrBadSecType
	}
	if actual {
		if p.Shares.LTE(zero) {
			return ErrBadNumShares
		}
	} else {
		if p.Pct.LTE(zero) || p.Pct.GT(oneHundred) {
			return ErrBadPct
		}
	}
	return nil
}

// Account represents a brokerage account
type Account struct {
	actual, desired                                       map[string]Position
	balance, lmtPct, lmtOpenPct, lmtClosepct, rebalThresh decimal.Decimal
	mut                                                   sync.Mutex
	sellOnClose                                           bool

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
	totalPct := new(decimal.Decimal)

	for _, pos := range p {
		if err := pos.validate(actual); err != nil {
			return err
		}
		if _, ok := m[pos.Sym]; ok {
			return ErrDupSym
		}
		if !actual {
			totalPct = totalPct.Add(&pos.Pct)
			if totalPct.GT(oneHundred) {
				return ErrPctOverflow
			}
		}
		if pos.Sym == CashSym {
			pos.Price = *one
		}
		m[pos.Sym] = pos
	}
	if !actual && totalPct.LT(oneHundred) {
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

func copyMap(m map[string]Position) map[string]Position {
	m2 := make(map[string]Position, maxPos)
	for k, v := range m {
		m2[k] = v
	}
	return m2
}

// Actual returns a copy of the map storing actual positions
func (a *Account) Actual() map[string]Position {
	a.mut.Lock()
	defer a.mut.Unlock()

	return copyMap(a.actual)
}

// Desired returns a copy of the map storing the desired positions
func (a *Account) Desired() map[string]Position {
	a.mut.Lock()
	defer a.mut.Unlock()

	return copyMap(a.desired)
}

func setPrice(m map[string]Position, sym string, price *decimal.Decimal) bool {
	p, ok := m[sym]
	if ok {
		p.Price = *price
		m[sym] = p
		return true
	}
	return false
}

// SetPrice sets the price on the symbol specified. It returns an error if the price or symbol
// is invalid or if the symbol cannot be found
func (a *Account) SetPrice(sym string, price *decimal.Decimal) error {
	if sym == "" || sym == CashSym {
		return ErrBadSym
	}
	if price.LTE(zero) {
		return ErrBadPrice
	}
	a.mut.Lock()
	defer a.mut.Unlock()

	found := setPrice(a.actual, sym, price)
	found2 := setPrice(a.desired, sym, price)

	if !found && !found2 {
		return ErrSymNotFound
	}
	return nil
}

// SetPriceStr sets the price of the symbol specified using a string. It returns an error if the
// price or symbol is invalid or if the symbol can't be found
func (a *Account) SetPriceStr(sym, price string) error {
	d, ok := decimal.New(price)
	if !ok {
		return ErrBadPrice
	}
	return a.SetPrice(sym, d)
}
