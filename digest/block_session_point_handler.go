package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-point/state"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func PreparePoint(bs *currencydigest.BlockSession, st base.State) (string, []mongo.WriteModel, error) {
	switch {
	case state.IsStateDesignKey(st.Key()):
		j, err := handlePointState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNamePoint, j, nil
	case state.IsStatePointBalanceKey(st.Key()):
		j, err := handlePointBalanceState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNamePointBalance, j, nil
	}

	return "", nil, nil
}

func handlePointState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if pointDoc, err := NewPointDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(pointDoc),
		}, nil
	}
}

func handlePointBalanceState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if pointBalanceDoc, err := NewPointBalanceDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(pointBalanceDoc),
		}, nil
	}
}
