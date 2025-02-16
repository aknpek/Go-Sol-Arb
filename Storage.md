# Error on All data on root drive
- Properly mounted n NVMe driver
- Root drive relieved
- Better system performance


Snapshots are:
- Point-in-time copies of account states
- Used for quick node restart
- Two types:
  1. Full snapshots (complete state)
  2. Incremental snapshots (changes only)

Location:
/mnt/ledger/snapshots/
/mnt/accounts/snapshots/

Storage Setup:
- NVMe Drives:
  1. nvme0n1: 500GB (Root drive)
  2. nvme1n1: 1TB (Accounts data)
  3. nvme2n1: 2TB (Ledger data)

Mount Points:
/mnt/accounts -> nvme1n1
- Used for: Account states, snapshots
- Size: 984G
- Used: 434G
- Available: 505G

/mnt/ledger -> nvme2n1
- Used for: Blockchain data, snapshots
- Size: 2.0T
- Used: 471G
- Available: 1.4T

Persistence:
Added to /etc/fstab:
/dev/nvme1n1 /mnt/accounts ext4 defaults 0 0
/dev/nvme2n1 /mnt/ledger ext4 defaults 0 0

Directory Structure:
/mnt/accounts/
├── run/        # Active accounts data
└── snapshot/   # Account snapshots

/mnt/ledger/
├── accounts_hash_cache/
├── accounts_index/
├── snapshots/
├── remote/
├── rocksdb/
└── other files...

Permissions:
- Owner: sol:sol
- Mode: 755 (rwxr-xr-x)

Useful Commands:
1. Check Drive Format:
```bash
# Check if drives are formatted
sudo file -s /dev/nvme1n1
sudo file -s /dev/nvme2n1
```

2. Mount Drives:
```bash
# Create mount points
sudo mkdir -p /mnt/accounts
sudo mkdir -p /mnt/ledger

# Mount drives
sudo mount /dev/nvme1n1 /mnt/accounts
sudo mount /dev/nvme2n1 /mnt/ledger
```

3. Set Permissions:
```bash
# Set ownership
sudo chown -R sol:sol /mnt/accounts
sudo chown -R sol:sol /mnt/ledger

# Set permissions
sudo chmod 755 /mnt/accounts
sudo chmod 755 /mnt/ledger
```

4. Check Storage Usage:
```bash
# Check mount points and space
df -h

# Check specific directories
du -h --max-depth=1 /mnt/accounts
du -h --max-depth=1 /mnt/ledger
```

5. Verify Service Access:
```bash
# Check service user access
sudo -u sol ls -la /mnt/ledger/snapshots/
sudo -u sol ls -la /mnt/accounts/

# Check service status
sudo systemctl status solana-validator
```

Performance Monitoring:
1. Monitor I/O Performance:
```bash
# Monitor disk I/O in real-time
iostat -xm 5 nvme1n1 nvme2n1

# Check I/O wait time
vmstat 1

# Monitor disk throughput
sudo iotop -o -P
```

2. Monitor Disk Latency:
```bash
# Check disk latency
sudo apt install ioping
ioping -c 10 /mnt/accounts
ioping -c 10 /mnt/ledger

# Detailed latency stats
sudo apt install blktrace
sudo blktrace -d /dev/nvme1n1 -o - | blkparse -i -
```

3. Check Disk Health:
```bash
# Check SMART status
sudo apt install smartmontools
sudo smartctl -a /dev/nvme1n1
sudo smartctl -a /dev/nvme2n1

# Monitor temperature
sudo nvme smart-log /dev/nvme1n1
sudo nvme smart-log /dev/nvme2n1
```

4. Monitor File System:
```bash
# Watch directory growth
watch -n 60 'du -sh /mnt/accounts/* /mnt/ledger/*'

# Monitor inode usage
df -ih /mnt/accounts /mnt/ledger

# Check largest files/directories
sudo find /mnt/ledger -type f -exec du -Sh {} + | sort -rh | head -n 20
sudo find /mnt/accounts -type f -exec du -Sh {} + | sort -rh | head -n 20
```

5. Performance Alerts:
```bash
# Set up disk space alert (add to crontab)
#!/bin/bash
THRESHOLD=90
USE=$(df /mnt/accounts | grep -vE '^Filesystem|tmpfs' | awk '{ print $5 }' | sed 's/%//g')
if [ $USE -gt $THRESHOLD ]; then
    echo "Warning: /mnt/accounts disk space is ${USE}%" | mail -s "Disk Space Alert" admin@example.com
fi
```

## Disk Writes

```bash
sudo apt-get install iotop 
```

```bash
sudo iotop -o -P -b -n 1
```

# Read Ahead on Disk
1. Adjust the Read Ahead for the Replay Phase
- Current setting: 256 sectors
- Each sector = 512 bytes
- Total read-ahead = 128 kb

2. Check the Read Ahead
Check Specific Highly IOWaited Volume
```bash
sudo blockdev --getra /dev/nvme2n1
```

Purpose: 
- Prefetches data into cache
- Reduces I/O operations
- Improves sequential read performance
- For private RPC replay phase, increasing the read-ahead can help significantly, because of heavy sequential reads from the ledger. 
- High I/O wait times (as we saw 12.2) 
- Lots of accounts loading 

3. Update after sync 

```bash
sudo blockdev --setra 256 /dev/nvme2n1 
```