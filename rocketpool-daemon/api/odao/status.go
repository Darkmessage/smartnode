package odao

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/dao/oracle"
	"github.com/rocket-pool/rocketpool-go/dao/proposals"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	rptypes "github.com/rocket-pool/rocketpool-go/types"

	"github.com/rocket-pool/smartnode/rocketpool-daemon/common/server"
	"github.com/rocket-pool/smartnode/shared/types/api"
)

// ===============
// === Factory ===
// ===============

type oracleDaoStatusContextFactory struct {
	handler *OracleDaoHandler
}

func (f *oracleDaoStatusContextFactory) Create(args url.Values) (*oracleDaoStatusContext, error) {
	c := &oracleDaoStatusContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *oracleDaoStatusContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterSingleStageRoute[*oracleDaoStatusContext, api.OracleDaoStatusData](
		router, "status", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type oracleDaoStatusContext struct {
	handler     *OracleDaoHandler
	rp          *rocketpool.RocketPool
	nodeAddress common.Address

	odaoMember *oracle.OracleDaoMember
	oSettings  *oracle.OracleDaoSettings
	odaoMgr    *oracle.OracleDaoManager
	dpm        *proposals.DaoProposalManager
}

func (c *oracleDaoStatusContext) Initialize() error {
	sp := c.handler.serviceProvider
	c.rp = sp.GetRocketPool()
	c.nodeAddress, _ = sp.GetWallet().GetAddress()

	// Requirements
	err := sp.RequireNodeRegistered()
	if err != nil {
		return err
	}

	// Bindings
	c.odaoMember, err = oracle.NewOracleDaoMember(c.rp, c.nodeAddress)
	if err != nil {
		return fmt.Errorf("error creating oracle DAO member binding: %w", err)
	}
	c.odaoMgr, err = oracle.NewOracleDaoManager(c.rp)
	if err != nil {
		return fmt.Errorf("error creating Oracle DAO manager binding: %w", err)
	}
	c.oSettings = c.odaoMgr.Settings
	c.dpm, err = proposals.NewDaoProposalManager(c.rp)
	if err != nil {
		return fmt.Errorf("error creating proposal manager binding: %w", err)
	}
	return nil
}

func (c *oracleDaoStatusContext) GetState(mc *batch.MultiCaller) {
	eth.AddQueryablesToMulticall(mc,
		c.odaoMember.Exists,
		c.odaoMember.InvitedTime,
		c.odaoMember.ReplacedTime,
		c.odaoMember.LeftTime,
		c.odaoMgr.MemberCount,
		c.dpm.ProposalCount,
	)
	c.oSettings.Proposal.ActionTime.AddToQuery(mc)
}

func (c *oracleDaoStatusContext) PrepareData(data *api.OracleDaoStatusData, opts *bind.TransactOpts) error {
	// Get the timestamp of the latest block
	latestHeader, err := c.rp.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error getting latest block header: %w", err)
	}
	currentTime := time.Unix(int64(latestHeader.Time), 0)
	actionWindow := c.oSettings.Proposal.ActionTime.Formatted()

	// Check action windows for the current member
	exists := c.odaoMember.Exists.Get()
	data.IsMember = exists
	if exists {
		data.CanLeave = isProposalActionable(actionWindow, c.odaoMember.LeftTime.Formatted(), currentTime)
		data.CanReplace = isProposalActionable(actionWindow, c.odaoMember.ReplacedTime.Formatted(), currentTime)
	} else {
		data.CanJoin = isProposalActionable(actionWindow, c.odaoMember.InvitedTime.Formatted(), currentTime)
	}

	// Total member count
	data.TotalMembers = c.odaoMgr.MemberCount.Formatted()

	// Get the proposals
	_, props, err := c.dpm.GetProposals(c.dpm.ProposalCount.Formatted(), false, nil)
	if err != nil {
		return fmt.Errorf("error getting Oracle DAO proposals: %w", err)
	}

	// Proposal info
	data.ProposalCounts.Total = len(props)
	for _, prop := range props {
		switch prop.State.Formatted() {
		case rptypes.ProposalState_Pending:
			data.ProposalCounts.Pending++
		case rptypes.ProposalState_Active:
			data.ProposalCounts.Active++
		case rptypes.ProposalState_Cancelled:
			data.ProposalCounts.Cancelled++
		case rptypes.ProposalState_Defeated:
			data.ProposalCounts.Defeated++
		case rptypes.ProposalState_Succeeded:
			data.ProposalCounts.Succeeded++
		case rptypes.ProposalState_Expired:
			data.ProposalCounts.Expired++
		case rptypes.ProposalState_Executed:
			data.ProposalCounts.Executed++
		}
	}
	return nil
}