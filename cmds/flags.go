package cmds

import (
	"fmt"
	
	"github.com/ProtoconNet/mitum-point/types"
)

type PointSymbolFlag struct {
	Symbol types.PointSymbol
}

func (v *PointSymbolFlag) UnmarshalText(b []byte) error {
	cid := types.PointSymbol(string(b))
	if err := cid.IsValid(nil); err != nil {
		return fmt.Errorf("invalid point symbol, %q, %w", string(b), err)
	}
	v.Symbol = cid

	return nil
}

func (v *PointSymbolFlag) String() string {
	return v.Symbol.String()
}
