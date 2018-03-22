package examples

import (
	"encoding/json"
	"fmt"
)

	
type jsonbyt struct {
    NUM  float64 `json:"num"`
    Star []string `json:"strs"`
}

// Jsontry show json
func Jsontry(){
	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
	var dat jsonbyt
	if err := json.Unmarshal(byt, &dat); err != nil {
        panic(err)
	}
	fmt.Println(dat)
}