package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *RegisterPointFact) unpack(enc encoder.Encoder,
	symbol, name, ts string,
) error {
	e := util.StringError(utils.ErrStringUnPack(*fact))

	fact.symbol = currencytypes.CurrencyID(symbol)
	fact.name = name

	big, err := common.NewBigFromString(ts)
	if err != nil {
		return e.Wrap(err)
	}
	fact.initialSupply = big

	return nil
}
