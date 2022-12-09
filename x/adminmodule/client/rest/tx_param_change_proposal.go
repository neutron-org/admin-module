package rest

import (
	"net/http"

	admintypes "github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	paramscutils "github.com/cosmos/cosmos-sdk/x/params/client/utils"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

// ParamChangeProposalReq defines a parameter change proposal request body.
type ParamChangeProposalReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string                        `json:"title" yaml:"title"`
	Description string                        `json:"description" yaml:"description"`
	Changes     paramscutils.ParamChangesJSON `json:"changes" yaml:"changes"`
	Proposer    sdk.AccAddress                `json:"proposer" yaml:"proposer"`
}

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the param
// change REST handler with a given sub-route.
func ParamChangeProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "param_change",
		Handler:  postParamChangeProposalHandlerFn(clientCtx),
	}
}

func postParamChangeProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ParamChangeProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := proposal.NewParameterChangeProposal(req.Title, req.Description, req.Changes.ToParamChanges())

		msg, err := admintypes.NewMsgSubmitProposal(content, req.Proposer)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
