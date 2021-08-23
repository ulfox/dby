package db

const (
	KeyExists            = "the given key already exists. Use update or upsert instead"
	NotAMap              = "object to be updated is not a map"
	KeyDoesNotExist      = "the given key does not exist"
	FileNotExist         = "the given file does not exist"
	DeleteObjKeyNotFound = "could not delete an existing key. Possible delete function bug"
	DictNotFile          = "a directory exists with that name"
)

type ObjectType int

const (
	UnknownObj ObjectType = iota
	MapObj
	ArrayObj
	ArrayMapObj
	MapStringString
	MapStringInterface
	MapStringArrayString
	ArrayMapStringArrayString
	ArrayMapStringString
	ArrayMapStringInterface
	ArrayMapStringArrayInterface
)
