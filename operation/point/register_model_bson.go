package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"go.mongodb.org/mongo-driver/bson"
)

func (fact RegisterModelFact) MarshalBSON() ([]byte, error) {
	m := fact.PointFact.marshalMap()

	m["symbol"] = fact.symbol
	m["name"] = fact.name
	m["decimal"] = fact.decimal
	m["initial_supply"] = fact.initialSupply

	return bsonenc.Marshal(m)
}

type RegisterModelFactBSONUnmarshaler struct {
	Symbol        string `bson:"symbol"`
	Name          string `bson:"name"`
	Decimal       string `bson:"decimal"`
	InitialSupply string `bson:"initial_supply"`
}

func (fact *RegisterModelFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	if err := fact.PointFact.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	var uf RegisterModelFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	if err := fact.unpack(enc, uf.Symbol, uf.Name, uf.Decimal, uf.InitialSupply); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	return nil
}

func (op *RegisterModel) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
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
