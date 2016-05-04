package main

import (
	"reflect"
	"testing"
)

func TestSign(t *testing.T) {
	cases := []struct {
		in   []byte
		want string
	}{
		{[]byte(`<xml>
   <appid>wx2b029c08a6232582</appid>
   <mch_id>1305176001</mch_id>
   <nonce_str>ec2316275641faa3aacf3cc599e8730f</nonce_str>
   <transaction_id>4004792001201604285304611529</transaction_id>
   <sign>7FE21BD5F4EB3E082E31FE0D88D2DAF4</sign>
</xml>`), "7FE21BD5F4EB3E082E31FE0D88D2DAF4"}}
	for _, c := range cases {
		a, _ := XMLStructToMap(c.in)
		got, _ := CalculateSignature(a, "111")
		if got != c.want {
			t.Errorf("in %q, out %q, want %q", c.in, got, c.want)
		}
	}
}

func TestXML(t *testing.T) {

	xmlByte := []byte(`<xml>
 <appid>wx2b029c08a6232582</appid>
 <mch_id>1305176001</mch_id>
 <nonce_str>ec2316275641faa3aacf3cc599e8730f</nonce_str>
 <transaction_id>4004792001201604285304611529</transaction_id>
 <sign>7FE21BD5F4EB3E082E31FE0D88D2DAF4</sign>
</xml>`)

	mapData, _ := XMLStructToMap(xmlByte)
	xmlByte2, _ := MapToXML(mapData)
	mapData2, _ := XMLStructToMap(xmlByte2)
	eq := reflect.DeepEqual(mapData, mapData2)
	if !eq {
		t.Errorf("xml convert not the same, xml:%q", string(xmlByte2))
	}
}
