## DB Yaml

Simple DB using yaml. A project for managing the content of yaml files.

Table of Contents
=================
- [DB Yaml](#db-yaml)
- [Features](#features)
- [Usage](#usage)
    * [Write to DB](#write-to-db)
    * [Query DB](#query-db)
      + [Get First Key](#get-first-key)
      + [Search for Keys](#search-for-keys)
    * [Query Path](#query-path)
      + [Query Path with Arrays](#query-path-with-arrays)
        - [Without trailing array](#without-trailing-array)
        - [With trailing array](#with-trailing-array)
    * [Delete Key By Path](#delete-key-by-path)
	* [Convert Utils](#convert-utils)
      + [Get map of strings from interface](#get-map-of-strings-from-interface)
        - [Get map directly from a GetPath object](#get-map-directly-from-a-getpath-object)
		- [Get map manually](#get-map-manually)
	  + [Get array of string from interface](#get-array-of-string-from-interface)
        - [Get array directly from a GetPath object](#get-array-directly-from-a-getpath-object)
		- [Get array manually](#get-array-manually)

## Features

The module can do

- Create/Load yaml files
- Update content
- Get values from keys
- Query for keys
- Delete keys
- Merge content

##  Usage

Simple examples for working with yaml files as db

### Initiate a new DB

Create a new local DB

```go
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/ulfox/dby/db"
)

func main() {
	logger := logrus.New()

	state, err := db.NewStorageFactory("local/db.yaml")
	if err != nil {
		logger.Fatalf(err.Error())
	}
}
```

### Write to DB

Insert a map to the local yaml file.

```go
err = state.Upsert(
	"some.path",
	map[string]string{
		"key-1": "value-1",
		"key-2": "value-2",
	},
)

if err != nil {
	logger.Fatalf(err.Error())
}
```

### Query DB

#### Get First Key

Get the value of the first key in the hierarchy (if any)

```go
val, err := state.GetFirst("key-1")
if err != nil {
	logger.Fatalf(err.Error())
}
logger.Info(val)
```

For example if we have the following structure

```yaml
key-1:
    key-2:
        key-3: "1"
    key-3: "2"
```

And we query for `key-3`, then we will get back **"2"** and not **"1"**
since `key-3` appears first on a higher layer with a value of **2**

#### Search for keys

Get all they keys (if any). This returns the full path for the key,
not the key values. To get the values check the next section **GetPath**

```go
keys, err := state.Get("key-1")
if err != nil {
	logger.Fatalf(err.Error())
}
logger.Info(keys)
```

From the previous example, this query would have returned

```yaml
["key-1.key-2.key-3", "key-1.key-3"]
```

### Query Path

Get the value from a given path (if any)

For example if we have in yaml file the following key-path

```yaml
key-1:
    key-2:
        key-3: someValue
```

Then to get someValue, issue

```go
keyPath, err := state.GetPath("key-1.key-2.key-3")
if err != nil {
	logger.Fatalf(err.Error())
}
logger.Info(keyPath)
```

#### Query Path with Arrays

We can also query paths that have arrays.

##### Without trailing array

```yaml
key-1:
    key-2:
        - key-3: 
            key-4: value-1
```

To get the value of `key-4`, issue

```go
keyPath, err := state.GetPath("key-1.key-2.[0].key-3.key-4")
if err != nil {
	logger.Fatalf(err.Error())
}
logger.Info(keyPath)
```

##### With trailing array

```yaml
key-1:
    key-2:
        - value-1
        - value-2
        - value-3
```

To get the first index of `key-2`, issue

```
keyPath, err := state.GetPath("key-1.key-2.[0]")
if err != nil {
	logger.Fatalf(err.Error())
}
logger.Info(keyPath)
```

### Delete Key By Path

To delete a single key for a given path, e.g. key-2
from the example above, issue

```go
err = state.Delete("key-1.key-2")
if err != nil {
	logger.Fatalf(err.Error())
}
```

### Convert Utils

Convert simply automate the need to
explicitly do assertion each time we need to access
an interface object.

Let us assume we have the following YAML structure

```yaml

to:
  array-1:
    key-1:
    - key-2: 2
    - key-3: 3
    - key-4: 4
  array-2:
  - 1
  - 2
  - 3
  - 4
  - 5
  array-3:
  - key-1: 1
  - key-2: 2

```

#### Get map of strings from interface

We can do this in two ways, get object by giving a path and assert the interface to `map[string]string`, or work manually our way to the object

##### Get map directly from a GetPath object

To get map **key-2: 2**, first get object via GetPath

```go

obj, err := state.GetPath("to.array-1.key-1.[0]")
if err != nil {
	logger.Fatalf(err.Error())
}
logger.Info(val)

```

Next, assert **obj** as `map[string]string`

```go

assertData := db.NewConvertFactory()

assertData.Input(val)
if assertData.Error != nil {
	logger.Fatalf(assertData.Error.Error())
}
vMap := assertData.GetMap()
logger.Info(vMap["key-2"])

```

##### Get map manually

We can get the map manually by using only **Convert** operations

```go

assertData := db.NewConvertFactory()

assertData.Input(state.Data).
	Key("to").
	Key("array-1").
	Key("key-1").Index(0)
if assertData.Error != nil {
	logger.Fatalf(assertData.Error.Error())
}
vMap := assertData.GetMap()
logger.Info(vMap["key-2"])

```

#### Get array of string from interface

Again here we can do it two ways as with the map example

##### Get array directly from a GetPath object

To get **array-2** as **[]string**, first get object via GetPath

```go

obj, err = state.GetPath("to.array-2")
if err != nil {
	logger.Fatalf(err.Error())
}
logger.Info(obj)

```

Next, assert **obj** as `[]string`

```go

assertData := db.NewConvertFactory()

assertData.Input(obj)
if assertData.Error != nil {
	logger.Fatalf(assertData.Error.Error())
}
vArray := assertData.GetArray()
logger.Info(vArray)

```

##### Get array manually

We can get the array manually by using only **Convert** operations

```go

assertData.Input(state.Data).
	Key("to").
	Key("array-2")
if assertData.Error != nil {
	logger.Fatalf(assertData.Error.Error())
}
vArray := assertData.GetArray()
logger.Info(vArray)

```
