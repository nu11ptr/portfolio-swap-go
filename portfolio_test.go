package portfolio_test

import (
	"math/big"
	"reflect"
	"testing"

	portfolio "github.com/nu11ptr/portfolio-swap"
)

var (
	badSym1    = portfolio.Position{Sym: ""}
	badSym2    = portfolio.Position{Sym: "bogus", SecType: portfolio.Cash}
	badSec1    = portfolio.Position{Sym: portfolio.CashSym, SecType: portfolio.Stock}
	badShares1 = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock, Shares: big.Rat{}}
	badShares2 = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock}
	badPct1    = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock, Pct: *big.NewRat(-1, 1)}
	badPct2    = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock, Pct: *big.NewRat(101, 1)}
	goodAct1   = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock, Shares: *big.NewRat(100, 1)}
	goodAct2   = portfolio.Position{Sym: "bogus2", SecType: portfolio.Stock, Shares: *big.NewRat(100, 1)}
	goodDes1   = portfolio.Position{Sym: "bogus", SecType: portfolio.Stock, Pct: *big.NewRat(40, 1)}
	goodDes2   = portfolio.Position{Sym: "bogus2", SecType: portfolio.Stock, Pct: *big.NewRat(60, 1)}
	goodDes3   = portfolio.Position{Sym: "bogus3", SecType: portfolio.Stock, Pct: *big.NewRat(1, 1)}
)

func TestSetActual(t *testing.T) {
	tests := []struct {
		name string
		p    []portfolio.Position
		err  error
	}{
		{"Blank", nil, nil},
		{"EmptySym", []portfolio.Position{badSym1}, portfolio.ErrBadSym},
		{"NonCashSym", []portfolio.Position{badSym2}, portfolio.ErrBadSym},
		{"NonCashType", []portfolio.Position{badSec1}, portfolio.ErrBadSecType},
		{"ZeroShares", []portfolio.Position{badShares1}, portfolio.ErrBadNumShares},
		{"DuplicateSym", []portfolio.Position{goodAct1, goodAct1}, portfolio.ErrDupSym},
		{"Good", []portfolio.Position{goodAct1, goodAct2}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			acct := portfolio.NewAccount(false, false)
			if err := acct.SetActual(test.p); err != test.err {
				t.Error("Got:", err, "Expected:", test.err)
			}
		})
	}
}

func TestSetDesired(t *testing.T) {
	tests := []struct {
		name string
		p    []portfolio.Position
		err  error
	}{
		{"EmptySym", []portfolio.Position{badSym1}, portfolio.ErrBadSym},
		{"NonCashSym", []portfolio.Position{badSym2}, portfolio.ErrBadSym},
		{"NonCashType", []portfolio.Position{badSec1}, portfolio.ErrBadSecType},
		{"NegPct", []portfolio.Position{badPct1}, portfolio.ErrBadPct},
		{"PctOverflow", []portfolio.Position{badPct2}, portfolio.ErrBadPct},
		{"DuplicateSym", []portfolio.Position{goodDes1, goodDes1}, portfolio.ErrDupSym},
		{"TotalPctOverflow", []portfolio.Position{goodDes1, goodDes2, goodDes3}, portfolio.ErrPctOverflow},
		{"TotalPctOverflow", []portfolio.Position{goodDes1}, portfolio.ErrPctUnderflow},
		{"Good", []portfolio.Position{goodDes1, goodDes2}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			acct := portfolio.NewAccount(false, false)
			if err := acct.SetDesired(test.p); err != test.err {
				t.Error("Got:", err, "Expected:", test.err)
			}
		})
	}
}

func TestActual(t *testing.T) {
	acct := portfolio.NewAccount(false, false)
	if err := acct.SetActual([]portfolio.Position{goodAct1, goodAct2}); err != nil {
		t.Error("Got:", err, "Expected:", nil)
	}

	expected := map[string]portfolio.Position{"bogus": goodAct1, "bogus2": goodAct2}
	actual := acct.Actual()
	if !reflect.DeepEqual(actual, expected) {
		t.Error("Got:", actual, "Expected:", expected)
	}
}

func TestDesired(t *testing.T) {
	acct := portfolio.NewAccount(false, false)
	if err := acct.SetDesired([]portfolio.Position{goodDes1, goodDes2}); err != nil {
		t.Error("Got:", err, "Expected:", nil)
	}

	expected := map[string]portfolio.Position{"bogus": goodDes1, "bogus2": goodDes2}
	actual := acct.Desired()
	if !reflect.DeepEqual(actual, expected) {
		t.Error("Got:", actual, "Expected:", expected)
	}
}

func TestSetPriceStr(t *testing.T) {
	acct := portfolio.NewAccount(false, false)
	if err := acct.SetDesired([]portfolio.Position{goodDes1, goodDes2}); err != nil {
		t.Error("Got:", err, "Expected:", nil)
	}

	tests := []struct {
		name, sym, price string
		err              error
	}{
		{"BadPrice", "bogus", "not_a_price", portfolio.ErrBadPrice},
		{"BadPrice2", "bogus", "-1.0", portfolio.ErrBadPrice},
		{"BadPrice3", "bogus", "0.0", portfolio.ErrBadPrice},
		{"BadSym1", "", "1.0", portfolio.ErrBadSym},
		{"BadSym2", portfolio.CashSym, "1.0", portfolio.ErrBadSym},
		{"SymNotFound", "bogus3", "1.0", portfolio.ErrSymNotFound},
		{"Good", "bogus", "1.0", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := acct.SetPriceStr(test.sym, test.price); err != test.err {
				t.Error("Got:", err, "Expected:", test.err)
			}

			// If we aren't looking for an error, validate the price was actually set
			if test.err == nil {
				expected := map[string]portfolio.Position{"bogus": goodDes1, "bogus2": goodDes2}
				// Manually copy the price into a new copy of the position and put back into the map
				p := expected[test.sym]
				r := new(big.Rat)
				pr, _ := r.SetString(test.price)
				p.Price = *pr
				expected[test.sym] = p

				actual := acct.Desired()
				if !reflect.DeepEqual(actual, expected) {
					t.Error("Got:", actual, "Expected:", expected)
				}
			}
		})
	}
}
