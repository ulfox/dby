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
  * [Document Management](#document-management)
      + [Add a new doc](#add-a-new-doc)
      + [Switch Doc](#switch-doc)
      + [Document names](#document-names)
        - [Name documents manually](#name-documents-manually)
        - [Name all documents automatically](#name-all-documents-automatically)
        - [Switch between docs by name](#switch-between-docs-by-name)
      + [Import Docs](#import-docs)
      + [Global Commands](#global-commands)
        - [Global Upsert](#global-upsert)
        - [Global Update](#global-update)
        - [Global GetFirst](#global-getfirst)
        - [Global FindKeys](#global-findkeys)
        - [Global GetPath](#global-getpath)
        - [Global Delete](#global-delete)
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

### Initiate a new stateful DB

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

The code above will create a new yaml file under **local** directory.

### Initiate a new stateless DB

```go
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/ulfox/dby/db"
)

func main() {
	logger := logrus.New()

	state, err := db.NewStorageFactory()
	if err != nil {
		logger.Fatalf(err.Error())
	}
}
```

Initiating a db without arguments will not create/write/read from a file. All operations
will be done in memory and unless the caller saves the data externally, all data will be lose
on termination


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
keys, err := state.FindKeys("key-1")
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

```go
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


### Document Management

DBy creates by default an array of documents called library. That is in fact an array of interfaces

When initiating DBy, document 0 (index 0) is creatd by default and any action is done to that document, unless we switch to a new one

#### Add a new doc

To add a new doc, issue

```go
err = state.AddDoc()
if err != nil {
  logger.Fatal(err)
}

```

**Note: Adding a new doc also switches the pointer to that doc. Any action will write/read from the new doc by default**

#### Switch Doc

To switch a different document, we can use **Switch** method that takes as an argument an index

For example to switch to doc 1 (second doc), issue

```go
err = state.Switch(1)
if err != nil {
  logger.Fatal(err)
}
```

#### Document names

When we work with more than 1 document, we may want to set names in order to easily switch between docs

We have 2 ways to name our documents

- Add a name to each document manually
- Add a name providing a path that exists in all documents

##### Name documents manually

To name a document manually, we can use the **SetName** method which takes 2 arguments

- name
- doc index

For example to name document with index 0, as myDoc

```go
err := state.SetName("myDoc", 0)
if err != nil {
  logger.Fatal(err)
}
```

##### Name all documents automatically

To name all documents automatically we need to ensure that the same path exists in all documents.

The method for updating all documents is called **SetNames** and takes 2 arguments

- Prefix: A path in the documents that will be used for the first name
- Suffix: A path in the documents that will be used for the last name

**Note: Docs that do not have the paths that are queried will not get a name**

This method best works with **Kubernetes** manifests, where all docs have a common set of fields. 

For example

```yaml
apiVersion: someApi-0
kind: someKind-0
metadata:
...
  name: someName-0
...
---
apiVersion: someApi-1
kind: someKind-1
metadata:
...
  name: someName-1
...
---
```

From above we could give a name for all our documents if we use **kind** + **metadata.name** for the name.

```go
err := state.SetNames("kind", "metadata.name")
if err != nil {
  logger.Fatal(err)
}
```

###### List all doc names

To get the name of all named docs, issue

```go
for i, j := range state.ListDocs() {
  fmt.Println(i, j)
}
```
Example output based on the previous **SetNames** example

```bash
0 service/listener-svc
1 poddisruptionbudget/listener-svc
2 horizontalpodautoscaler/caller-svc
3 deployment/caller-svc
4 service/caller-svc
5 poddisruptionbudget/caller-svc
6 horizontalpodautoscaler/listener-svc
7 deployment/listener-svc
```

##### Switch between docs by name

To switch to a doc by using the doc's name, issue

```go
err = state.SwitchDoc("PodDisruptionBudget/caller-svc")
if err != nil {
  logger.Fatal(err)
}
```

#### Import Docs

We can import a set of docs with **ImportDocs** method. For example if we have the following yaml

```yaml
apiVersion: someApi-0
kind: someKind-0
metadata:
...
  name: someName-0
...
---
apiVersion: someApi-1
kind: someKind-1
metadata:
...
  name: someName-1
...
---
```

We can import it by giving the path of the file

```go
err = state.ImportDocs("file-name.yaml")
if err != nil {
  logger.Fatal(err)
}
```

#### Global Commands

Wrappers for working with all documents

##### Global Upsert

We can use upsert to update or create keys on all documents

```go
err = state.UpsertGlobal(
  "some.path",
  "v0.3.0",
)
if err != nil {
  logger.Fatal(err)
}

```

##### Global Update

Global update works as **GlobalUpsert** but it skips documents that
miss a path rather than creating the path on those docs.

##### Global GetFirst

To get the value of the first key in the hierarchy for each document, issue

```go
valueOfDocs, err := state.GetFirstGlobal("keyName")
if err != nil {
	logger.Fatalf(err.Error())
}
logger.Info(valueOfDocs)
```

This returns a `map[int]interface{}` object. The key is the index of each document and it's value
is the value of the first key in the hierarchy in that document

##### Global FindKeys

To get all the paths for a given from all documents, issue

```go
mapOfPaths, err := state.FindKeysGlobal("keyName")
if err != nil {
	logger.Fatalf(err.Error())
}
logger.Info(mapOfPaths)
```

This returns a `map[int][]string` object. The key is the index of each document and it's value
is a list of paths that have the queried key

##### Global GetPath

To get a path that exists in all documents, issue

```go
valueOfDocs, err := state.GetPathGlobal("key-1.key-2.key-3")
if err != nil {
	logger.Fatalf(err.Error())
}
logger.Info(valueOfDocs)
```

This returns a `map[int]interface{}` object. The key is the index of each document and it's value
is the value of the specific key in that document

##### Global Delete

To delete a path from all documents, issue

```go
err := state.DeleteGlobal("key-1.key-2.key-3")
if err != nil {
	logger.Fatalf(err.Error())
}
```

The above will delete all the paths that match the queried path from each doc

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
if assertData.GetError() != nil {
	logger.Fatal(assertData.GetError())
}
vMap, err := assertData.GetMap()
if err != nil {
	logger.Fatal(err)
}
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
if assertData.GetError() != nil {
	logger.Fatal(assertData.GetError())
}
vMap, err := assertData.GetMap()
if err != nil {
	logger.Fatal(err)
}
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
if assertData.GetError() != nil {
	logger.Fatal(assertData.GetError())
}
vArray, err := assertData.GetArray()
if err != nil {
	logger.Fatal(err)
}
logger.Info(vArray)

```

##### Get array manually

We can get the array manually by using only **Convert** operations

```go

assertData.Input(state.Data).
	Key("to").
	Key("array-2")
if assertData.GetError() != nil {
	logger.Fatal(assertData.GetError())
}
vArray, err := assertData.GetArray()
if err != nil {
	logger.Fatal(err)
}
logger.Info(vArray)

```
