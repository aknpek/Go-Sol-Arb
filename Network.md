
## UDP buffer size 

Current:
udp_rmem_min = 4096 (4KB) ❌ Too small!

Recommended:
net.ipv4.udp_rmem_min = 16777216    # 16MB
net.core.rmem_max = 134217728        # 128MB
net.core.rmem_default = 134217728    # 128MB


# Set each value individually
sudo sysctl -w net.core.rmem_max=134217728
sudo sysctl -w net.core.wmem_max=134217728
sudo sysctl -w net.core.rmem_default=134217728
sudo sysctl -w net.ipv4.udp_rmem_min=16777216


# Monitor

- Monitor established connections
```bash
watch -n 1 'netstat -ant | grep :8899 | grep ESTABLISHED | wc -l'
```

- Check buffer pressure 
```bash
netstat -s | grep -i buffer
```

- Watch memory usage
```bash
free -m 
```

- Check network interface stats
```bash
ip -s link show ens5
```

- Monitor Network throughput
```bash
iftop -i ens5
```