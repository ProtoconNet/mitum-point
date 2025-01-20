package digest

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-point/state"
	"github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum2/base"
	utilm "github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DefaultColNamePoint        = "digest_point"
	DefaultColNamePointBalance = "digest_point_bl"
)

func Point(st *cdigest.Database, contract string) (*types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design *types.Design
	var sta base.State
	var err error
	if err := st.MongoClient().GetByFilter(
		DefaultColNamePoint,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}

			design, err = state.StateDesignValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, utilm.ErrNotFound.Errorf("point design, contract %s", contract)
	}

	return design, nil
}

func PointBalance(st *cdigest.Database, contract, account string) (*common.Big, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("address", account)

	var amount common.Big
	var sta base.State
	var err error
	if err := st.MongoClient().GetByFilter(
		DefaultColNamePointBalance,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}

			amount, err = state.StatePointBalanceValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		//return nil, mitumutil.ErrNotFound.Errorf("token balance by contract %s, account %s", contract, account)
	}

	return &amount, nil
}
