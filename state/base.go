package state

import (
	"fmt"
	"strings"

	"github.com/ProtoconNet/mitum2/base"
)

var PointPrefix = "point"

func StateKeyPointPrefix(contract base.Address) string {
	return fmt.Sprintf("%s:%s", PointPrefix, contract)
}

type StateKeyGenerator struct {
	contract base.Address
}

func NewStateKeyGenerator(contract base.Address) StateKeyGenerator {
	return StateKeyGenerator{
		contract,
	}
}

func (g StateKeyGenerator) Design() string {
	return StateKeyDesign(g.contract)
}

func (g StateKeyGenerator) PointBalance(address base.Address) string {
	return StateKeyPointBalance(g.contract, address)
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, PointPrefix) && strings.HasSuffix(key, DesignSuffix)
}

func IsStatePointBalanceKey(key string) bool {
	return strings.HasPrefix(key, PointPrefix) && strings.HasSuffix(key, PointBalanceSuffix)
}
