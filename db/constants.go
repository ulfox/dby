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

func getObjectType(o interface{}) objectType {
	_, isMap := o.(map[interface{}]interface{})
	if isMap {
		return 1
	}

	_, isArray := o.([]interface{})
	if isArray {
		return 2
	}

	return unknownObj
}
