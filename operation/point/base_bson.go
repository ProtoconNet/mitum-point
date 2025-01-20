package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

func (fact PointFact) marshalMap() map[string]interface{} {
	return map[string]interface{}{
		"_hint":    fact.Hint().String(),
		"sender":   fact.sender,
		"contract": fact.contract,
		"currency": fact.currency,
		"hash":     fact.BaseFact.Hash().String(),
		"token":    fact.BaseFact.Token(),
	}
}

func (fact PointFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    fact.Hint().String(),
			"sender":   fact.sender,
			"contract": fact.contract,
			"currency": fact.currency,
			"hash":     fact.BaseFact.Hash().String(),
			"token":    fact.BaseFact.Token(),
		},
	)
}

type PointFactBSONUnmarshaler struct {
	Hint     string `bson:"_hint"`
	Sender   string `bson:"sender"`
	Contract string `bson:"contract"`
	Currency string `bson:"currency"`
}

func (fact *PointFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubf common.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf PointFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	if err := fact.unpack(enc, uf.Sender, uf.Contract, uf.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	return nil
}
