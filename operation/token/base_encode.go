package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *TokenFact) unmarshal(enc encoder.Encoder,
	sa, ca, tid, cid string,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*fact))

	fact.tokenID = types.CurrencyID(tid)
	fact.currency = types.CurrencyID(cid)

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.contract = a
	}

	return nil
}
