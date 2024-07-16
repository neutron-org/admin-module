package adminmodule

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/client/v2/autocli"
)

var _ autocli.HasAutoCLIConfig = AppModule{}

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: "cosmos.adminmodule.adminmodule.Query",
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Admins",
					Use:       "admins",
					Short:     "Query list of admins",
				},
				{
					RpcMethod: "ArchivedProposals",
					Use:       "archivedproposals",
					Short:     "Query archived proposals",
				},
				{
					RpcMethod: "ArchivedProposalsLegacy",
					Use:       "archivedproposalslegacy",
					Short:     "Query legacy archived proposals",
				},
			},
			EnhanceCustomCommand: true,
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: "cosmos.adminmodule.adminmodule.Msg",
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "DeleteAdmin",
					Use:            "delete-admin [admin]",
					Short:          "Delete admin",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "admin"}},
				},
				{
					RpcMethod:      "AddAdmin",
					Use:            "add-admin [admin]",
					Short:          "Add new admin",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "admin"}},
				},
				{
					RpcMethod: "SubmitProposalLegacy",
					Skip:      true,
				},
			},
			EnhanceCustomCommand: true,
		},
	}
}
