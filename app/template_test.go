package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

func BenchmarkTemplateCach(b *testing.B) {
	tms, err := newTemplateCache(TEMPLATES_PATH)
	fmt.Printf("Templates: %#v\nError: %s\n", tms, err)
	for _, tm := range tms {
		fmt.Printf("template's name: %s\nBase: %s\n", tm.Name(), filepath.Base(tm.Name()))
		fmt.Printf("tm: %v\n", tm)
		fmt.Printf("tm's tree: %v\n", tm.Tree)
		fmt.Printf("tm.Tree.Name: %v\n", tm.Tree.Name)

	}
}
