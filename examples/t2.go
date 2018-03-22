package examples

import "fmt"

// Gethello get t1.go hello
func Gethello(s string){
	fmt.Println(" package examples 1")
	fmt.Printf(hello("func hello")+` in t1.go , t2.go can use func hello in Gethello `)
}