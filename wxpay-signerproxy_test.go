package main

import "testing"

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
