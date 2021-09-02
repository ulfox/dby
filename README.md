# DB Yaml

Simple DB using yaml. A project that for managing the 
content of yaml files.

Features:
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

#### GetSingle

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
