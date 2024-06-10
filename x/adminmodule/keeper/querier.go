package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new adminmodule Querier instance
func NewQuerier(keeper *Keeper) Querier {
	return Querier{Keeper: *keeper}
}

// nolint: unparam
func queryAdmins(ctx sdk.Context, _ []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	admins := keeper.GetAdmins(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, admins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrortypes.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// nolint: unparam
func queryArchivedProposals(ctx sdk.Context, _ []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	proposals := keeper.GetArchivedProposals(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, proposals)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrortypes.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
