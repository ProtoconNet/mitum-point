package cmds

import (
	"fmt"
	"github.com/ProtoconNet/mitum-point/types"
)

type PointIDFlag struct {
	CID types.PointID
}

func (v *PointIDFlag) UnmarshalText(b []byte) error {
	cid := types.PointID(string(b))
	if err := cid.IsValid(nil); err != nil {
		return fmt.Errorf("invalid point id, %q, %w", string(b), err)
	}
	v.CID = cid

	return nil
}

func (v *PointIDFlag) String() string {
	return v.CID.String()
}
