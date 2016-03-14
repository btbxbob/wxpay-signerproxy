package main

import (
	"flag"
	//"fmt"
	"io"
	//"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("127.0.0.1", ":1718", "http service address") // Q=17, R=18

func main() {
	flag.Parse()
	//http.Handle("/", http.HandlerFunc(QR))
	err := http.ListenAndServe(*addr, http.HandlerFunc(mainHandler))
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		//io.WriteString(w, "please use post method.<br />")
		//io.WriteString(w, req.URL.Path)
		log.Println(req.URL.String())
		newReq, err := http.NewRequest(req.Method, req.URL.String(), nil)
		newReq.Header = req.Header
		newReq.URL.Host = "api.mch.weixin.qq.com"
		newReq.Host = newReq.URL.Host
		newReq.URL.Scheme = "https"
		//newReq.URL.Path = req.URL.Path

		client := &http.Client{}
		resp, err := client.Do(newReq)

		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println(resp.StatusCode)

		defer resp.Body.Close()

		//body, err := ioutil.ReadAll(resp.Body)
		//log.Println(string(body))

		//resp.Write(w)
		io.Copy(w, resp.Body)

		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}
}

/*
func QR(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		io.WriteString(w, "please use post method.")
		//templ.Execute(w, req.FormValue("s"))
		return
	}
}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET"><input maxLength=1024 size=70
name=s value="" title="Text to QR Encode"><input type=submit
value="Show QR" name=qr>
</form>
</body>
</html>
`
*/
