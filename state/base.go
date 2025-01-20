package state

import (
	"fmt"
	"strings"
)

var PointPrefix = "point"

func StateKeyPointPrefix(contract string) string {
	return fmt.Sprintf("%s:%s", PointPrefix, contract)
}

type StateKeyGenerator struct {
	contract string
}

func NewStateKeyGenerator(contract string) StateKeyGenerator {
	return StateKeyGenerator{
		contract,
	}
}

func (g StateKeyGenerator) Design() string {
	return StateKeyDesign(g.contract)
}

func (g StateKeyGenerator) PointBalance(address string) string {
	return StateKeyPointBalance(g.contract, address)
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, PointPrefix) && strings.HasSuffix(key, DesignSuffix)
}

func IsStatePointBalanceKey(key string) bool {
	return strings.HasPrefix(key, PointPrefix) && strings.HasSuffix(key, PointBalanceSuffix)
}
