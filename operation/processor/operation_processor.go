package processor

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/operation/point"
	"github.com/ProtoconNet/mitum-point/utils"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

const (
	DuplicationTypeSender   currencytypes.DuplicationType = "sender"
	DuplicationTypeCurrency currencytypes.DuplicationType = "currency"
)

func CheckDuplication(opr *currencyprocessor.OperationProcessor, op mitumbase.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	var did string
	var didtype currencytypes.DuplicationType
	var newAddresses []mitumbase.Address
	var err error

	switch t := op.(type) {
	case currency.CreateAccount:
		fact, ok := t.Fact().(currency.CreateAccountFact)
		if !ok {
			return errors.Errorf("expected CreateAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case currency.UpdateKey:
		fact, ok := t.Fact().(currency.UpdateKeyFact)
		if !ok {
			return errors.Errorf("expected UpdateKeyFact, not %T", t.Fact())
		}
		did = fact.Target().String()
		didtype = DuplicationTypeSender
	case currency.Transfer:
		fact, ok := t.Fact().(currency.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case extension.CreateContractAccount:
		fact, ok := t.Fact().(extension.CreateContractAccountFact)
		if !ok {
			return errors.Errorf("expected CreateContractAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
	case extension.Withdraw:
		fact, ok := t.Fact().(extension.WithdrawFact)
		if !ok {
			return errors.Errorf("expected WithdrawFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case currency.RegisterCurrency:
		fact, ok := t.Fact().(currency.RegisterCurrencyFact)
		if !ok {
			return errors.Errorf("expected RegisterCurrencyFact, not %T", t.Fact())
		}
		did = fact.Currency().Currency().String()
		didtype = DuplicationTypeCurrency
	case currency.UpdateCurrency:
		fact, ok := t.Fact().(currency.UpdateCurrencyFact)
		if !ok {
			return errors.Errorf("expected UpdateCurrencyFact, not %T", t.Fact())
		}
		did = fact.Currency().String()
		didtype = DuplicationTypeCurrency
	case currency.Mint:
	case point.Mint:
		fact, ok := t.Fact().(point.MintFact)
		if !ok {
			return errors.Errorf("expected MintFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case point.RegisterPoint:
		fact, ok := t.Fact().(point.RegisterPointFact)
		if !ok {
			return errors.Errorf("expected RegisterPointFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case point.Burn:
		fact, ok := t.Fact().(point.BurnFact)
		if !ok {
			return errors.Errorf("expected BurnFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case point.Approve:
		fact, ok := t.Fact().(point.ApproveFact)
		if !ok {
			return errors.Errorf("expected ApproveFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case point.Transfer:
		fact, ok := t.Fact().(point.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case point.TransferFrom:
		fact, ok := t.Fact().(point.TransferFromFact)
		if !ok {
			return errors.Errorf("expected TransferFromFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	default:
		return nil
	}

	if err != nil {
		return err
	}

	if did != "" {
		if _, found := opr.Duplicated[did]; found {
			switch didtype {
			case DuplicationTypeSender:
				return errors.Errorf("violates only one sender in proposal")
			case DuplicationTypeCurrency:
				return errors.Errorf("duplicate currency id, %q found in proposal", did)
			default:
				return errors.Errorf("violates duplication in proposal")
			}
		}

		opr.Duplicated[did] = didtype
	}

	if len(newAddresses) > 0 {
		if err := opr.CheckNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}

func checkDuplicateSender(op mitumbase.Operation) (string, currencytypes.DuplicationType, error) {
	fact, ok := op.Fact().(point.PointFact)
	if !ok {
		return "", "", errors.Errorf(utils.ErrStringTypeCast(point.PointFact{}, op.Fact()))
	}
	return fact.Sender().String(), DuplicationTypeSender, nil
}

func GetNewProcessor(opr *currencyprocessor.OperationProcessor, op mitumbase.Operation) (mitumbase.OperationProcessor, bool, error) {
	switch i, err := opr.GetNewProcessorFromHintset(op); {
	case err != nil:
		return nil, false, err
	case i != nil:
		return i, true, nil
	}

	switch t := op.(type) {
	case currency.CreateAccount,
		currency.UpdateKey,
		currency.Transfer,
		extensioncurrency.CreateContractAccount,
		extensioncurrency.Withdraw,
		currency.RegisterCurrency,
		currency.UpdateCurrency,
		currency.Mint,
		point.RegisterPoint,
		point.Mint,
		point.Burn,
		point.Approve,
		point.Transfer,
		point.TransferFrom:
		return nil, false, errors.Errorf("%T needs SetProcessor", t)
	default:
		return nil, false, nil
	}
}
