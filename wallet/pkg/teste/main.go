package main

import "fmt"

func main() {
	evento := []string{"teste1", "teste2", "teste3", "teste4"}
	evento = append(evento[:0], evento[1:]...)
	// 0 1 2 3 4
	fmt.Println(evento)
}
