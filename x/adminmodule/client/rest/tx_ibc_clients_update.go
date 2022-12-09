package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

func emptyBadRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Legacy REST Routes are not supported for IBC proposals")
	}
}

func IbcUpgradeProposalEmptyRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "ibc-upgrade",
		Handler:  emptyBadRequestHandlerFn(clientCtx),
	}
}

func ClientUpdateProposalEmptyRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update-client",
		Handler:  emptyBadRequestHandlerFn(clientCtx),
	}
}
