package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// CONFIGURATION
const (
	RPC_URL = "http://localhost:8899" 
	RAYDIUM_PROGRAM_ID = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
	POLL_INTERVAL = 10 * time.Second
)

// PREDEFINED POOLS TO MONITOR (E.G SOL/USDC, BONK/USDC, ...)
var TARGET_POOLS = []solana.PublicKey{
	solana.MustPublicKeyFromBase58("58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2"), // SOL/USDC
	solana.MustPublicKeyFromBase58("HJPjoWUrhoZzkNfRpHuieeFk9WcZWjwy6PBjZ81ngndJ"),  // BONK/USDC
}

type PoolReserves struct { 
	Address solana.PublicKey
	ReserveA float64 // Token A (e.g. SOL)
	ReserveB float64 // Token B (e.g. USDC)
	Price float64 // Price = Reserve B / Reserve A
}

func main() {
	client := rpc.New(RPC_URL)
	ctx := context.Background()

	for { 
		// Fetch and decade all target pools
		pools := make([]PoolReserves, 0)
		accounts, err := client.GetMultipleAccounts(ctx, TARGET_POOLS...)
		if err != nil { 
			log.Printf("Error fetching accounts: %v", err)
			time.Sleep(POLL_INTERVAL)
			continue
		}

		for _, account := range accounts.Value { 	
			if account == nil { // POOL account not found.
				continue
			}

			reserves, err := decodeRaydiumPool(*acc)

			if err != nil { 
				log.Printf("Error decoding pool: %v", err)
				continue
			}
			
			pools = append(pools, reserves)
			log.Printf("Pool %s: 1 TokenA = %.4f TokenB", TARGET_POOLS[i].String(), reserves.Price)

		}

		// Check for arbitrage opportunities between two pools
		if len(pools) >= 2 { 
			checkArbitrage(pools[0], pools[1])
		}

		time.Sleep(POLL_INTERVAL)
	}

}

// Decode Raydium pool reserves from account data
func decodeRaydiumPool(acc rpc.Account) (PoolReserves, error) { 

	data := acc.Data.GetBinary()
	if len(data) < 17 { 
		return PoolReserves{}, fmt.Errorf("invalid pool data length")
	}

	// Reserve A (u64 at offset 1-8)
	reserveA := solana.Uint64FromBytes(data[1:9])

	// Reserve B (u64 at offset 9-17)
	reserveB := solana.Uint64FromBytes(data[9:17])

	// Adjust for token decimals (e.g. SOL=9, USDC=6)
	adjReserve := float64(reserveA) / 1e9
	adjReserveB := float64(reserveB) / 1e6

	return PoolReserves{
		Address: acc.PublicKey,
		ReserveA: adjReserve,
		ReserveB: adjReserveB,
		Price: adjReserveB / adjReserve,
	}, nil
}


// Simple arbitrage check between two pools
func checkArbitrage(pool1, pool2 PoolReserves) { 
	spread := abs(pool1.Price - pool2.Price)

	if spread > 0.001 { // .1% threshold
		log.Printf("ARBITRAGE DETECTED: %.4f vs %.4f (%.2f%%)", pool1.Price, pool2.Price, spread * 100)
	}
}

func abs(x float64) float64 { 
	if x < 0 { 
		return -x
	}
	return x
}