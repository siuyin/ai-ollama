package main

import (
	"bufio"
	"fmt"
	"os"
)
func main() {
	fmt.Println("hello")

	sc:=bufio.NewScanner(os.Stdin)
	sc.Scan()
	fmt.Println(sc.Text())
}