# Iwinswap ERC20 Token System

A high-performance, concurrency-safe, in-memory registry for ERC20 token metadata in Go.

---

## Overview

The Iwinswap ERC20 Token System provides a robust and highly efficient way to manage collections of ERC20-like token data. It is designed from the ground up for performance-critical applications that require safe, concurrent access to token metadata without sacrificing speed.

The system is built on a Data-Oriented Design (DOD) core for maximum efficiency, wrapped in a thread-safe service layer for ease of use and safety.

---

## Key Features

* **High-Performance Data Layout**
  Uses a Struct of Arrays (SoA) layout to ensure cache-friendly iteration and fast access in data-heavy workloads.

* **Stable, Permanent Identifiers**
  Every token is assigned a permanent `uint64` ID that never changes or gets reused, enabling safe long-term references.

* **O(1) Operations**
  All core actions—adding, deleting, and looking up tokens by ID or address—have constant time complexity.

* **Concurrency-Safe API**
  A `sync.RWMutex` allows for high reader concurrency while ensuring atomic, safe writes.

* **Efficient Deletion**
  Implements a "swap-and-pop" strategy to keep the data dense and deletions fast.

* **Rigorously Tested**
  Includes unit tests, fuzz tests, race detection, and benchmarks.

---

## Installation

```bash
go get github.com/Iwinswap/iwinswap-token-analyzer
```

---

## Usage

Here’s a basic example using the `TokenSystem`:

```go
package main

import (
	"fmt"
	"log"

	"github.com/Iwinswap/iwinswap-erc20-token-system"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	tokenSystem := token.NewTokenSystem()

	addressA := common.HexToAddress("0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984") // UNI
	addressB := common.HexToAddress("0x7Fc66500c84A76Ad7e9c93437bFc5Ac33E2DDaE9") // AAVE

	idA, err := tokenSystem.AddToken(addressA, "Uniswap", "UNI", 18)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added 'Uniswap' with stable ID: %d\n", idA)

	idB, err := tokenSystem.AddToken(addressB, "Aave", "AAVE", 18)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added 'Aave' with stable ID: %d\n", idB)

	view, err := tokenSystem.GetTokenByID(idA)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved token by ID %d: %s (%s)\n", view.ID, view.Name, view.Symbol)

	allTokens := tokenSystem.View()
	fmt.Printf("Total tokens in registry: %d\n", len(allTokens))

	err = tokenSystem.DeleteToken(idB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted token with ID: %d\n", idB)

	allTokens = tokenSystem.View()
	fmt.Printf("Total tokens after deletion: %d\n", len(allTokens))
}
```

---

## Architecture

The system is split into two clear layers:

* **`token.go` (TokenRegistry)**
  The internal engine, built with Data-Oriented Design (SoA, transformation functions). This layer is not concurrency-safe by itself.

* **`tokensystem.go` (TokenSystem)**
  The public, thread-safe API. Wraps `TokenRegistry` and ensures all operations are safe using `sync.RWMutex`.

This hybrid model ensures maximum performance internally while providing safety at the API level.

---

## Running Tests

Run all unit tests with race detection:

```bash
go test -v -race ./...
```

Run fuzz tests (recommended to run for at least a minute):

```bash
go test -fuzz=FuzzDeleteToken
```

Run benchmarks:

```bash
go test -bench=.
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
