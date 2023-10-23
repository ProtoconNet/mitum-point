package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type RegisterPointFactJSONMarshaler struct {
	PointFactJSONMarshaler
	Symbol        currencytypes.CurrencyID `json:"symbol"`
	Name          string                   `json:"name"`
	InitialSupply common.Big               `json:"initial_supply"`
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

func (fact *RegisterPointFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*fact))

	if err := fact.PointFact.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf RegisterPointFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unpack(enc,
		uf.Symbol,
		uf.Name,
		uf.InitialSupply,
	)
}
