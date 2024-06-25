package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *RegisterModelFact) unpack(_ encoder.Encoder,
	symbol, name, decimal, initialSupply string,
) error {
	fact.symbol = types.PointSymbol(symbol)
	fact.name = name

	big, err := common.NewBigFromString(decimal)
	if err != nil {
		return err
	}
	fact.decimal = big

	big, err = common.NewBigFromString(initialSupply)
	if err != nil {
		return err
	}
	fact.initialSupply = big

	return nil
}
