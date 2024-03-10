package security

import (
	"bytes"
	"fmt"

	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/urfave/cli/v2"

	"github.com/rocket-pool/smartnode/rocketpool-cli/utils"
	"github.com/rocket-pool/smartnode/rocketpool-cli/utils/client"
	"github.com/rocket-pool/smartnode/rocketpool-cli/utils/tx"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/input"
)

var voteSupportFlag *cli.StringFlag = &cli.StringFlag{
	Name:    "support",
	Aliases: []string{"s"},
	Usage:   "Whether to support the proposal ('yes' or 'no')",
}

func voteOnProposal(c *cli.Context) error {
	// Get RP client
	rp, err := client.NewClientFromCtx(c).WithReady()
	if err != nil {
		return err
	}

	// Get security council proposals
	proposals, err := rp.Api.Security.Proposals()
	if err != nil {
		return err
	}

	// Get security council members
	allMembers, err := rp.Api.Security.Members()
	if err != nil {
		return err
	}

	// Get votable proposals
	votableProposals := []api.SecurityProposalDetails{}
	for _, proposal := range proposals.Data.Proposals {
		if proposal.State == types.ProposalState_Active && !proposal.MemberVoted {
			votableProposals = append(votableProposals, proposal)
		}
	}

	// Check for votable proposals
	if len(votableProposals) == 0 {
		fmt.Println("No proposals can be voted on.")
		return nil
	}

	// Get selected proposal
	var selectedProposal api.SecurityProposalDetails
	if c.Uint64(proposalFlag.Name) != 0 {
		// Get matching proposal
		selectedId := c.Uint64(proposalFlag.Name)
		found := false
		for _, proposal := range votableProposals {
			if proposal.ID == selectedId {
				selectedProposal = proposal
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("Proposal %d can not be voted on.", selectedId)
		}
	} else {
		// Prompt for proposal selection
		var memberID string
		options := make([]string, len(votableProposals))
		for pi, proposal := range votableProposals {
			for _, member := range allMembers.Data.Members {
				if bytes.Equal(proposal.ProposerAddress.Bytes(), member.Address.Bytes()) {
					memberID = member.ID
				}
			}
			options[pi] = fmt.Sprintf(
				"proposal %d (message: '%s', payload: %s, end time: %s, votes required: %.2f, votes for: %.2f, votes against: %.2f, proposed by: %s (%s))",
				proposal.ID,
				proposal.Message,
				proposal.PayloadStr,
				utils.GetDateTimeStringOfTime(proposal.EndTime),
				proposal.VotesRequired,
				proposal.VotesFor,
				proposal.VotesAgainst,
				memberID,
				proposal.ProposerAddress)
		}
		selected, _ := utils.Select("Please select a proposal to vote on:", options)
		selectedProposal = votableProposals[selected]
	}

	// Get support status
	var support bool
	var supportLabel string
	if c.String(voteSupportFlag.Name) != "" {
		// Parse support status
		var err error
		support, err = input.ValidateBool("support", c.String(voteSupportFlag.Name))
		if err != nil {
			return err
		}
	} else {
		// Prompt for support status
		support = utils.Confirm("Would you like to vote in support of the proposal?")
	}
	if support {
		supportLabel = "in support of"
	} else {
		supportLabel = "against"
	}

	// Build the TX
	response, err := rp.Api.Security.VoteOnProposal(selectedProposal.ID, support)
	if err != nil {
		return err
	}

	// Verify
	if !response.Data.CanVote {
		fmt.Println("Cannot vote on proposal:")
		if response.Data.JoinedAfterCreated {
			fmt.Println("You cannot vote on proposals created before you joined the security council.")
		}
		if response.Data.AlreadyVoted {
			fmt.Println("You already voted on this proposal.")
		}
		if response.Data.DoesNotExist {
			fmt.Println("The proposal does not exist.")
		}
		if response.Data.InvalidState {
			fmt.Println("The proposal is not in a voteable state.")
		}
		return nil
	}

	// Run the TX
	err = tx.HandleTx(c, rp, response.Data.TxInfo,
		fmt.Sprintf("Are you sure you want to vote %s proposal %d? Your vote cannot be changed later.", supportLabel, selectedProposal.ID),
		"voting on proposal",
		"Submitting vote...",
	)
	if err != nil {
		return err
	}

	// Log & return
	fmt.Printf("Successfully voted %s proposal %d.\n", supportLabel, selectedProposal.ID)
	return nil
}