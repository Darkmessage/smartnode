package pdao

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"

	"github.com/rocket-pool/smartnode/rocketpool/common/beacon"
	"github.com/rocket-pool/smartnode/rocketpool/common/server"
	"github.com/rocket-pool/smartnode/shared/config"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/input"
)

// ===============
// === Factory ===
// ===============

type protocolDaoOverrideVoteOnProposalContextFactory struct {
	handler *ProtocolDaoHandler
}

func (f *protocolDaoOverrideVoteOnProposalContextFactory) Create(vars map[string]string) (*protocolDaoOverrideVoteOnProposalContext, error) {
	c := &protocolDaoOverrideVoteOnProposalContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("id", vars, input.ValidatePositiveUint, &c.proposalID),
		server.ValidateArg("vote", vars, input.ValidateVoteDirection, &c.voteDirection),
	}
	return c, errors.Join(inputErrs...)
}

func (f *protocolDaoOverrideVoteOnProposalContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterSingleStageRoute[*protocolDaoOverrideVoteOnProposalContext, api.ProtocolDaoOverrideVoteOnProposalData](
		router, "proposal/override-vote", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type protocolDaoOverrideVoteOnProposalContext struct {
	handler     *ProtocolDaoHandler
	cfg         *config.RocketPoolConfig
	rp          *rocketpool.RocketPool
	bc          beacon.Client
	nodeAddress common.Address

	proposalID      uint64
	voteDirection   types.VoteDirection
	node            *node.Node
	existingVoteDir types.VoteDirection
	pdaoMgr         *protocol.ProtocolDaoManager
	proposal        *protocol.ProtocolDaoProposal
}

func (c *protocolDaoOverrideVoteOnProposalContext) Initialize() error {
	sp := c.handler.serviceProvider
	c.cfg = sp.GetConfig()
	c.rp = sp.GetRocketPool()
	c.bc = sp.GetBeaconClient()
	c.nodeAddress, _ = sp.GetWallet().GetAddress()

	// Requirements
	err := sp.RequireNodeRegistered()
	if err != nil {
		return err
	}

	// Bindings
	c.node, err = node.NewNode(c.rp, c.nodeAddress)
	if err != nil {
		return fmt.Errorf("error creating node %s binding: %w", c.nodeAddress.Hex(), err)
	}
	c.pdaoMgr, err = protocol.NewProtocolDaoManager(c.rp)
	if err != nil {
		return fmt.Errorf("error creating protocol DAO manager binding: %w", err)
	}
	c.proposal, err = protocol.NewProtocolDaoProposal(c.rp, c.proposalID)
	if err != nil {
		return fmt.Errorf("error creating proposal binding: %w", err)
	}
	return nil
}

func (c *protocolDaoOverrideVoteOnProposalContext) GetState(mc *batch.MultiCaller) {
	core.AddQueryablesToMulticall(mc,
		c.pdaoMgr.ProposalCount,
		c.proposal.State,
		c.proposal.TargetBlock,
	)
	c.proposal.GetAddressVoteDirection(mc, &c.existingVoteDir, c.nodeAddress)
}

func (c *protocolDaoOverrideVoteOnProposalContext) PrepareData(data *api.ProtocolDaoOverrideVoteOnProposalData, opts *bind.TransactOpts) error {
	// Get the voting power for the node as of this proposal
	err := c.rp.Query(func(mc *batch.MultiCaller) error {
		c.node.GetVotingPowerAtBlock(mc, &data.VotingPower, c.proposal.TargetBlock.Formatted())
		return nil
	}, nil)
	if err != nil {
		return fmt.Errorf("error getting node voting power at block %d: %w", c.proposal.TargetBlock.Formatted(), err)
	}

	data.DoesNotExist = (c.proposalID > c.pdaoMgr.ProposalCount.Formatted())
	data.InvalidState = (c.proposal.State.Formatted() != types.ProtocolDaoProposalState_ActivePhase2)
	data.AlreadyVoted = (c.existingVoteDir != types.VoteDirection_NoVote)
	data.InsufficientPower = (data.VotingPower.Cmp(common.Big0) == 0)
	data.CanVote = !(data.DoesNotExist || data.InvalidState || data.AlreadyVoted || data.InsufficientPower)

	// Get the tx
	if data.CanVote && opts != nil {
		txInfo, err := c.proposal.OverrideVote(c.voteDirection, opts)
		if err != nil {
			return fmt.Errorf("error getting TX info for Vote: %w", err)
		}
		data.TxInfo = txInfo
	}
	return nil
}
