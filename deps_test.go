package allsrvc_test

import (
	"os"
	"testing"

	"golang.org/x/mod/modfile"
)

func TestDependencies(t *testing.T) {
	b, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatal(err)
	}

	mf, err := modfile.Parse("go.mod", b, nil)
	if err != nil {
		t.Fatal(err)
	}

	allowed := map[string]bool{
		"github.com/jsteenb2/errors": true,
		"golang.org/x/mod":           true,
	}
	for range mf.Require {

	}

	var forbiddenAdded []string
	for _, mod := range mf.Require {
		if mod == nil {
			continue
		}
		if p := mod.Mod.Path; !allowed[p] {
			forbiddenAdded = append(forbiddenAdded, p)
		}
	}
	if len(forbiddenAdded) > 0 {
		t.Fatalf("forbidden modules have been added to the pkg, please remove:\n%v", forbiddenAdded)
	}
}
