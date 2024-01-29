package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type MintFactJSONMarshaler struct {
	PointFactJSONMarshaler
	Receiver base.Address `json:"receiver"`
	Amount   common.Big   `json:"amount"`
}

func (fact MintFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MintFactJSONMarshaler{
		PointFactJSONMarshaler: fact.PointFact.JSONMarshaler(),
		Receiver:               fact.receiver,
		Amount:                 fact.amount,
	})
}

type MintFactJSONUnMarshaler struct {
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (fact *MintFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*fact))

	if err := fact.PointFact.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf MintFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unpack(enc,
		uf.Receiver,
		uf.Amount,
	)
}
