package resolve

import (
	"fmt"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func TestConfigurerResolve(t *testing.T) {
	configurer := &Configurer{}
	cfgRoot := config.New()
	configurer.RegisterFlags(nil, "", cfgRoot)
	// This initializes the overrides array with capacity > length to simulate
	// the case where overrides has been written to (e.g. when there is a parent BAZEL file that
	// specifies overrides) and Go has expanded the underlying array (this generally causes
	// cap(overrides) > len(overrides))
	initialOverrides := make([]overrideSpec, 0, 1)
	getResolveConfig(cfgRoot).overrides = initialOverrides

	cfg1 := cfgRoot.Clone()
	f := &rule.File{
		Directives: []rule.Directive{{
			Key:   "resolve",
			Value: fmt.Sprintf("go go github.com/golang/protobuf/ptypes1 //:good1"),
		}},
	}
	configurer.Configure(cfg1, "test", f)
	resolvedOverrides1 := getResolveConfig(cfg1).overrides

	// simulate walking a different subdirectory
	cfg2 := cfgRoot.Clone()
	f2 := &rule.File{
		Directives: []rule.Directive{{
			Key:   "resolve",
			Value: fmt.Sprintf("go go github.com/golang/protobuf/ptypes2 //:good2"),
		}},
	}
	configurer.Configure(cfg2, "test", f2)
	resolvedOverrides2 := getResolveConfig(cfg2).overrides

	if len(resolvedOverrides1) != 1 {
		t.Errorf("unexpected content for resolvedOverrides1: %v", resolvedOverrides1)
	}
	if resolvedOverrides1[0].imp.Imp != "github.com/golang/protobuf/ptypes1" {
		t.Errorf("expected %s but got %s", "github.com/golang/protobuf/ptypes1", resolvedOverrides1[0].imp.Imp)
	}
	if len(resolvedOverrides2) != 1 {
		t.Errorf("unexpected content for resolvedOverrides2: %v", resolvedOverrides2)
	}
	if resolvedOverrides2[0].imp.Imp != "github.com/golang/protobuf/ptypes2" {
		t.Errorf("expected %s but got %s", "github.com/golang/protobuf/ptypes2", resolvedOverrides2[0].imp.Imp)
	}
}
