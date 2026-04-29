package commands

import "testing"

func TestDefaultCatalogFindsNativeAndLegacyCommands(t *testing.T) {
	catalog := DefaultCatalog()
	if len(catalog.Commands) < 100 {
		t.Fatalf("expected broad legacy command coverage, got %d commands", len(catalog.Commands))
	}
	cmd, ok := catalog.Find("convert:c-to-f")
	if !ok {
		t.Fatal("convert:c-to-f missing")
	}
	if !cmd.Native {
		t.Fatal("convert:c-to-f should be a native example module")
	}
	if _, ok := catalog.Find("system:kill"); !ok {
		t.Fatal("system:kill missing")
	}
}

func TestCategoriesAreStable(t *testing.T) {
	categories := DefaultCatalog().Categories()
	if len(categories) == 0 {
		t.Fatal("expected categories")
	}
	for i := 1; i < len(categories); i++ {
		if categories[i-1] > categories[i] {
			t.Fatalf("categories should be sorted: %q before %q", categories[i-1], categories[i])
		}
	}
}
