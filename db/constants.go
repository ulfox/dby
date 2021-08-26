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
