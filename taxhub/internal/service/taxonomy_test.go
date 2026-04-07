package service

import (
	"os"
	"testing"
)

func writeTempCSV(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "taxonomy*.csv")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadCSV_Get(t *testing.T) {
	tx := NewTaxonomy()
	path := writeTempCSV(t, "tax_id,name\n9606,Homo sapiens\n")

	if err := tx.LoadCSV(path); err != nil {
		t.Fatal(err)
	}

	name, ok := tx.Get("9606")
	if !ok || name != "Homo sapiens" {
		t.Errorf("expected Homo sapiens, got %q ok=%v", name, ok)
	}
}

func TestGet_notFound(t *testing.T) {
	tx := NewTaxonomy()
	_, ok := tx.Get("0")
	if ok {
		t.Error("expected not found for empty taxonomy")
	}
}

func TestCount(t *testing.T) {
	tx := NewTaxonomy()
	path := writeTempCSV(t, "tax_id,name\n1,Bacteria\n2,Archaea\n")

	if err := tx.LoadCSV(path); err != nil {
		t.Fatal(err)
	}

	if tx.Count() != 2 {
		t.Errorf("expected 2, got %d", tx.Count())
	}
}

func TestLoadCSV_invalidHeader(t *testing.T) {
	tx := NewTaxonomy()
	path := writeTempCSV(t, "only_one_column\n9606\n")

	if err := tx.LoadCSV(path); err == nil {
		t.Error("expected error for invalid header")
	}
}
