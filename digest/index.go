package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var pointServiceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "point_service_contract_height"),
	},
}

var pointBalanceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "address", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "point_balance_contract_address_height"),
	},
}

var DefaultIndexes = cdigest.DefaultIndexes

func init() {
	DefaultIndexes[DefaultColNamePoint] = pointServiceIndexModels
	DefaultIndexes[DefaultColNamePointBalance] = pointBalanceIndexModels
}
