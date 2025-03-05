### Find devices

Encontra dispositivos Daikin na rede local.

```
λλ go run *.go --timeout=10s
Device Hostname=DAIKINXXAXX.local. IP=192.168.0.XX Port=80 APN=DAIKIN:XXAXXXXXXXXC

JSON format 
λλ go run *.go --timeout=10s --json=true
[
  {
    "port": 80,
    "apn": "DAIKIN:XXAXXXXXXXXC",
    "hostname": "DAIKINXXAXXX.local.",
    "ip": "192.168.0.XX"
  }
]
```