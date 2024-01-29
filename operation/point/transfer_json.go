package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type TransferFactJSONMarshaler struct {
	PointFactJSONMarshaler
	Receiver base.Address `json:"receiver"`
	Amount   common.Big   `json:"amount"`
}

func (fact TransferFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferFactJSONMarshaler{
		PointFactJSONMarshaler: fact.PointFact.JSONMarshaler(),
		Receiver:               fact.receiver,
		Amount:                 fact.amount,
	})
}

type TransferFactJSONUnMarshaler struct {
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (fact *TransferFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*fact))

	if err := fact.PointFact.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf TransferFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unpack(enc,
		uf.Receiver,
		uf.Amount,
	)
}
