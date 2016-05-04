// A tool to check params and take log for Weixin Pay Merchents.
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Configuration to load from json
type Configuration struct {
	// listen address, like:"0.0.0.0:80"
	Listen string
	// Key 秘钥，在商户平台的API安全里设置
	Key string
	// IsLoad 判断是否已经载入配置
	IsLoad bool
}

var config = Configuration{}

func main() {
	// init logger
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// log.SetOutput(&lumberjack.Logger{
	// 	Filename:   "log.txt",
	// 	MaxSize:    20,
	// 	MaxBackups: 9,
	// })
	// configFile Load config json.
	config.IsLoad = false
	configFile, _ := os.Open("config.json")
	decoder := json.NewDecoder(configFile)

	err := decoder.Decode(&config)
	if err != nil {
		log.Println("error:", err)
	}
	log.Println("listen address:", config.Listen)
	config.IsLoad = true
	err = http.ListenAndServeTLS(config.Listen, "server.crt", "server.key", http.HandlerFunc(mainHandler))
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// mainHandler HTTP Handler func.
func mainHandler(w http.ResponseWriter, req *http.Request) {
	//if req.Method != "POST" {
	log.Println(req.URL.String())
	body, err := ioutil.ReadAll(req.Body)
	tmp, err := XMLStructToMap(body)
	log.Printf("%#v", tmp)

	newReq, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewReader(body))
	newReq.Header = req.Header
	newReq.URL.Host = "api.mch.weixin.qq.com"
	newReq.Host = newReq.URL.Host
	//newReq.Body = body
	newReq.URL.Scheme = "https"

	client := &http.Client{}
	resp, err := client.Do(newReq)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(resp.StatusCode)

	defer resp.Body.Close()

	io.Copy(w, resp.Body)

	if err != nil {
		log.Fatal(err.Error())
	}

	return
	//}
}
