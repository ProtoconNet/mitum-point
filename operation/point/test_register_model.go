package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	pointtypes "github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum2/base"
)

type TestRegisterPointProcessor struct {
	*test.BaseTestOperationProcessorNoItem[RegisterModel]
}

func NewTestRegisterPointProcessor(tp *test.TestProcessor) TestRegisterPointProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[RegisterModel](tp)
	return TestRegisterPointProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestRegisterPointProcessor) Create() *TestRegisterPointProcessor {
	t.Opr, _ = NewRegisterModelProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestRegisterPointProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestRegisterPointProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestRegisterPointProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestRegisterPointProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestRegisterPointProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestRegisterPointProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestRegisterPointProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestRegisterPointProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestRegisterPointProcessor) LoadOperation(fileName string,
) *TestRegisterPointProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestRegisterPointProcessor) Print(fileName string,
) *TestRegisterPointProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestRegisterPointProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address,
	symbol pointtypes.PointSymbol, name string, initialSupply int64, currency types.CurrencyID,
) *TestRegisterPointProcessor {
	op := NewRegisterModel(
		NewRegisterModelFact(
			[]byte("Point"),
			sender,
			contract,
			currency,
			symbol,
			name,
			common.NewBig(initialSupply),
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestRegisterPointProcessor) RunPreProcess() *TestRegisterPointProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestRegisterPointProcessor) RunProcess() *TestRegisterPointProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestRegisterPointProcessor) IsValid() *TestRegisterPointProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestRegisterPointProcessor) Decode(fileName string) *TestRegisterPointProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
