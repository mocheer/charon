// +build wasm
package main

import (
	"fmt"
	"syscall/js"
)

func foo(args []js.Value) {
	fmt.Println("hellow wasm")
	fmt.Println(args)
}

func main() {
	js.Global().Set("foo", foo)
	js.Global().Set("value", 100)
	select {}
}
