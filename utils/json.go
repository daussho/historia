package utils

import "encoding/json"

func JsonStringify(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
