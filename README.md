## event-listener 
event-listener 用来解析 solidity 合约中`event XXX(a, b, c)`这种方式(即emit)产生的日志

- v1: 通过指定开始区块号分批爬块的方式, 处理event
- v2: 通过流的方式在线处理每一个transaction的日志 （todo）

