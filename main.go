package main

import (
	"github.com/chinanf-boy/annie-explain/examples"
	"flag"
)

func init(){
}
func main(){
	flag.Parse()
	exampleIndex := flag.Args()
	if(!(len(exampleIndex) > 0)){
		return 
	}
	if(exampleIndex[0] == "1"){
		examples.Gethello("local")
	}
	if(exampleIndex[0] == "pb"){
		examples.Pbtry()
	}	
}