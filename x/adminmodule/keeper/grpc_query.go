package keeper

import (
	"github.com/cosmos/admin-module/v2/x/adminmodule/types"
)

var _ types.QueryServer = Keeper{}

type Querier struct {
	Keeper
}
