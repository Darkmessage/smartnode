package security

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/dao/proposals"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/dao/security"
	"github.com/rocket-pool/rocketpool-go/rocketpool"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/smartnode/shared/types/api"
)

// ===============
// === Factory ===
// ===============

type securityLeaveContextFactory struct {
	handler *SecurityCouncilHandler
}

func (f *securityLeaveContextFactory) Create(args url.Values) (*securityLeaveContext, error) {
	c := &securityLeaveContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *securityLeaveContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterSingleStageRoute[*securityLeaveContext, api.SecurityLeaveData](
		router, "leave", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type securityLeaveContext struct {
	handler     *SecurityCouncilHandler
	rp          *rocketpool.RocketPool
	nodeAddress common.Address

	scMgr     *security.SecurityCouncilManager
	scMember  *security.SecurityCouncilMember
	dpm       *proposals.DaoProposalManager
	pSettings *protocol.ProtocolDaoSettings
}

func (c *securityLeaveContext) Initialize() (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	c.rp = sp.GetRocketPool()
	c.nodeAddress, _ = sp.GetWallet().GetAddress()

	// Requirements
	err := sp.RequireOnSecurityCouncil()
	if err != nil {
		return types.ResponseStatus_InvalidChainState, err
	}
	// Bindings
	c.scMember, err = security.NewSecurityCouncilMember(c.rp, c.nodeAddress)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error creating security council member binding: %w", err)
	}
	c.dpm, err = proposals.NewDaoProposalManager(c.rp)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error creating DAO proposal manager binding: %w", err)
	}
	pdaoMgr, err := protocol.NewProtocolDaoManager(c.rp)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error creating protocol DAO manager binding: %w", err)
	}
	c.pSettings = pdaoMgr.Settings
	c.scMgr, err = security.NewSecurityCouncilManager(c.rp, c.pSettings)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error creating security council manager binding: %w", err)
	}
	return types.ResponseStatus_Success, nil
}

func (c *securityLeaveContext) GetState(mc *batch.MultiCaller) {
	eth.AddQueryablesToMulticall(mc,
		c.scMember.Exists,
		c.scMember.LeftTime,
		c.pSettings.Security.ProposalActionTime,
	)
}

func (c *securityLeaveContext) PrepareData(data *api.SecurityLeaveData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	leftTime := c.scMember.LeftTime.Formatted()
	actionTime := c.pSettings.Security.ProposalActionTime.Formatted()
	data.ProposalExpired = time.Until(leftTime.Add(actionTime)) < 0
	data.IsNotMember = !c.scMember.Exists.Get()
	data.CanLeave = !(data.ProposalExpired || data.IsNotMember)

	// Get the tx
	if data.CanLeave && opts != nil {
		txInfo, err := c.scMgr.Leave(opts)
		if err != nil {
			return types.ResponseStatus_Error, fmt.Errorf("error getting TX info for Leave: %w", err)
		}
		data.TxInfo = txInfo
	}
	return types.ResponseStatus_Success, nil
}
