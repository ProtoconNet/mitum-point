package digest

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-point/state"
	"github.com/ProtoconNet/mitum-point/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defaultColNameAccount         = "digest_ac"
	defaultColNameContractAccount = "digest_ca"
	defaultColNameBalance         = "digest_bl"
	defaultColNameCurrency        = "digest_cr"
	defaultColNameOperation       = "digest_op"
	defaultColNameBlock           = "digest_bm"
	defaultColNamePoint           = "digest_point"
	defaultColNamePointBalance    = "digest_point_bl"
)

func Point(st *currencydigest.Database, contract string) (*types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design *types.Design
	var sta mitumbase.State
	var err error
	if err := st.MongoClient().GetByFilter(
		defaultColNamePoint,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.Encoders())
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
		return nil, mitumutil.ErrNotFound.Errorf("point design, contract %s", contract)
	}

	return design, nil
}

func PointBalance(st *currencydigest.Database, contract, account string) (*common.Big, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("address", account)

	var amount common.Big
	var sta mitumbase.State
	var err error
	if err := st.MongoClient().GetByFilter(
		defaultColNamePointBalance,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.Encoders())
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
