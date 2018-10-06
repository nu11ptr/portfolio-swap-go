package portfolio_test

import (
	"math/big"
	"testing"

	portfolio "github.com/nu11ptr/portfolio-swap"
)

var (
	badSym1    = portfolio.Position{Sym: ""}
	badSym2    = portfolio.Position{Sym: "bogus", SecType: portfolio.Cash}
	badSec1    = portfolio.Position{Sym: portfolio.CashSym, SecType: portfolio.Stock}
	badShares1 = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock, Shares: new(big.Rat)}
	badShares2 = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock}
	badPct1    = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock, Pct: big.NewRat(-1, 1)}
	badPct2    = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock, Pct: big.NewRat(101, 1)}
	goodAct1   = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock, Shares: big.NewRat(100, 1)}
	goodAct2   = portfolio.Position{Sym: "bogus2", SecType: portfolio.Stock, Shares: big.NewRat(100, 1)}
	goodPct1   = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock, Pct: big.NewRat(40, 1)}
	goodPct2   = portfolio.Position{Sym: "bogus2", SecType: portfolio.Stock, Pct: big.NewRat(60, 1)}
	goodPct3   = portfolio.Position{Sym: "bogus3", SecType: portfolio.Stock, Pct: big.NewRat(1, 1)}
)

func TestSetActual(t *testing.T) {
	tests := []struct {
		name string
		p    []portfolio.Position
		e    error
	}{
		{"Blank", nil, nil},
		{"EmptySym", []portfolio.Position{badSym1}, portfolio.ErrBadSym},
		{"NonCashSym", []portfolio.Position{badSym2}, portfolio.ErrBadSym},
		{"NonCashType", []portfolio.Position{badSec1}, portfolio.ErrBadSecType},
		{"ZeroShares", []portfolio.Position{badShares1}, portfolio.ErrBadNumShares},
		{"NilShares", []portfolio.Position{badShares2}, portfolio.ErrBadNumShares},
		{"DuplicateSym", []portfolio.Position{goodAct1, goodAct1}, portfolio.ErrDupSym},
		{"Good", []portfolio.Position{goodAct1, goodAct2}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			acct := portfolio.NewAccount(false, false)
			if err := acct.SetActual(test.p); err != test.e {
				t.Error("Got", err, "Expected", test.e)
			}
		})
	}
}

func TestSetDesired(t *testing.T) {
	tests := []struct {
		name string
		p    []portfolio.Position
		e    error
	}{
		{"EmptySym", []portfolio.Position{badSym1}, portfolio.ErrBadSym},
		{"NonCashSym", []portfolio.Position{badSym2}, portfolio.ErrBadSym},
		{"NonCashType", []portfolio.Position{badSec1}, portfolio.ErrBadSecType},
		{"NegPct", []portfolio.Position{badPct1}, portfolio.ErrBadPct},
		{"PctOverflow", []portfolio.Position{badPct2}, portfolio.ErrBadPct},
		{"DuplicateSym", []portfolio.Position{goodPct1, goodPct1}, portfolio.ErrDupSym},
		{"Good", []portfolio.Position{goodPct1, goodPct2}, nil},
		{"TotalPctOverflow", []portfolio.Position{goodPct1, goodPct2, goodPct3}, portfolio.ErrPctOverflow},
		{"TotalPctOverflow", []portfolio.Position{goodPct1}, portfolio.ErrPctUnderflow},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			acct := portfolio.NewAccount(false, false)
			if err := acct.SetDesired(test.p); err != test.e {
				t.Error("Got", err, "Expected", test.e)
			}
		})
	}
}
