# How to Connect to Ec2

- The permissions on the private key are too open (0644), therefore, we need to give the correct permissions.
```bash
chmod 400 ./solana-bot-key.pem
```

- If first time connection being established, please say "yes"
```bash
ssh -i ./solana-bot-key.pem ubuntu@44.195.40.169
```

# Solana RPC Client

- Check the solana version
```bash
solana --version
```

## Memory Check 

```bash
free -h
```

## CPU Check 

```bash
lscpu
```

## Test RPC Node Run
- CURL command test.
```bash
curl http://localhost:8899 -X POST -H "Content-Type: application/json" -d '
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "getBlockHeight"
}
'
```
- Expected Response Shape
```json
{"jsonrpc":"2.0","result":297999981,"id":1}
```

# Solana Arbitrage Bot

## Prerequisites
- Go 1.x
- Solana CLI tools

## Installation