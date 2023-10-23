package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type BurnFactJSONMarshaler struct {
	PointFactJSONMarshaler
	Target base.Address `json:"target"`
	Amount common.Big   `json:"amount"`
}

func (fact BurnFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(BurnFactJSONMarshaler{
		PointFactJSONMarshaler: fact.PointFact.JSONMarshaler(),
		Target:                 fact.target,
		Amount:                 fact.amount,
	})
}

type BurnFactJSONUnMarshaler struct {
	Target string `json:"target"`
	Amount string `json:"amount"`
}

func (fact *BurnFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*fact))

	if err := fact.PointFact.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf BurnFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unpack(enc,
		uf.Target,
		uf.Amount,
	)
}
