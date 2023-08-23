package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc                       codec.Codec
		storeKey                  storetypes.StoreKey
		memKey                    storetypes.StoreKey
		rtr                       govv1beta1types.Router
		msgServiceRouter          *baseapp.MsgServiceRouter
		IsProposalTypeWhitelisted func(govv1beta1types.Content) bool
		IsModuleWhiteliested      func(typeUrl string) bool
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	rtr govv1beta1types.Router,
	msgServiceRouter *baseapp.MsgServiceRouter,
	isProposalTypeWhitelisted func(govv1beta1types.Content) bool,
	IsModuleWhiteliested func(typeUrl string) bool,
) *Keeper {
	return &Keeper{
		cdc:                       cdc,
		storeKey:                  storeKey,
		memKey:                    memKey,
		rtr:                       rtr,
		msgServiceRouter:          msgServiceRouter,
		IsProposalTypeWhitelisted: isProposalTypeWhitelisted,
		IsModuleWhiteliested:      IsModuleWhiteliested,
	}
}

// RouterLegacy returns the adminmodule Keeper's govtypeRouter
func (k Keeper) RouterLegacy() govv1beta1types.Router {
	return k.rtr
}

// Router returns the adminmodule Keeper's Router
func (k Keeper) Router() *baseapp.MsgServiceRouter {
	return k.msgServiceRouter
}

// Router returns the adminmodule Keeper's Logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Router returns the adminmodule Keeper's Codec
func (k Keeper) Codec() codec.Codec {
	return k.cdc
}
