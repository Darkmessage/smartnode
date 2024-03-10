package network

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/deposit"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/tokens"

	"github.com/rocket-pool/smartnode/rocketpool-daemon/common/server"
	"github.com/rocket-pool/smartnode/shared/types/api"
)

// ===============
// === Factory ===
// ===============

type networkStatsContextFactory struct {
	handler *NetworkHandler
}

func (f *networkStatsContextFactory) Create(args url.Values) (*networkStatsContext, error) {
	c := &networkStatsContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *networkStatsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterSingleStageRoute[*networkStatsContext, api.NetworkStatsData](
		router, "stats", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type networkStatsContext struct {
	handler *NetworkHandler
	rp      *rocketpool.RocketPool

	depositPool   *deposit.DepositPoolManager
	nodeMgr       *node.NodeManager
	mpMgr         *minipool.MinipoolManager
	networkMgr    *network.NetworkManager
	reth          *tokens.TokenReth
	smoothingPool *core.Contract
}

func (c *networkStatsContext) Initialize() error {
	sp := c.handler.serviceProvider
	c.rp = sp.GetRocketPool()

	// Requirements
	err := sp.RequireEthClientSynced()
	if err != nil {
		return err
	}

	// Bindings
	c.depositPool, err = deposit.NewDepositPoolManager(c.rp)
	if err != nil {
		return fmt.Errorf("error getting deposit pool manager binding: %w", err)
	}
	c.nodeMgr, err = node.NewNodeManager(c.rp)
	if err != nil {
		return fmt.Errorf("error getting node manager binding: %w", err)
	}
	c.mpMgr, err = minipool.NewMinipoolManager(c.rp)
	if err != nil {
		return fmt.Errorf("error getting minipool manager binding: %w", err)
	}
	c.reth, err = tokens.NewTokenReth(c.rp)
	if err != nil {
		return fmt.Errorf("error getting rETH token binding: %w", err)
	}
	c.smoothingPool, err = c.rp.GetContract(rocketpool.ContractName_RocketSmoothingPool)
	if err != nil {
		return fmt.Errorf("error getting rETH token binding: %w", err)
	}
	c.networkMgr, err = network.NewNetworkManager(c.rp)
	if err != nil {
		return fmt.Errorf("error creating network prices binding: %w", err)
	}
	return nil
}

func (c *networkStatsContext) GetState(mc *batch.MultiCaller) {
	eth.AddQueryablesToMulticall(mc,
		c.depositPool.Balance,
		c.mpMgr.TotalQueueCapacity,
		c.networkMgr.EthUtilizationRate,
		c.networkMgr.NodeFee,
		c.nodeMgr.NodeCount,
		c.mpMgr.MinipoolCount,
		c.mpMgr.FinalisedMinipoolCount,
		c.networkMgr.RplPrice,
		c.nodeMgr.TotalRplStake,
		c.reth.ExchangeRate,
	)
}

func (c *networkStatsContext) PrepareData(data *api.NetworkStatsData, opts *bind.TransactOpts) error {
	// Handle the details
	data.DepositPoolBalance = c.depositPool.Balance.Get()
	data.MinipoolCapacity = c.mpMgr.TotalQueueCapacity.Get()
	data.StakerUtilization = c.networkMgr.EthUtilizationRate.Raw()
	data.NodeFee = c.networkMgr.NodeFee.Raw()
	data.NodeCount = c.nodeMgr.NodeCount.Formatted()
	data.RplPrice = c.networkMgr.RplPrice.Raw()
	data.TotalRplStaked = c.nodeMgr.TotalRplStake.Get()
	data.RethPrice = c.reth.ExchangeRate.Raw()

	// Get the total effective RPL stake
	effectiveRplStaked, err := c.nodeMgr.GetTotalEffectiveRplStake(c.nodeMgr.NodeCount.Formatted(), nil)
	if err != nil {
		return fmt.Errorf("error getting total effective RPL stake: %w", err)
	}
	data.EffectiveRplStaked = effectiveRplStaked

	// Get the minipool counts by status
	data.FinalizedMinipoolCount = c.mpMgr.FinalisedMinipoolCount.Formatted()
	minipoolCounts, err := c.mpMgr.GetMinipoolCountPerStatus(c.mpMgr.MinipoolCount.Formatted(), nil)
	if err != nil {
		return fmt.Errorf("error getting minipool counts per status: %w", err)
	}
	data.InitializedMinipoolCount = minipoolCounts.Initialized.Uint64()
	data.PrelaunchMinipoolCount = minipoolCounts.Prelaunch.Uint64()
	data.StakingMinipoolCount = minipoolCounts.Staking.Uint64()
	data.WithdrawableMinipoolCount = minipoolCounts.Withdrawable.Uint64()
	data.DissolvedMinipoolCount = minipoolCounts.Dissolved.Uint64()

	// Get the number of nodes opted into the smoothing pool
	spCount, err := c.nodeMgr.GetSmoothingPoolRegisteredNodeCount(c.nodeMgr.NodeCount.Formatted(), nil)
	if err != nil {
		return fmt.Errorf("error getting smoothing pool opt-in count: %w", err)
	}
	data.SmoothingPoolNodes = spCount

	// Get the smoothing pool balance
	data.SmoothingPoolAddress = *c.smoothingPool.Address
	smoothingPoolBalance, err := c.rp.Client.BalanceAt(context.Background(), *c.smoothingPool.Address, nil)
	if err != nil {
		return fmt.Errorf("error getting smoothing pool balance: %w", err)
	}
	data.SmoothingPoolBalance = smoothingPoolBalance

	// Get the TVL
	activeMinipools := data.InitializedMinipoolCount +
		data.PrelaunchMinipoolCount +
		data.StakingMinipoolCount +
		data.WithdrawableMinipoolCount +
		data.DissolvedMinipoolCount
	tvl := big.NewInt(int64(activeMinipools))
	tvl.Mul(tvl, big.NewInt(32))
	tvl.Add(tvl, data.DepositPoolBalance)
	tvl.Add(tvl, data.MinipoolCapacity)
	rplWorth := big.NewInt(0).Mul(data.TotalRplStaked, data.RplPrice)
	tvl.Add(tvl, rplWorth)
	data.TotalValueLocked = tvl

	return nil
}