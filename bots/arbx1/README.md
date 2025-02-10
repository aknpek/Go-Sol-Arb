# Jito Searcher Setup
**Answer:**
- The jito.NewSearcher setup does not require manual gRPC installation. The Jito-Go SDK handles gRPC internally.

**Cost:**
- Free to use: Connecting to Jito’s Block Engine (via amsterdam.block-engine.jito.wtf:443) is free.

- Transaction costs: You pay only for bundles via tips (≥1,000 lamports) and Solana network fees.

- Purpose: This setup is for submitting MEV bundles (e.g., arbitrage transactions) with priority.

# Private RPC Node vs. External RPC Endpoints
- Jito Block Engine: Required for MEV bundle submission (separate from standard RPC).

- DEX APIs: Still useful for metadata (e.g., token names, pool types) without parsing raw data.

# Arbitrage Options

1. Why Monitor 1,000+ Pools?

- Arbitrage opportunities emerge from price discrepancies across pools. The more pools you monitor, the higher your chances of finding profitable trades. Here’s why:

**Factor Explanation**
- Opportunity Density More pools = More token pairs (e.g., SOL/USDC, BONK/USDT) = More chances for price differences.
Liquidity Spread Low-liquidity pools often have larger spreads (e.g., 5% vs 0.3% on high-liquidity pools).
DEX Diversity	Different DEXes (Raydium, Orca, Meteora) may price the same token differently.

1. Implement triangular arbitrage (e.g., SOL → USDC → BONK → SOL).

2. Use Jito’s tip system instead of front-running small trades.

3. Batch Requests (getMultipleAccounts):

Effectiveness: Fetching 100 pools in 1 RPC call reduces requests by 99%.
Limitation: Still requires parsing all accounts. Combine with caching:

```go
// Batch fetch 100 pools
accounts, _ := client.GetMultipleAccounts(ctx, poolAddresses)
// Decode and cache all
for _, acc := range accounts {
    decoded := decodePool(acc)
    redis.Set(acc.PublicKey.String(), decoded, 100*time.Millisecond)
}
```


**Example Workflow:**
- Geyser sends a pool reserve update.
- Bot decodes and processes the data.
- Bot stores the processed reserves in Redis.
- Subsequent requests for reserves (within 100ms) use the cached data.