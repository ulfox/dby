## ToDo

### Features

- Add multiple keys

### Improvements

#### Optimize more the cache

Description: Methods like getPath(), get(), deletePath() can be further optimized by using cache to avoid multiple copies and loops. 

#### Reduce write operations

Description: Currently we do a write on Upsert() and Delete() methods. We could instead work in ram and write periodically in the persistant
storage.

#### MergeDBs should not shadow target

Description: Currently DBy supports yaml merges. That is we can merge a source file with the content of the DBy yaml file. The issue we currently
have is that the merge replaces same paths with the path from the source file rather than merging both paths.

##### Example

DBy yaml content

```
key-1:
  key-2:
    key-3: value-3

some:
  other:
    path: test
```

Some local yaml file

```
key-1:
  key-2:
    key-4: value-4

```

Current Merged DBy

```
key-1:
  key-2:
    key-4: value-4

some:
  other:
    path: test
```

Expected Merge

```
key-1:
  key-2:
    key-4: value-4
    key-3: value-3

some:
  other:
    path: test
```

#### Better errors

Description: Provide more details on errors and have better wrapping during error propagation

### Tests

- Cover cache state during operations
- Cover different upsert structures.
    - []string
    - []map[string]string
    - []interface{}
    - []map[string][]string
    - []map[string][]interface{}


