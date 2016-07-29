package main

import (
	"fmt"
	"os"

	"github.com/sg3des/carrot"
)

func init() {
	os.Chdir("../")
}

func main() {
	if err := carrot.Open("db"); err != nil {
		panic(err)
	}

	u0 := &carrot.Users{Name: "bunny", Number: 324232}
	fmt.Println("write:", u0)

	u0.Write()
	fmt.Println("writed item resieved id:", u0.ID)

	var u1 = new(carrot.Users)
	if err := u1.Read(u0.ID); err != nil {
		panic(err)
	}
	fmt.Println("read:", u1)

	carrot.Close()
}
