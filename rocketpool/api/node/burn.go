package node

import (
	"fmt"
	"math/big"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/tokens"

	"github.com/rocket-pool/smartnode/shared/types/api"
)

type nodeBurnHandler struct {
	amountWei *big.Int
	reth      *tokens.TokenReth
	balance   *big.Int
}

func (h *nodeBurnHandler) CreateBindings(ctx *callContext) error {
	rp := ctx.rp
	var err error

	h.reth, err = tokens.NewTokenReth(rp)
	if err != nil {
		return fmt.Errorf("error creating reth binding: %w", err)
	}
	return nil
}

func (h *nodeBurnHandler) GetState(ctx *callContext, mc *batch.MultiCaller) {
	node := ctx.node
	h.reth.GetBalance(mc, &h.balance, node.Details.Address)
	h.reth.GetTotalCollateral(mc)
}

func (h *nodeBurnHandler) PrepareResponse(ctx *callContext, response *api.NodeBurnResponse) error {
	opts := ctx.opts

	// Check for validity
	response.InsufficientBalance = (h.amountWei.Cmp(h.balance) > 0)
	response.InsufficientCollateral = (h.amountWei.Cmp(h.reth.Details.TotalCollateral) > 0)
	response.CanBurn = !(response.InsufficientBalance || response.InsufficientCollateral)

	// Get tx info
	if response.CanBurn {
		txInfo, err := h.reth.Burn(h.amountWei, opts)
		if err != nil {
			return fmt.Errorf("error getting TX info for Burn: %w", err)
		}
		response.TxInfo = txInfo
	}
	return nil
}
