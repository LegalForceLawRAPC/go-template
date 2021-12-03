package utils

import (
	"log"
	"runtime/debug"
)

func Recover() {
	if r := recover(); r != nil {
		log.Println("Recovering from panic:", string(debug.Stack()))
	}
}

// defer utils.Recover()
