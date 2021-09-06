### Update labels for all kubernetes manifests

In the **manifests directory** we have a **deployment.yaml** that we will import and update

In this example we will import the manifest and update the version for all documents from **v0.2.0** to **v0.3.0**

```go

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ulfox/dby/db"
)

func main() {
	logger := logrus.New()
	state, err := db.NewStorageFactory("local/db.yaml")
	if err != nil {
		logger.Fatal(err)
	}
	err = state.ImportDocs("docs/examples/manifests/deployment.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	// Automatically update all document names based on "kind/metadata.name" values
	state.SetNames("kind", "metadata.name")

	// Set the paths we want to update
	paths := []string{
		"spec.selector.matchLabels.version",
		"metadata.labels.version",
		"spec.selector.version",
		"spec.template.selector.matchLabels.version",
		"spec.template.metadata.labels.version",
	}

	// UpdateGlobal is a global command that updates all fields
	// that match the given path. Documents that do not have the
	// specific path will not be updated
	//
	// If we wanted to update or create the path then we could issue
	// UpsertGlobal() instead. Using that command however for Kubernetes
	// manifests is not recommended since you may end up having
	// manifests with fields that are not supported by the resource API
	for _, j := range paths {
		err = state.UpdateGlobal(
			j,
			"v0.3.0",
		)
		if err != nil {
			logger.Fatal(err)
		}

	}

	// List Docs by name
	for _, j := range state.ListDocs() {
		// Switch to a doc by name
		err = state.SwitchDoc(j)
		if err != nil {
			logger.Fatal(err)
		}

		// Get the metadata
		val, err := state.GetPath("metadata.labels.version")
		if err != nil {
			// We use continue here because HorizontalPodAutoscaler does not have labels set
			// so no update was done and no path exists for us to get
			continue
		}
		logger.Infof("%s has version: %s", j, val)
	}
}

```


Example output


```bash
INFO[0000] poddisruptionbudget/caller-svc has version: v0.3.0 
INFO[0000] deployment/listener-svc has version: v0.3.0  
INFO[0000] service/listener-svc has version: v0.3.0     
INFO[0000] poddisruptionbudget/listener-svc has version: v0.3.0 
INFO[0000] deployment/caller-svc has version: v0.3.0    
INFO[0000] service/caller-svc has version: v0.3.0 
```