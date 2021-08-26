package db

import "fmt"

// Inormational Error constants. Used during a return ... errors.New()
const (
	notAMap         = "target object is not a map"
	notArrayObj     = "received a non array object but expected []interface{}"
	keyDoesNotExist = "the given key [%s] does not exist"
	fileNotExist    = "the given file [%s] does not exist"
	dictNotFile     = "can not create file [%s], a directory exists with that name"
	notAnIndex      = "object (%s) is not an index. Index example: some.path.[someInteger].someKey"
	arrayOutOfRange = "index value (%s) is bigger than the length (%s) of the array to be indexed"
	invalidKeyPath  = "the key||path [%s] that was given is not valid"
)

func errString(b, k string) string {
	return fmt.Sprintf(b, k)
}