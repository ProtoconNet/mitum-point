package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *RegisterPointFact) unpack(enc encoder.Encoder,
	symbol, name, ts string,
) error {
	fact.symbol = types.PointID(symbol)
	fact.name = name

	big, err := common.NewBigFromString(ts)
	if err != nil {
		return err
	}
	fact.initialSupply = big

	return nil
}
