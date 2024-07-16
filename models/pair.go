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
	token2PairArr map[string][]*Pair
	// bToken2PairArr map[string][]*Pair
}

func NewPairMngr() *PairMngr {
	m := PairMngr{
		token2PairArr: make(map[string][]*Pair),
	}
	symbolPairs := GetAllSymbols()
	for _, sp := range symbolPairs {
		arr := strings.Split(sp.String(), "-")
		aToken := arr[0]
		bToken := arr[1]
		pair := NewPair(aToken, bToken)
		// A token map
		if len(m.token2PairArr[aToken]) == 0 {
			m.token2PairArr[aToken] = []*Pair{}
		}
		m.token2PairArr[aToken] = append(m.token2PairArr[aToken], pair)
		// B token map
		if len(m.token2PairArr[bToken]) == 0 {
			m.token2PairArr[bToken] = []*Pair{}
		}
		m.token2PairArr[bToken] = append(m.token2PairArr[bToken], pair)

		fmt.Println("aToken2PairArr pair added: " + sp)
	}
	return &m
}

func findPair(arr []*Pair, outToken string) *Pair {
	for _, pair := range arr {
		if pair.aToken == outToken || pair.bToken == outToken {
			return pair
		}
	}
	return nil
}

func (m *PairMngr) Resolve(inToken, outToken string) *Pair {
	if inToken == outToken {
		return nil // illegal pair same token
	}
	// attempt inToken as aToken
	pairArr, ok := m.token2PairArr[inToken]
	if !ok {
		return nil // pair not found for these two tokens
	}
	return findPair(pairArr, outToken)
}
