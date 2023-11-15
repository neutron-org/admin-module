package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "adminmodule"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_adminmodule"

	AdminKey = "Admin-"

	ArchiveKey = "ArchiveSdk47-"

	ArchiveLegacyKey = "Archive-"
)

var (
	ProposalsKeyLegacyPrefix = []byte{0x00}
	ProposalIDKeyLegacy      = []byte{0x03}
	ProposalsKeyPrefix       = []byte{0x04}
	ProposalIDKey            = []byte{0x06}
)

// GetProposalIDBytes returns the byte representation of the proposalID
func GetProposalIDBytes(proposalID uint64) (proposalIDBz []byte) {
	proposalIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(proposalIDBz, proposalID)
	return
}

// GetProposalIDFromBytes returns proposalID in uint64 format from a byte array
func GetProposalIDFromBytes(bz []byte) (proposalID uint64) {
	return binary.BigEndian.Uint64(bz)
}

// ProposalKey gets a specific proposal from the store
func ProposalKey(proposalID uint64) []byte {
	return append(ProposalsKeyPrefix, GetProposalIDBytes(proposalID)...)
}

// ProposalLegacyKey gets a specific proposal from the store
func ProposalLegacyKey(proposalID uint64) []byte {
	return append(ProposalsKeyLegacyPrefix, GetProposalIDBytes(proposalID)...)
}
