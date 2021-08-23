package db

// Inormational Error constants. Used during a return ... errors.New()
const (
	keyExists            = "the given key already exists. Use update or upsert instead"
	notAMap              = "object to be updated is not a map"
	notArrayObj          = "received a non array object but expected []interface{}"
	keyDoesNotExist      = "the given key does not exist"
	fileNotExist         = "the given file does not exist"
	deleteObjKeyNotFound = "could not delete an existing key. Possible delete function bug"
	dictNotFile          = "a directory exists with that name"
	notAnIndex           = "object is not an index. Index example: some.path.[someInteger].someKey"
	arrayOutOfRange      = "index value was bigger than the array to be indexed"
)

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
