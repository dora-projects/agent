package internal

import (
	"testing"
)

func TestGetSubHost(t *testing.T) {
	host := "1312312.123.web.nancode.cn"
	str := GetSubHost(host, "web.nancode.cn")

	expect := "1312312.123"

	if str != expect {
		t.Errorf("expect: %v, but got: %v", expect, str)
	}
}

func TestParseConfig(t *testing.T) {
	testData := `{
    "release": "",
    "filepath": "./build",
    "index": "index.html",
    "proxy": {
      "/api": "http://localhost:3000"
    }
  }
`

	config, err := ParseSiteConfig([]byte(testData))

	if err != nil {
		t.Logf("%#v", config)
		t.Fatal(err)
	}

	//t.Logf("%#v", config)
}
