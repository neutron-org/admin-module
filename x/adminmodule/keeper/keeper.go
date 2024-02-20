package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc                       codec.Codec
		storeKey                  storetypes.StoreKey
		memKey                    storetypes.StoreKey
		legacyRouter              govv1beta1types.Router
		msgServiceRouter          *baseapp.MsgServiceRouter
		isProposalTypeWhitelisted func(govv1beta1types.Content) bool
		isMessageWhitelisted      func(message sdk.Msg) bool
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	legacyRouter govv1beta1types.Router,
	msgServiceRouter *baseapp.MsgServiceRouter,
	isProposalTypeWhitelisted func(govv1beta1types.Content) bool,
	isMessageWhitelisted func(msg sdk.Msg) bool,
) *Keeper {
	return &Keeper{
		cdc:                       cdc,
		storeKey:                  storeKey,
		memKey:                    memKey,
		legacyRouter:              legacyRouter,
		msgServiceRouter:          msgServiceRouter,
		isProposalTypeWhitelisted: isProposalTypeWhitelisted,
		isMessageWhitelisted:      isMessageWhitelisted,
	}
}

// RouterLegacy returns the adminmodule Keeper's govtypeRouter
func (k Keeper) RouterLegacy() govv1beta1types.Router {
	return k.legacyRouter
}

// Router returns the adminmodule Keeper's Router
func (k Keeper) Router() *baseapp.MsgServiceRouter {
	return k.msgServiceRouter
}

// Logger returns the adminmodule Keeper's Logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codec returns the adminmodule Keeper's Codec
func (k Keeper) Codec() codec.Codec {
	return k.cdc
}

// IsProposalTypeWhitelisted returns the adminmodule Keeper's isProposalTypeWhitelisted
func (k Keeper) IsProposalTypeWhitelisted() func(govv1beta1types.Content) bool {
	return k.isProposalTypeWhitelisted
}

// IsMessageWhitelisted returns the adminmodule Keeper's isMessageWhitelisted
func (k Keeper) IsMessageWhitelisted() func(msg sdk.Msg) bool {
	return k.isMessageWhitelisted
}
