package main

import (
	"fmt"
)


func main() {
	//generate random jwt secret
	secret := make([]byte, 32)
	for i := range secret {
		secret[i] = byte(65 + i) // Just a simple pattern for demonstration
	}
	fmt.Printf("Generated JWT secret: %s\n", secret)

	fmt.Println("Hello, World!")
}

