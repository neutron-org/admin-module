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
		cdc                           codec.Codec
		storeKey                      storetypes.StoreKey
		memKey                        storetypes.StoreKey
		rtr                           govv1beta1types.Router
		msgServiceRouter              *baseapp.MsgServiceRouter
		IsProposalTypeWhitelisted     func(govv1beta1types.Content) bool
		RegisteredModulesUpdateParams map[string]RegisteredModuleUpdateParams
		// this line is used by starport scaffolding # ibc/keeper/attribute
	}
)

type RegisteredModuleUpdateParams struct {
	// Unique parameter struct of given module
	ParamsMsg interface{}
	// Unique parameters update struct of given module, implements sdk.Msg
	// satisfying adminmodule.SumbitProposal(msgs []sdk.Msg ,<..>)

	UpdateParamsMsg sdk.Msg
}

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	rtr govv1beta1types.Router,
	msgServiceRouter *baseapp.MsgServiceRouter,
	isProposalTypeWhitelisted func(govv1beta1types.Content) bool,
	RegisteredModulesUpdate map[string]RegisteredModuleUpdateParams,
// this line is used by starport scaffolding # ibc/keeper/parameter
) *Keeper {
	return &Keeper{
		cdc:                           cdc,
		storeKey:                      storeKey,
		memKey:                        memKey,
		rtr:                           rtr,
		msgServiceRouter:              msgServiceRouter,
		IsProposalTypeWhitelisted:     isProposalTypeWhitelisted,
		RegisteredModulesUpdateParams: RegisteredModulesUpdate,
		// this line is used by starport scaffolding # ibc/keeper/return
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

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
