package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *ApproveFact) unpack(
	enc encoder.Encoder,
	sd string,
	bits []byte,
) error {
	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return err
	}
	fact.sender = sender

	hits, err := enc.DecodeSlice(bits)
	if err != nil {
		return err
	}

	items := make([]ApproveItem, len(hits))
	for i, hinter := range hits {
		item, ok := hinter.(ApproveItem)
		if !ok {
			return common.ErrTypeMismatch.Wrap(errors.Errorf("expected ApproveItem, not %T", hinter))
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
