# How to Connect to Ec2

- The permissions on the private key are too open (0644), therefore, we need to give the correct permissions.
```bash
chmod 400 ./solana-bot-key.pem
```

- If first time connection being established, please say "yes"
```bash
ssh -i ./solana-bot-key.pem ubuntu@44.195.40.169
```

### Solana RPC Client

- Check the solana version
```bash
solana --version
```

### Memory Check 

```bash
free -h
```

### CPU Check 

```bash
lscpu
```

### Test RPC Node Run
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

# RPC Running Validator Settings: 

### 1. System Service File
- Open the Validator service
```bash
sudo nano /etc/systemd/system/solana-validator.service
```
- Do You Need --full-rpc-api for an Arbitrage Bot?

No, unless your bot specifically requires historical transaction data or obscure RPC methods. Most arbitrage bots only need:
- getRecentBlockhash
- sendTransaction
- getAccountInfo
- getSlot


```ini
[Unit]
Description=Solana RPC Validator
After=network.target remote-fs.target syslog.target
Wants=network.target remote-fs.target

[Service]
Type=simple
User=sol
WorkingDirectory=/home/sol
ExecStart=/home/sol/bin/validator.sh
Restart=always
RestartSec=10
LimitNOFILE=1000000
LimitNPROC=500000
LimitMEMLOCK=infinity
StandardOutput=append:/home/sol/solana-validator.log
StandardError=append:/home/sol/solana-validator.log

# ✅ Memory and CPU Usage Limits
MemoryMax=220G       # Hard limit (validator will be killed if it exceeds)
MemoryHigh=180G      # Soft limit (triggers memory cleanup)
CPUQuota=80%         # Restrict to 80% of total CPU
IOSchedulingClass=2  # Improves disk I/O priority

[Install]
WantedBy=multi-user.target
```

### 2. Bash Script; 

**Where is the RPC Setup Located**
- If you run the bash script it will start with the flags
- When the script runs, it starts the Solana validator as a foreground process (attached to your terminal session)
- If you close the terminal, the validator will stop running

```bash
sudo nano /home/sol/bin/validator.sh
```

**Running from the Bash Script**

<details>
<summary>Bash Script Details Here</summary>


**Updated Version v0.0.1**
```validator.sh
#!/bin/bash
exec /home/sol/.local/share/solana/install/releases/2.1.5/solana-release/bin/agave-validator \
    --identity /home/sol/validator-keypair.json \
    --known-validator 7Np41oeYqPefeNQEHSv1UDhYrehxin3NStELsSKCT4K2 \
    --known-validator GdnSyH3YtwcxFvQrVVJMm1JhTS4QVX7MFsX56uJLUfiZ \
    --known-validator DE1bawNcRJB9rVm3buyMVfr8mBEoyyu73NBovf2oXJsJ \
    --known-validator CakcnaRDHka2gXyfbEd2d3xsvkJkqsLw2akB3zsN1D2S \
    --only-known-rpc \
    --full-rpc-api \
    --no-voting \
    --ledger /mnt/ledger \
    --accounts /mnt/accounts \
    --log /home/sol/solana-rpc-mainnet.log \
    --rpc-port 8899 \
    --rpc-bind-address 0.0.0.0 \
    --private-rpc \
    --dynamic-port-range 8000-8020 \
    --entrypoint entrypoint.mainnet-beta.solana.com:8001 \
    --entrypoint entrypoint2.mainnet-beta.solana.com:8001 \
    --entrypoint entrypoint3.mainnet-beta.solana.com:8001 \
    --entrypoint entrypoint4.mainnet-beta.solana.com:8001 \
    --entrypoint entrypoint5.mainnet-beta.solana.com:8001 \
    --expected-genesis-hash 5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d \
    --wal-recovery-mode skip_any_corrupted_record \
    --limit-ledger-size \
    --enable-rpc-transaction-history \
    --enable-cpi-and-log-storage
```

**Updated Version v0.0.2**

```bash
exec /home/sol/.local/share/solana/install/releases/2.1.5/solana-release/bin/agave-validator \
    --identity /home/sol/validator-keypair.json \
    --known-validator 7Np41oeYqPefeNQEHSv1UDhYrehxin3NStELsSKCT4K2 \
    --known-validator GdnSyH3YtwcxFvQrVVJMm1JhTS4QVX7MFsX56uJLUfiZ \
    --known-validator DE1bawNcRJB9rVm3buyMVfr8mBEoyyu73NBovf2oXJsJ \
    --known-validator CakcnaRDHka2gXyfbEd2d3xsvkJkqsLw2akB3zsN1D2S \
    --only-known-rpc \
    --full-rpc-api \
    --no-voting \
    --ledger /mnt/ledger \
    --accounts /mnt/accounts \
    --log /home/sol/solana-rpc-mainnet.log \
    --rpc-port 8899 \
    --rpc-bind-address 0.0.0.0 \
    --private-rpc \
    --dynamic-port-range 8000-8020 \
    --entrypoint entrypoint.mainnet-beta.solana.com:8001 \
    --entrypoint entrypoint2.mainnet-beta.solana.com:8001 \
    --entrypoint entrypoint3.mainnet-beta.solana.com:8001 \
    --entrypoint entrypoint4.mainnet-beta.solana.com:8001 \
    --entrypoint entrypoint5.mainnet-beta.solana.com:8001 \
    --expected-genesis-hash 5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d \
    --wal-recovery-mode skip_any_corrupted_record \
    --limit-ledger-size \
    --enable-rpc-transaction-history \
    --enable-cpi-and-log-storage \
    --accounts-db-cache-limit-mb=8192 \
    --accounts-index-memory-limit-mb=16384
```

</details>

### How to start RPC from Bash

```bash
sudo /home/sol/bin/validator.sh
```

- It can be started from the binary or validator.sh \
✅ Starts the validator with Solana's default settings. \
❌ Might not include all optimizations and custom configurations. \
❌ Requires manual input of flags each time you start it. 

- Manually starting the bash script \
✅ Uses all your custom configurations (like ledger location, RPC settings, entry points). \
✅ Ensures every startup is consistent with the same parameters. \
✅ Can be automated using systemd to restart on failures.

#### Should You Always Use validator.sh?
🔹 YES, if you want full control over your setup and better performance.
🔹 If you run it manually (binary-only), you may miss important configurations.
🔹 If you restart the server, a systemd-managed validator.sh will ensure auto-restart. 


### RPC Reco
✅ One important fact that rotating the log files that will not overload the disk space. \
✅ Limit memory usage in systemd (MemoryMax=220G, CPUQuota=80%). \
✅ Enable log rotation to prevent logs from filling disk. \
✅ Enable 32GB swap space to avoid crashes (sudo fallocate -l 32G /swapfile). \
✅ Remove --enable-rpc-transaction-history from validator.sh to save RAM.

🔹 Explanation of New Limits
✅ MemoryMax=220G → If Solana exceeds 220GB RAM, it will be killed to prevent system crashes. \
✅ MemoryHigh=180G → Triggers memory cleanup when usage crosses 180GB. \
✅ CPUQuota=80% → Limits validator to 80% CPU usage, avoiding system overload. \
✅ IOSchedulingClass=2 → Optimizes disk access, reducing latency issues.


### RPC SyncUP Status Check
```bash
sudo tail -f /home/sol/solana-rpc-mainnet.log
```
- Check the RPC sync up status after the private RPC stop running.
- After Slot re-building is done, RPC should start
```bash
netstat -tulnp | grep 8899
```

### Check Run Status

1. Check if the systemd is handling the validator
```bash
sudo systemctl status solana-validator
```

2. If it is not active, enable it.
```bash
sudo systemctl daemon-reload
sudo systemctl enable solana-validator
sudo systemctl restart solana-validator
```

### Check Solana Validator Progress

```bash
journalctl -u solana-validator | less
```
- Monitor the progress of the validator history.

```bash
sudo journalctl -u solana-validator -f --since "2m ago"
```

## RPC Node Usages
### Check the Storages
- Check the file storages
```bash
df -h
```

### Check Memory Usage
- Check the memory usages
```bash
top
```

### Create a Swap to reduce of getting OOM Killed
- Create swap memory from the storage
```bash
sudo fallocate -l 64G /mnt/swapfile
sudo chmod 600 /mnt/swapfile
sudo mkswap /mnt/swapfile
sudo swapon /mnt/swapfile
```

2. What is 
```bash
free -h
```

## RPC Errors

## How to Check Error Logs of the Validator:

```bash
sudo tail -n 1000 /home/sol/solana-rpc-mainnet.log | grep -i "error\|fatal\|panic"
```
### Can I limit the Slots Downloaded? 
- Run a minimal validator in --limit-ledger-size mode. This only prunes old blocks after you’ve verified them, but you can’t skip verifying them in the first place.
- The --limit-ledger-size flag does help keep disk usage lower, but you still must download and replay each slot to confirm the chain’s state at some point in time. After you’ve processed them, the node discards older ledger data to save space, but it doesn’t skip them entirely.


### 1. Add Solana CLI to Secure Path
The solana-keygen command is a key management tool for Solana. It helps generate and inspect keypairs used for validators, wallets, and accounts.

- How to get the our validator public key? 
```bash
sudo /home/ubuntu/.local/share/solana/install/active_release/bin/solana-keygen pubkey /home/
sol/validator-keypair.json
```

### 2. Network Buffer Size Errors (Keep the Large Network Buffers?)
- (512gb machine) Since you have abundant memory, you can keep the net.core.rmem_max, net.core.wmem_max, etc. set high. This helps prevent packet drops under heavy network 
load.

```bash
# Set 1GB buffers (maximum supported)
sudo sysctl -w net.core.rmem_max=1073741824
sudo sysctl -w net.core.rmem_default=1073741824
sudo sysctl -w net.core.wmem_max=1073741824
sudo sysctl -w net.core.wmem_default=1073741824

# Set TCP specific buffers
sudo sysctl -w net.ipv4.tcp_rmem='4096 87380 1073741824'
sudo sysctl -w net.ipv4.tcp_wmem='4096 87380 1073741824'
```

```bash
sudo bash -c "cat >/etc/sysctl.d/21-agave-validator.conf <<EOF
# Increase UDP buffer sizes
net.core.rmem_default = 134217728
net.core.rmem_max = 134217728
net.core.wmem_default = 134217728
net.core.wmem_max = 134217728

# Increase memory mapped files limit
vm.max_map_count = 1000000

# Increase number of allowed open file descriptors
fs.nr_open = 1000000
EOF"
```

#### Temporary Adjustment
```bash
sudo sysctl -w net.core.netdev_max_backlog=50000
sudo sysctl -w net.ipv4.tcp_max_syn_backlog=50000
```

#### Persistent Adjustment

```bash
sudo bash -c "cat >>/etc/sysctl.d/99-network.conf <<EOF
# Increase backlog queues
net.core.netdev_max_backlog=50000
net.ipv4.tcp_max_syn_backlog=50000
EOF"
```

### 3. Permission Denied Errors with RPC

Running as a Root vs Sol User
- Right now, it's running as root instead of the sol user. This could cause permission issues with files inside /home/sol/.


```bash
ted Solana Validator.
Feb 13 08:40:13 ip-172-31-9-222 systemd[8819]: solana-validator.service: Failed to determine user credentials: No such process
Feb 13 08:40:13 ip-172-31-9-222 systemd[8819]: solana-validator.service: Failed at step USER spawning /home/ubuntu/bin/validator.sh: No such process
Feb 13 08:40:13 ip-172-31-9-222 systemd[1]: solana-validator.service: Main process exited, code=exited, status=217/USER
Feb 13 08:40:13 ip-172-31-9-222 systemd[1]: solana-validator.service: Failed with result 'exit-code'.
Feb 13 08:40:14 ip-172-31-9-222 systemd[1]: solana-validator.service: Scheduled restart job, restart counter is at 6.
Feb 13 08:40:14 ip-172-31-9-222 systemd[1]: Stopped Solana Validator.
Feb 13 08:40:15 ip-172-31-9-222 systemd[1]: Started Solana Validator.
```

- Check permissions SOL users
```bash
ls -l /home/sol/bin/validator.sh  # Check script permissions
```

- Give the permissions
```bash
# Fix ownership and permissions for the sol user
sudo chown -R sol:sol /home/sol  # Ensure sol owns its home directory
sudo chmod 755 /home/sol  # Allow directory traversal
sudo chmod +x /home/sol/bin/validator.sh  # Make the script executable
```

- Fix validator.sh Paths
    - Environment="PATH=...:/home/sol/.local/share/solana/install/active_release/bin"

```bash
sudo -u sol ln -sf /home/sol/.local/share/solana/install/releases/2.1.5/solana-release /home/sol/.local/share/solana/install/active_release
```
