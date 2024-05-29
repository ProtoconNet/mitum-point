package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type RegisterPointFactJSONMarshaler struct {
	PointFactJSONMarshaler
	Symbol        types.PointID `json:"symbol"`
	Name          string        `json:"name"`
	InitialSupply common.Big    `json:"initial_supply"`
}

func (fact RegisterPointFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterPointFactJSONMarshaler{
		PointFactJSONMarshaler: fact.PointFact.JSONMarshaler(),
		Symbol:                 fact.symbol,
		Name:                   fact.name,
		InitialSupply:          fact.initialSupply,
	})
}

type RegisterPointFactJSONUnMarshaler struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	InitialSupply string `json:"initial_supply"`
}

func (fact *RegisterPointFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	if err := fact.PointFact.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	var uf RegisterPointFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	if err := fact.unpack(enc, uf.Symbol, uf.Name, uf.InitialSupply); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}
