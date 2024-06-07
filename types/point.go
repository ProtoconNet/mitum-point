package types

import (
	"github.com/pkg/errors"
	"regexp"
)

var (
	MinLengthPointID = 3
	MaxLengthPointID = 10
	ReValidPointID   = regexp.MustCompile(`^[A-Z0-9][A-Z0-9_\.\!\$\*\@]*[A-Z0-9]$`)
	ReSpcecialChar   = regexp.MustCompile(`^[^\s:/?#\[\]@]*$`)
)

type PointID string

func (cid PointID) Bytes() []byte {
	return []byte(cid)
}

func (cid PointID) String() string {
	return string(cid)
}

func (cid PointID) IsValid([]byte) error {
	if l := len(cid); l < MinLengthPointID || l > MaxLengthPointID {
		return errors.Errorf(
			"invalid length of point id, %d <= %d <= %d", MinLengthPointID, l, MaxLengthPointID)
	} else if !ReValidPointID.Match([]byte(cid)) {
		return errors.Errorf("wrong point id, %v", cid)
	}

	return nil
}
