## ToDo

### Features

- Add multiple keys
- Remote backends: Allow to work with yaml files that are on a remote location
  - S3
  - Google Storage
- Export to STDOUT method

### Improvements

### Remove nested loops/ifs for better readability

Description: Some blocks have nested loops/ifs (e.g. for{if{if{}}}) which are not good
either for maintenace or readability.

### Consolidate repeating code blocks

Description: There are code blocks that appear on different methods. While they are
not exactly identical, their structure is. We could consolidate these blocks to a
more abstract functions in order to reduce repitivness and work needed for future
changes.

#### Introduce double linked lists for path seeking

Description: When we want to find valid paths for a given key, or we want to get 
the content of a path, we use getPath() method which is a loop that calls itself 
until it finds a path that leads to the desired key.

Linked lists would be a better solution for storing our path since we can easily
know our full path and also have easier access to the content for any given branch.

#### Optimize more the cache

Description: Methods like getPath(), get(), deletePath() can be further optimized by 
using cache to avoid multiple copies and loops. 

#### Reduce write operations

Description: Currently we do a write on Upsert() and Delete() methods. We could instead 
work in ram and write periodically in the persistant storage.

#### MergeDBs should not shadow target

Description: Currently DBy supports yaml merges. That is we can merge a source file with 
the content of the DBy yaml file. The issue we currently have is that the merge replaces 
same paths with the path from the source file rather than merging both paths.

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

### Tests

- Cover cache state during operations
- Cover different upsert structures.
    - []string
    - []map[string]string
    - []interface{}
    - []map[string][]string
    - []map[string][]interface{}


