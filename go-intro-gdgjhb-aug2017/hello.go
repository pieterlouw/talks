package main

import (
	"fmt"
	"time"
)

func main() {
	var 界 string

	fmt.Println("Hello 世界")

	//You can even use unicode for variable names
	界 = "Wow!"

	fmt.Printf("%s\n", 界)

	time.Sleep(3 * time.Second)
}
