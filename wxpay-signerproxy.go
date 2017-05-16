// A tool to check params and take log for Weixin Pay Merchents.
package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/beevik/ntp"
	"github.com/miekg/dns"
)

// Configuration to load from json
type Configuration struct {
	// 监听地址, like:"0.0.0.0:80"
	Listen struct {
		// http协议监听地址
		HTTP string
		// https协议监听地址
		HTTPS string
	}
	// Key 秘钥，在商户平台的API安全里设置
	Key string
	// UseCert 是否使用证书
	UseCert bool `json:"use_cert"`
	// 证书配置
	Cert struct {
		// CertFile 证书文件路径
		CertFile string
		// KeyFile 秘钥文件路径
		KeyFile string
		// Ca CA证书路径
		Ca string
	}
	// Resign 是否重签名
	Resign bool
	// LogToFile 是否记录到文件
	LogToFile bool `json:"log_to_file"`
	// Diagnosis 是否诊断
	Diagnosis bool
	// IsLoad 判断是否已经载入配置
	IsLoad bool
}

var config Configuration
var cert tls.Certificate
var caCert []byte

func init() {
	// init logger
	config.IsLoad = false
	configFile, _ := os.Open("config.json")
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)

	err := decoder.Decode(&config)
	if err != nil {
		log.Println("config error:", err)
	}
	log.Println("listen address:", config.Listen.HTTP, config.Listen.HTTPS)
	config.IsLoad = true

	if config.LogToFile {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		l := &lumberjack.Logger{
			Filename:  "log.txt",
			MaxSize:   20,
			LocalTime: true,
		}
		l.Rotate()
		log.SetOutput(l)
	}

	//载入证书
	if config.UseCert {
		// Load client cert
		cert, err = tls.LoadX509KeyPair(config.Cert.CertFile, config.Cert.KeyFile)
		if err != nil {
			log.Fatal(err)
		}

		// Load CA cert
		caCert, err = ioutil.ReadFile(config.Cert.Ca)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	if config.Diagnosis {
		// do some test here
		// 0. Collect Infos
		log.Println("[INFO] ============")
		log.Println("[INFO] 0. Address Info")
		localAddrs, err := net.InterfaceAddrs()
		log.Printf("%v", localAddrs)
		// 1. DNS lookup test
		log.Println("[INFO] ============")
		log.Println("[INFO] 1. DNS Lookup test")
		addr, err := net.LookupHost("api.mch.weixin.qq.com")
		log.Printf("[INFO] Local DNS lookup result: %v", addr)

		c := dns.Client{}
		m := dns.Msg{}
		m.SetQuestion("api.mch.weixin.qq.com.", dns.TypeA)
		r, t, err := c.Exchange(&m, "119.29.29.29:53")
		if err != nil {
			log.Printf("[FATAL] Get ip from 119.29.29.29 failed. err:%s", err)
		}
		log.Printf("[INFO] Request took time: %v", t)
		if len(r.Answer) == 0 {
			log.Fatal("No results")
		}
		log.Printf("[INFO] Remote DNS lookup result:%v", r)
		//2. local time and remote time
		log.Println("[INFO] ============")
		log.Println("[INFO] 2. Time test")
		localTime := time.Now()
		log.Printf("[INFO] Local Time:%v", localTime.Unix())
		log.Println("[INFO] Getting remote time, should not wait too long")
		remoteTime, err := ntp.Time("cn.pool.ntp.org")
		if err != nil {
			log.Printf("[FATAL] Get time from cn.pool.ntp.org failed. err:%s", err)
		}
		log.Printf("[INFO] Remote Time:%v", remoteTime.Unix())
		if math.Abs(float64(localTime.Unix()-remoteTime.Unix())) > 10 {
			log.Printf("[WARNING] Time needs Sync:%v", remoteTime.Unix())
		}

		log.Println("[INFO] Test done.")
	}
	go func() {
		err2 := http.ListenAndServe(config.Listen.HTTP, http.HandlerFunc(mainHandler))
		if err2 != nil {
			log.Fatal("ListenAndServe:", err2)
		}
	}()

	err := http.ListenAndServeTLS(config.Listen.HTTPS, "server.crt", "server.key", http.HandlerFunc(mainHandler))
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// mainHandler HTTP Handler func.
func mainHandler(w http.ResponseWriter, req *http.Request) {
	//if req.Method != "POST" {
	log.Println(req.URL.String())
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	mapData, err := XMLStructToMap(body)
	log.Printf("[INFO][IN] Request content:%v", string(body))

	// sign
	eq, err := CompareSignature(body)
	if !eq {
		log.Println("[WARNING] Sign not match.")
		log.Printf("original sign: %s", mapData["sign"])
		newSign, err2 := CalculateSignature(mapData, config.Key)
		log.Printf("calculated sign: %s", newSign)
		if config.Resign {
			mapData["sign"] = newSign
			body, err2 = MapToXML(mapData)
			if err2 != nil {
				log.Fatal(err.Error())
			}
			log.Printf("[INFO][IN] Renewed request content:%v", string(body))
		}
	}

	newReq, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewReader(body))
	newReq.Header = req.Header
	newReq.URL.Host = "api.mch.weixin.qq.com"
	newReq.Host = newReq.URL.Host
	newReq.URL.Scheme = "https"

	var client *http.Client
	//cert
	if config.UseCert {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client = &http.Client{Transport: transport}
	} else {
		client = &http.Client{}
	}

	resp, err := client.Do(newReq)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(resp.StatusCode)

	//io.Copy(w, resp.Body)
	respBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[INFO][IN] Respond content:%v", string(respBody))
	w.Write(respBody)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	return
	//}
}
