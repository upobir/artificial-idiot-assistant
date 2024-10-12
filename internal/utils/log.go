package utils

import (
	"fmt"
	"log"
	"time"
)

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("[%s] %s\n", time.Now().Format("2006/01/02 15:04:05"), string(bytes))
}

func InitializeLog() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}
