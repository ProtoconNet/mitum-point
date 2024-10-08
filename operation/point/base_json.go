package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type PointFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender   base.Address             `json:"sender"`
	Contract base.Address             `json:"contract"`
	Currency currencytypes.CurrencyID `json:"currency"`
}

func (fact PointFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PointFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Currency:              fact.currency,
	})
}

type PointFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender   string `json:"sender"`
	Contract string `json:"contract"`
	Currency string `json:"currency"`
}

func (fact *PointFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf PointFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, uf.Sender, uf.Contract, uf.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (fact PointFact) JSONMarshaler() PointFactJSONMarshaler {
	return PointFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Currency:              fact.currency,
	}
}
