package db

type objectType int

const (
	unknownObj objectType = iota
	mapObj
	arrayObj
	arrayMapObj
	mapStringString
	mapStringInterface
	mapStringArrayString
	arrayMapStringArrayString
	arrayMapStringString
	arrayMapStringInterface
	arrayMapStringArrayInterface
)

// Inormational Error constants. Used during a return err
const (
	notAMap         = "target object is not a map"
	notArrayObj     = "received a non array object but expected []interface{}"
	keyDoesNotExist = "the given key [%s] does not exist"
	fileNotExist    = "the given file [%s] does not exist"
	dictNotFile     = "can not create file [%s], a directory exists with that name"
	notAnIndex      = "object (%s) is not an index. Index example: some.path.[someInteger].someKey"
	arrayOutOfRange = "index value (%s) is bigger than the length (%s) of the array to be indexed"
	invalidKeyPath  = "the key||path [%s] that was given is not valid"
	emptyKey        = "path [%s] contains an empty key"
	libOutOfIndex   = "lib out of index"
	docNotExists    = "doc [%s] does not exist in lib"
	fieldNotString  = "[%s] with value [%s] is not a string"
	notAType        = "value is not a %s"
)

func getObjectType(o interface{}) objectType {
	_, isMap := o.(map[interface{}]interface{})
	if isMap {
		return 1
	}

	_, isArray := o.([]interface{})
	if isArray {
		return 2
	}

	_, isMapStringInterface := o.(map[string]interface{})
	if isMapStringInterface {
		return 5
	}

	return unknownObj
}
