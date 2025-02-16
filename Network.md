

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