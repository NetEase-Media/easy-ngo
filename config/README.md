config

Protocol
- etcd:
- file:
- http:
- env:
- consul:
- ftp:

config模块在加载的时候，读取启动参数-c，解析-c参数，根据不同的协议，调用不同的实现