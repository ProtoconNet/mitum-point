package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (d *Design) unpack(enc encoder.Encoder, ht hint.Hint, symbol, name string, bp []byte) error {
	e := util.StringError(utils.ErrStringUnPack(*d))

	d.BaseHinter = hint.NewBaseHinter(ht)
	d.symbol = types.CurrencyID(symbol)
	d.name = name

	if hinter, err := enc.Decode(bp); err != nil {
		return e.Wrap(err)
	} else if p, ok := hinter.(Policy); !ok {
		return e.Wrap(errors.Errorf(utils.ErrStringTypeCast(Policy{}, hinter)))
	} else {
		d.policy = p
	}

	return nil
}
