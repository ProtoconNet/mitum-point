package digest

import (
	"net/http"
	"time"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-point/types"
)

func (hd *Handlers) handlePoint(w http.ResponseWriter, r *http.Request) {
	cachekey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handlePointInGroup(contract)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cachekey, time.Millisecond*100)
		}
	}
}

func (hd *Handlers) handlePointInGroup(contract string) (interface{}, error) {
	switch design, err := Point(hd.database, contract); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildPointHal(contract, *design)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildPointHal(contract string, design types.Design) (cdigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathPoint, "contract", contract)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(design, cdigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handlePointBalance(w http.ResponseWriter, r *http.Request) {
	cachekey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	account, err, status := cdigest.ParseRequest(w, r, "address")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handlePointBalanceInGroup(contract, account)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cachekey, time.Millisecond*100)
		}
	}
}

func (hd *Handlers) handlePointBalanceInGroup(contract, account string) (interface{}, error) {
	switch amount, err := PointBalance(hd.database, contract, account); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildPointBalanceHal(contract, account, amount)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildPointBalanceHal(contract, account string, amount *common.Big) (cdigest.Hal, error) {
	var hal cdigest.Hal

	if amount == nil {
		hal = cdigest.NewEmptyHal()
	} else {
		h, err := hd.combineURL(HandlerPathPointBalance, "contract", contract, "address", account)
		if err != nil {
			return nil, err
		}

		hal = cdigest.NewBaseHal(struct {
			Amount common.Big `json:"amount"`
		}{Amount: *amount}, cdigest.NewHalLink(h, nil))
	}

	return hal, nil
}
