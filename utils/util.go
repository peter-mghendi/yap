package utils

import(
	"log"
)

// Try handles top-level errors
func Try(e error) {
	if e != nil {
		log.Println(e)
	}
}