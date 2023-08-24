package token

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
)

func (fact TransferFact) MarshalBSON() ([]byte, error) {
	m := fact.TokenFact.marshalMap()

	m["receiver"] = fact.receiver
	m["amount"] = fact.amount

	return bsonenc.Marshal(m)
}

type TransferFactBSONUnmarshaler struct {
	Receiver string `bson:"receiver"`
	Amount   string `bson:"amount"`
}

func (fact *TransferFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*fact))

	if err := fact.TokenFact.DecodeBSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf TransferFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unmarshal(enc,
		uf.Receiver,
		uf.Amount,
	)
}
