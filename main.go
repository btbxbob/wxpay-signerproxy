package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

type Configuration struct {
	Listen string
}

func main() {
	//read json file
	config_file, _ := os.Open("config.json")
	decoder := json.NewDecoder(config_file)
	config := Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Println("error:", err)
	}
	log.Println("listen address:", config.Listen)
	flag.Parse()
	err = http.ListenAndServe(config.Listen, http.HandlerFunc(mainHandler))
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	//只处理Post方法
	if req.Method != "POST" {
		log.Println(req.URL.String())
		newReq, err := http.NewRequest(req.Method, req.URL.String(), nil)
		newReq.Header = req.Header
		newReq.URL.Host = "api.mch.weixin.qq.com"
		newReq.Host = newReq.URL.Host
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
	}
}
