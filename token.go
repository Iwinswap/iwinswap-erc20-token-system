package token

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// ErrTokenNotFound is returned when a token lookup fails.
	ErrTokenNotFound = errors.New("token not found")
	// ErrAlreadyExists is returned when trying to add a token that is already in the registry.
	ErrAlreadyExists = errors.New("token already exists")
)

// TokenView is a safe, structured representation of a token's data for external use.
type TokenView struct {
	ID                   uint64         `json:"id"`
	Address              common.Address `json:"address"`
	Name                 string         `json:"name"`
	Symbol               string         `json:"symbol"`
	Decimals             uint8          `json:"decimals"`
	FeeOnTransferPercent float64        `json:"feeOnTransferPercent"`
	GasForTransfer       uint64         `json:"gasForTransfer"`
}

// TokenRegistry manages a collection of token data using a Struct-of-Arrays layout.
type TokenRegistry struct {
	// --- Physical data storage (Struct of Arrays) ---
	address              []common.Address
	name                 []string
	symbol               []string
	decimals             []uint8
	feeOnTransferPercent []float64
	gasForTransfer       []uint64
	id                   []uint64 // Stores the stable ID for each index

	// --- Mapping layers to separate logical ID from physical index ---
	nextID      uint64                    // A counter to generate new, permanent IDs
	idToIndex   map[uint64]int            // Maps a permanent ID to its current slice index
	addressToID map[common.Address]uint64 // Maps an address to its permanent ID
}

// NewTokenRegistry creates and initializes a new, empty TokenRegistry.
func NewTokenRegistry() *TokenRegistry {
	return &TokenRegistry{
		// Initialize with a capacity to reduce initial reallocations
		address:              make([]common.Address, 0, 128),
		name:                 make([]string, 0, 128),
		symbol:               make([]string, 0, 128),
		decimals:             make([]uint8, 0, 128),
		feeOnTransferPercent: make([]float64, 0, 128),
		gasForTransfer:       make([]uint64, 0, 128),
		id:                   make([]uint64, 0, 128),

		nextID:      1, // Start IDs at 1 to avoid confusion with zero-values
		idToIndex:   make(map[uint64]int),
		addressToID: make(map[common.Address]uint64),
	}
}

// AddToken adds a new token to the registry and assigns it a new, permanent ID.
func addToken(addr common.Address, name, symbol string, decimals uint8, registry *TokenRegistry) (uint64, error) {
	if _, exists := registry.addressToID[addr]; exists {
		return 0, ErrAlreadyExists
	}

	newID := registry.nextID
	registry.nextID++

	newIndex := len(registry.address)
	registry.address = append(registry.address, addr)
	registry.name = append(registry.name, name)
	registry.symbol = append(registry.symbol, symbol)
	registry.decimals = append(registry.decimals, decimals)
	registry.feeOnTransferPercent = append(registry.feeOnTransferPercent, 0)
	registry.gasForTransfer = append(registry.gasForTransfer, 0)
	registry.id = append(registry.id, newID)

	registry.idToIndex[newID] = newIndex
	registry.addressToID[addr] = newID

	return newID, nil
}

// deleteToken removes a token using the "swap-and-pop" algorithm.
func deleteToken(idToDelete uint64, registry *TokenRegistry) error {
	indexToDelete, ok := registry.idToIndex[idToDelete]
	if !ok {
		return ErrTokenNotFound
	}

	addressToDelete := registry.address[indexToDelete]
	lastIndex := len(registry.address) - 1

	if indexToDelete != lastIndex {
		lastID := registry.id[lastIndex]
		registry.address[indexToDelete] = registry.address[lastIndex]
		registry.name[indexToDelete] = registry.name[lastIndex]
		registry.symbol[indexToDelete] = registry.symbol[lastIndex]
		registry.decimals[indexToDelete] = registry.decimals[lastIndex]
		registry.feeOnTransferPercent[indexToDelete] = registry.feeOnTransferPercent[lastIndex]
		registry.gasForTransfer[indexToDelete] = registry.gasForTransfer[lastIndex]
		registry.id[indexToDelete] = lastID
		registry.idToIndex[lastID] = indexToDelete
	}

	registry.address = registry.address[:lastIndex]
	registry.name = registry.name[:lastIndex]
	registry.symbol = registry.symbol[:lastIndex]
	registry.decimals = registry.decimals[:lastIndex]
	registry.feeOnTransferPercent = registry.feeOnTransferPercent[:lastIndex]
	registry.gasForTransfer = registry.gasForTransfer[:lastIndex]
	registry.id = registry.id[:lastIndex]

	delete(registry.idToIndex, idToDelete)
	delete(registry.addressToID, addressToDelete)

	return nil
}

// updateToken updates the mutable data for a token.
func updateToken(id uint64, feeOnTransferPercent float64, gasForTransfer uint64, registry *TokenRegistry) error {
	index, ok := registry.idToIndex[id]
	if !ok {
		return ErrTokenNotFound
	}
	registry.feeOnTransferPercent[index] = feeOnTransferPercent
	registry.gasForTransfer[index] = gasForTransfer
	return nil
}

// getTokenByID returns a safe, structured view of a single token by its permanent ID.
func getTokenByID(id uint64, registry *TokenRegistry) (TokenView, error) {
	index, ok := registry.idToIndex[id]
	if !ok {
		return TokenView{}, ErrTokenNotFound
	}

	return TokenView{
		ID:                   id,
		Address:              registry.address[index],
		Name:                 registry.name[index],
		Symbol:               registry.symbol[index],
		Decimals:             registry.decimals[index],
		FeeOnTransferPercent: registry.feeOnTransferPercent[index],
		GasForTransfer:       registry.gasForTransfer[index],
	}, nil
}

// getTokenByAddress finds a token by its address and returns its view.
func getTokenByAddress(addr common.Address, registry *TokenRegistry) (TokenView, error) {
	id, ok := registry.addressToID[addr]
	if !ok {
		return TokenView{}, ErrTokenNotFound
	}
	index := registry.idToIndex[id]
	return TokenView{
		ID:                   id,
		Address:              registry.address[index],
		Name:                 registry.name[index],
		Symbol:               registry.symbol[index],
		Decimals:             registry.decimals[index],
		FeeOnTransferPercent: registry.feeOnTransferPercent[index],
		GasForTransfer:       registry.gasForTransfer[index],
	}, nil
}

// viewRegistry returns a slice of views for all active tokens in the registry.
func viewRegistry(registry *TokenRegistry) []TokenView {
	length := len(registry.address)
	views := make([]TokenView, length)
	for i := 0; i < length; i++ {
		views[i] = TokenView{
			ID:                   registry.id[i],
			Address:              registry.address[i],
			Name:                 registry.name[i],
			Symbol:               registry.symbol[i],
			Decimals:             registry.decimals[i],
			FeeOnTransferPercent: registry.feeOnTransferPercent[i],
			GasForTransfer:       registry.gasForTransfer[i],
		}
	}
	return views
}
