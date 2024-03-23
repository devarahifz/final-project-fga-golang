package helpers

import (
	"fmt"
	"os"
)

func GetEnv(key string) string {
	var value, ok = os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s env has not been set", key))
	}

	return value
}
