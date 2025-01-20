package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"go.mongodb.org/mongo-driver/bson"
)

func (fact MintFact) MarshalBSON() ([]byte, error) {
	m := fact.PointFact.marshalMap()

	m["receiver"] = fact.receiver
	m["amount"] = fact.amount

	return bsonenc.Marshal(m)
}

type MintFactBSONUnmarshaler struct {
	Receiver string `bson:"receiver"`
	Amount   string `bson:"amount"`
}

func (fact *MintFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	if err := fact.PointFact.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	var uf MintFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	if err := fact.unpack(enc, uf.Receiver, uf.Amount); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	return nil
}

func (op *Mint) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
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
