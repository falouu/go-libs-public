package main

import "fmt"

func main() {

	fmt.Println("Printing first line immiedatelly. Waiting for confirmation...")
	fmt.Scanln()
	fmt.Print("Printing half of the second line and waiting... ")
	fmt.Scanln()
	fmt.Println("then printing the rest.")
}
