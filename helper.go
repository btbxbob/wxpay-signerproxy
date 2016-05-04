package main

import (
	"encoding/xml"
	"log"
)

type nameValue struct {
	XMLName xml.Name `xml:""`
	Value   string   `xml:",cdata"`
}

type xmlRootStruct struct {
	Elements []nameValue `xml:",any"`
}

// XMLStructToMap turn XML into struct
func XMLStructToMap(xmlRawData []byte) (result map[string]string, err error) {
	result = make(map[string]string)
	v := new(xmlRootStruct)
	err = xml.Unmarshal(xmlRawData, v)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	//log.Printf("v = %#v\n", v)
	//v = &main.xmlRootStruct{Elements:[]main.nameValue{main.nameValue{XMLName:xml.Name{Space:"", Local:"appid"}, Value:"wx2b029c08a6232582"}, main.nameValue{XMLName:xml.Name{Space:"", Local:"mch_id"}, Value:"1305176001"}, main.nameValue{XMLName:xml.Name{Space:"", Local:"nonce_str"}, Value:"ec2316275641faa3aacf3cc599e8730f"}, main.nameValue{XMLName:xml.Name{Space:"", Local:"transaction_id"}, Value:"4004792001201604285304611529"}, main.nameValue{XMLName:xml.Name{Space:"", Local:"sign"}, Value:"CCE53B58591F386DA7D0FEE640EE15CA"}}}
	for i := 0; i < len(v.Elements); i++ {
		result[v.Elements[i].XMLName.Local] = v.Elements[i].Value
	}
	return result, nil
}

// CalculateSignature return sign string
func CalculateSignature(fields map[string]string, key string) (result string) {
	return result
}
