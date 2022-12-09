package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/tx"

	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"

	admintypes "github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// PlanRequest defines a proposal for a new upgrade plan.
type PlanRequest struct {
	BaseReq       rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title         string       `json:"title" yaml:"title"`
	Description   string       `json:"description" yaml:"description"`
	UpgradeName   string       `json:"upgrade_name" yaml:"upgrade_name"`
	UpgradeHeight int64        `json:"upgrade_height" yaml:"upgrade_height"`
	UpgradeInfo   string       `json:"upgrade_info" yaml:"upgrade_info"`
}

// CancelRequest defines a proposal to cancel a current plan.
type CancelRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
}

func SoftwareUpgradeProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "upgrade",
		Handler:  newPostPlanHandler(clientCtx),
	}
}

func CancelUpgradeProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "cancel_upgrade",
		Handler:  newCancelPlanHandler(clientCtx),
	}
}

func newPostPlanHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PlanRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		plan := types.Plan{Name: req.UpgradeName, Height: req.UpgradeHeight, Info: req.UpgradeInfo}
		content := types.NewSoftwareUpgradeProposal(req.Title, req.Description, plan)
		msg, err := admintypes.NewMsgSubmitProposal(content, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func newCancelPlanHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CancelRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewCancelSoftwareUpgradeProposal(req.Title, req.Description)

		msg, err := admintypes.NewMsgSubmitProposal(content, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
