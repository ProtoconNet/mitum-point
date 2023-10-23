package state

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type TokenBalanceStateValueJSONMarshaler struct {
	hint.BaseHinter
	Amount common.Big `json:"amount"`
}

func (s PointBalanceStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TokenBalanceStateValueJSONMarshaler{
		BaseHinter: s.BaseHinter,
		Amount:     s.amount,
	})
}

type TokenBalanceStateValueJSONUnmarshaler struct {
	Amount string `json:"amount"`
}

func (s *PointBalanceStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u TokenBalanceStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	big, err := common.NewBigFromString(u.Amount)
	if err != nil {
		return e.Wrap(err)
	}
	s.amount = big

	return nil
}
