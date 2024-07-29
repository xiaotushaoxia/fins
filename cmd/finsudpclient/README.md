# simple fins client

# usage
```
Usage of finsclient.exe:
  -cn int
        client network(0-255)
  -cnd int
        client node(0-255)   
  -cu int
        client unit(0-255)
  -ip string
        plc server ip (default "127.0.0.1")
  -port int
        plc server udp port (default 9600)
  -sn int
        plc server network(0-255)
  -snd int
        plc server node(0-255)
  -su int
        plc server unit(0-255)

```

# Output
```bash
>> help
support memory type: D for DM Area, A for Auxiliary Area, H for Holding Bit Area, W for Work Area
support data type: b for Bit, B for Byte, s for String, w for Word
read usage:  r <memory type> <data type> <address> <count>
write usage: w <memory type> <data type> <address> <values>
set/reset usage: set/reset <memory type> <address> <offset>
single cmd usage: `close` for close client conn; `rc` for read clock
>> w A w 100 1 2 3 4 5 6
write success
>> r A w 100 6
read success:  [1 2 3 4 5 6]
>> r A B 100 6
read success:  [0 1 0 2 0 3 0 4 0 5 0 6]
```