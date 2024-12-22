package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
)

func (fact TransferFromFact) MarshalBSON() ([]byte, error) {
	m := fact.PointFact.marshalMap()

	m["receiver"] = fact.receiver
	m["target"] = fact.target
	m["amount"] = fact.amount

	return bsonenc.Marshal(m)
}

type TransferFromFactBSONUnmarshaler struct {
	Receiver string `bson:"receiver"`
	Target   string `bson:"target"`
	Amount   string `bson:"amount"`
}

func (fact *TransferFromFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	if err := fact.PointFact.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	var uf TransferFromFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	if err := fact.unpack(enc, uf.Receiver, uf.Target, uf.Amount); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	return nil
}

func (op *TransferFrom) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *op)
	}

	op.BaseOperation = ubo

	var ueo extras.BaseOperationExtensions
	if err := ueo.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *op)
	}

	op.BaseOperationExtensions = &ueo

	return nil
}
