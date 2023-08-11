
1. 添加默认路由（一种特殊的静态路由），下一跳192.168.116.1，metric 为1
```
  ip route add default via 192.168.116.1 metric 1
      or
  ip route add 0.0.0.0/0 via 192.168.116.1 metric 1
```
2. 添加静态路由 192.168.208.0 网段的流量下一跳 192.168.116.1 metric 1
```
  ip route add 192.168.208.0/24 proto static via 192.168.116.1 metric 1
```
3. 添加直连路由，出接口ip为192.168.116.108
```
  ip route add 192.168.116.0/24 proto kernel src 192.168.116.108
```
4. 添加主机路由，出接口ip为192.168.116.108，网卡名为eth0
```
  ip route add 192.168.116.108/32 dev eth0
```