package client

import (
	"github.com/cosmos/admin-module/x/adminmodule/client/cli"
	"github.com/cosmos/admin-module/x/adminmodule/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// Param change proposal handler.
var ParamChangeProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd, rest.ParamChangeProposalRESTHandler)

// Software upgrade proposal handler.
var SoftwareUpgradeProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal, rest.SoftwareUpgradeProposalRESTHandler)

// Cancel software upgrade proposal handler.
var CancelUpgradeProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCancelUpgradeProposal, rest.CancelUpgradeProposalRESTHandler)

// Community pool spend proposal handler.
var CommunityPoolSpendProposalHandler = govclient.NewProposalHandler(cli.NewSubmitPoolSpendProposalTxCmd, rest.CommunityPoolSpendProposalRESTHandler)

// IBC Client upgrade proposal handler.
var IBCClientUpgradeProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitIbcClientUpgradeProposal, rest.IbcUpgradeProposalEmptyRESTHandler)

// IBC Client update proposal handler.
var IBCClientUpdateProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpdateClientProposal, rest.ClientUpdateProposalEmptyRESTHandler)
