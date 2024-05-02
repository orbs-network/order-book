package models

import (
	"fmt"
	"strings"
)

type Pair struct {
	aToken string
	bToken string
}

func NewPair(aToken, bToken string) *Pair {
	if len(aToken)*len(bToken) == 0 {
		return nil
	}
	return &Pair{
		aToken: strings.ToUpper(aToken),
		bToken: strings.ToUpper(bToken),
	}
}

func (p *Pair) String() string {
	return fmt.Sprintf("%s-%s", p.aToken, p.bToken)
}
func (p *Pair) Symbol() Symbol {
	return Symbol(p.String())
}

func (p *Pair) GetMakerSide(takerInToken string) Side {
	if takerInToken == p.aToken {
		return BUY
	} else {
		return SELL
	}
}

//////////////////////////////////////////////////////////////////////

type PairMngr struct {
	aToken2PairArr map[string][]*Pair
}

func NewPairMngr() *PairMngr {
	m := PairMngr{
		aToken2PairArr: make(map[string][]*Pair),
	}
	symbolPairs := GetAllSymbols()
	for _, sp := range symbolPairs {
		arr := strings.Split(sp.String(), "-")
		m.aToken2PairArr[arr[0]] = []*Pair{NewPair(arr[0], arr[1])}
	}
	return &m
}

func bTokenOfArr(arr []*Pair, bToken string) *Pair {
	for _, pair := range arr {
		if pair.bToken == bToken {
			return pair
		}
	}
	return nil
}

func (m *PairMngr) Resolve(xToken, yToken string) *Pair {
	// attempt x as aToken
	pairArr, ok := m.aToken2PairArr[xToken]
	if !ok {
		// attempt y as aToken
		pairArr, ok = m.aToken2PairArr[yToken]
		if !ok {
			return nil // pair not found for these two tokens
		}
		return bTokenOfArr(pairArr, xToken)
	}
	return bTokenOfArr(pairArr, yToken)
}
