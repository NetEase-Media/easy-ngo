# accesslog 格式说明
格式遵循apache语法：
- %a - 远程IP地址
- %A - 本地IP地址
- %b - 发送的字节数(Bytes sent), 不包括HTTP headers的字节，如果为0则展示'-'
- %B - 发送的字节数(Bytes sent), 不包括HTTP headers的字节
- %h - 远程主机名称(如果resolveHosts为false则展示IP)
- %H - 请求协议
- %l - 远程用户名，始终为'-'
- %m - 请求的方法(GET, POST等)
- %p - 接受请求的本地端口`暂不支持`
- %q - 查询字符串，如果存在，有一个前置的'?'
- %r - 请求的第一行(包括请求方法和请求的URI)
- %s - response的HTTP状态码(200,404等)
- %S - 用户的session ID`暂不支持`
- %t - 日期和时间，Common Log Format格式
- %u - 被认证的远程用户, 不存在则展示'-'
- %U - 请求URL路径
- %v - 本地服务名`暂不支持`
- %D - 处理请求的时间，单位为毫秒
- %T - 处理请求的时间，单位为秒
- %I - 当前请求的线(协)程名`暂不支持`

另外，Access Log中也支持cookie，请求header，响应headers，Session或者其他在ServletRequest中的对象的信息。格式遵循apache语法：
- %{xxx}i 请求headers的信息
- %{xxx}o 响应headers的信息
- %{xxx}c 请求cookie的信息
- %{xxx}r xxx是gin.Context.Keys的一个key 

内置别名
- common = `%h %l %u %t "%r" %>s %b`
- combined = `%h %l %u %t "%r" %>s %b "%{Referer}i" "%{User-agent}i"`