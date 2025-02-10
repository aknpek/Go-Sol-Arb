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

1. Check if go installed;

```bash
go 
```

2. Cloning the bot code;
```bash
# Create a projects directory in your home folder
mkdir -p ~/projects/solana-bots

# Move into that directory
cd ~/projects/solana-bots

# Clone your repository here
git clone https://github.com/aknpek/Go-Sol-Arb.git
```

## Your directory structure will look like:
/home/ubuntu/ \
└── projects/ \
    └── solana-bots/ \
        └── your-repo-name/ \
            ├── bots/ \
            │   └── arbx1/ \
            │       ├── cmd/ \
            │       │   └── main.go \
            │       ├── go.mod \ 
            │       ├── go.sum \ 
            │       └── README.md \
            └── .gitignore \



## How to run?

1. Install the Bot dependencies
```bash
go mod tidy
```

2. Build and run the bot
```bash
cd bots/
```

## Possible Errors

1. Connection Refused

```bash
curl: (7) Failed to connect to localhost port 8899 after 0 ms: Connection refused
```

- Open in 2nd terminal
```bash
solana-test-validator
```