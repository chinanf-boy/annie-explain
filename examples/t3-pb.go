package examples

import (
	"github.com/cheggaaa/pb"
	"os"
	"io"
)

func Pbtry(){
	var file *os.File
	myDataLen := 100
	bar := pb.New(myDataLen).SetUnits(pb.U_BYTES)
	bar.Start()
	
	// my io.Reader
	r, _ := os.Open("./readme.md")
	
	// my io.Writer
	file, _ = os.Create("./try-pb.md")
	
	// create multi writer
	writer := io.MultiWriter(file, bar)
	
	// and copy
	io.Copy(writer, r)
	
	bar.Finish()
}