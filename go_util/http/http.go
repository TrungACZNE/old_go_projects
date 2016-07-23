package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

var QueryStringReplacer *strings.Replacer

func init() {
	QueryStringReplacer = strings.NewReplacer(
		" ", "%20",
		"&", "%26",
		"'", "%27",
		"\"", "%22",
		"?", "%3F",
		"=", "%3D",
		"+", "%2B",
		"%", "%25",
	)
}

func escape(s string) string {
	return QueryStringReplacer.Replace(s)
}

func QueryString(dict map[string]string) string {
	qs := ""
	c := 0
	for k, v := range dict {
		if c == 0 {
			qs += "?"
		}

		qs += k + "=" + escape(v)

		if c != len(dict)-1 {
			qs += "&"
		}
		c += 1
	}
	return qs
}

func GetAndUnmarshal(url string, i interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	} else {
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		} else {
			err := json.Unmarshal(bytes, i)
			if err != nil {
				return err
			}
			return nil
		}
	}
}
