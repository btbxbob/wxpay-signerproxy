# wxpay-signerproxy

一个可以重签名传入xml，并转发给api.mch.weixin.qq.com的反向代理
并且集成了一些常见的诊断功能，包括：

- 检查传入的签名是否匹配
- 检查对api.mch.weixin.qq.com的DNS解析是否和119.29.29.29的一致
- 检查本地时间和cn.pool.ntp.org上的Ntp服务的是否一致
- 滚动日志输出

# 使用说明

1. 创建一个配置文件`config.json`，格式如下:
```
    {
        "listen": {
            "https": "127.0.0.1:1718",
            "http": "127.0.0.1:1717"
        },
        "key": "--填写商户号对应的key--",
        "use_cert": true,
        "cert": {
            "certfile": "--证书文件路径--",
            "keyfile": "--秘钥文件路径--",
            "ca": "--根证书--"
        },
        "resign": true,
        "log_to_file": false,
        "diagnosis": true
    }
```
  释义见下：
```
type Configuration struct {
	// 监听地址, like:"0.0.0.0:80"
	Listen struct {
    // http协议监听地址
		HTTP  string
    // https协议监听地址
		HTTPS string
	}
	// Key 秘钥，在商户平台的API安全里设置
	Key string
	// UseCert 是否使用证书
	UseCert bool `json:"use_cert"`
  // 证书配置
	Cert    struct {
    // 证书文件路径
		CertFile string
    // 秘钥文件路径
		KeyFile  string
    // CA证书路径
		Ca       string
	}
  // Resign 是否重签名
	Resign bool
	// LogToFile 日志功能开关
	LogToFile bool `json:"log_to_file"`
	// Diagnosis 诊断功能开关
	Diagnosis bool
}
```

2. 运行wxpay-signerproxy

3. 将远系统代码请求的接口地址中的api.mch.weixin.qq.com更换为监听地址，这里要注意协议。  
   由于这个工具本身没有有效的根证书签名，https请求的话需要忽略服务端证书有效性。或者改为使用http协议。

> 从wxpay-signerproxy到api.mch.weixin.qq.com的请求是https协议的。

4. 进行支付，并收集日志。

# 待加入的特性

- 诊断请求延迟
- 检查证书

# wxpay-signerproxy
A tool too resign incoming request and rewrite to wxpay server.

- Check sign
- Check DNS record
- Check local time

# TODO

- Check delay
- Check cert
