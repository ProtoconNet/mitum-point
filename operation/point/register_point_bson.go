package point

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/util"
)

func (fact RegisterPointFact) MarshalBSON() ([]byte, error) {
	m := fact.PointFact.marshalMap()

	m["symbol"] = fact.symbol
	m["name"] = fact.name
	m["initial_supply"] = fact.initialSupply

	return bsonenc.Marshal(m)
}

type RegisterPointFactBSONUnmarshaler struct {
	Symbol        string `bson:"symbol"`
	Name          string `bson:"name"`
	InitialSupply string `bson:"initial_supply"`
}

func (fact *RegisterPointFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*fact))

	if err := fact.PointFact.DecodeBSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf RegisterPointFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unpack(enc,
		uf.Symbol,
		uf.Name,
		uf.InitialSupply,
	)
}
