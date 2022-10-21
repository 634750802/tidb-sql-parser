package analyze

import "testing"

func AssertColumnEquals(t *testing.T, column *Column, to *Column) {
	if column.Name != to.Name {
		t.Errorf("expect column.Name to be %s, actual: %s", to.Name, column.Name)
	}
	if column.Type != to.Type {
		t.Errorf("expect column.Type to be %d, actual: %d", to.Type, column.Type)
	}
	if column.Nullable != to.Nullable {
		t.Errorf("expect column.Nullable to be %t, actual: %t", to.Nullable, column.Nullable)
	}
}
