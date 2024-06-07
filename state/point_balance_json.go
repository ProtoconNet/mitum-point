package state

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type PointBalanceStateValueJSONMarshaler struct {
	hint.BaseHinter
	Amount common.Big `json:"amount"`
}

func (s PointBalanceStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PointBalanceStateValueJSONMarshaler{
		BaseHinter: s.BaseHinter,
		Amount:     s.Amount,
	})
}

type PointBalanceStateValueJSONUnmarshaler struct {
	Amount string `json:"amount"`
}

func (s *PointBalanceStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u PointBalanceStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	big, err := common.NewBigFromString(u.Amount)
	if err != nil {
		return e.Wrap(err)
	}
	s.Amount = big

	return nil
}
