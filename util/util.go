package util

import (
	"github.com/pmoule/go2hal/hal"
	"log"
)

// JSONify the resource
func JSONify(root hal.Resource) []byte {

	encoder := new(hal.Encoder)
	bytes, err := encoder.ToJSON(root)

	if err != nil {
		log.Println(err)
		return nil
	}
	return bytes
}
