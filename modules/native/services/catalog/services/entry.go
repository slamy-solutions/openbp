package services

import "fmt"

func MakeCatalogEntryCollectionName(namespace string, name string) string {
	return fmt.Sprintf("native_catalog_%s_%s", namespace, name)
}
