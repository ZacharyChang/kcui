package log

import (
	"log"
	"os"
)

func init() {
	f, err := os.OpenFile("kcui.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//defer f.Close()

	log.SetOutput(f)
}
