package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/pkg/errors"
	"regexp"
)

var (
	MinLengthPointSymbol = 3
	MaxLengthPointSymbol = 10
	ReValidPointSymbol   = regexp.MustCompile(`^[A-Z0-9][A-Z0-9_\.\!\$\*\@]*[A-Z0-9]$`)
	ReSpcecialChar       = regexp.MustCompile(`^[^\s:/?#\[\]@]*$`)
)

type PointSymbol string

func (ps PointSymbol) Bytes() []byte {
	return []byte(ps)
}

func (ps PointSymbol) String() string {
	return string(ps)
}

func (ps PointSymbol) IsValid([]byte) error {
	if l := len(ps); l < MinLengthPointSymbol || l > MaxLengthPointSymbol {
		return common.ErrValOOR.Wrap(errors.Errorf(
			"invalid length of point symbol, %d <= %d <= %d", MinLengthPointSymbol, l, MaxLengthPointSymbol))
	} else if !ReValidPointSymbol.Match([]byte(ps)) {
		return common.ErrValueInvalid.Wrap(errors.Errorf("wrong point symbol, %v", ps))
	}

	return nil
}
