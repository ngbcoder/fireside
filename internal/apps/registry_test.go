package apps

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRegistry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "apps.json")
	data := `{"apps":[{"id":"yt","name":"YouTube","kind":0,"url":"https://youtube.com/tv"}]}`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	r, err := LoadRegistry(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Apps) != 1 {
		t.Fatalf("expected 1 app, got %d", len(r.Apps))
	}
	a, ok := r.Find("yt")
	if !ok || a.Name != "YouTube" {
		t.Fatalf("find failed: %+v", a)
	}
}
