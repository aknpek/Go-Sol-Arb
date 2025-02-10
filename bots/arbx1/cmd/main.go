package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
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

// PREDEFINED POOLS WITH TOKEN DECIMALS
var POOL_CONFIGS = map[solana.PublicKey]PoolConfig{
	solana.MustPublicKeyFromBase58("58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2"): {
		TokenA: TokenConfig{
			Mint:     solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112"), // SOL
			Decimals: 9,
		},
		TokenB: TokenConfig{
			Mint:     solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"), // USDC
			Decimals: 6,
		},
	},
	solana.MustPublicKeyFromBase58("HJPjoWUrhoZzkNfRpHuieeFk9WcZWjwy6PBjZ81ngndJ"): {
		TokenA: TokenConfig{
			Mint:     solana.MustPublicKeyFromBase58("DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"), // BONK
			Decimals: 5,
		},
		TokenB: TokenConfig{
			Mint:     solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"), // USDC
			Decimals: 6,
		},
	},


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
		// Fetch and decode all target pools
		pools := make([]PoolReserves, 0)
		accounts, err := client.GetMultipleAccounts(ctx, TARGET_POOLS...)
		if err != nil { 
			log.Printf("Error fetching accounts: %v", err)
			time.Sleep(POLL_INTERVAL)
			continue
		}

		for i, account := range accounts.Value { 	
			if account == nil { // POOL account not found.
				continue
			}

			reserves, err := decodeRaydiumPool(account)

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
func decodeRaydiumPool(acc *rpc.Account) (PoolReserves, error) {
	data := acc.Data.GetBinary()
	if len(data) < 17 {
		return PoolReserves{}, fmt.Errorf("invalid pool data length")
	}

	config, exists := POOL_CONFIGS[acc.PublicKey]
	if !exists {
		return PoolReserves{}, fmt.Errorf("unknown pool configuration")
	}

	reserveA := binary.LittleEndian.Uint64(data[1:9])
	reserveB := binary.LittleEndian.Uint64(data[9:17])

	adjReserveA := float64(reserveA) / math.Pow10(int(config.TokenA.Decimals))
	adjReserveB := float64(reserveB) / math.Pow10(int(config.TokenB.Decimals))

	return PoolReserves{
		Address:  acc.PublicKey, // Use PublicKey not Owner
		ReserveA: adjReserveA,
		ReserveB: adjReserveB,
		Price:    adjReserveB / adjReserveA,
	}, nil
}


// Simple arbitrage check between two pools
func checkArbitrage(pool1, pool2 PoolReserves) {
	spread := math.Abs(pool1.Price - pool2.Price)
	if spread > 0.001 {
		config1 := POOL_CONFIGS[pool1.Address]
		config2 := POOL_CONFIGS[pool2.Address]
		
		log.Printf("ARBITRAGE DETECTED: %s (%.4f) vs %s (%.4f) [%.2f%%]",
			config1.TokenA.Mint.String(), pool1.Price,
			config2.TokenA.Mint.String(), pool2.Price,
			spread*100,
		)
	}
}
