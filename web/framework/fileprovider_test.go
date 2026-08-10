package framework

import (
	"io"
	"testing"
	"testing/fstest"
)

func TestFileProviderListsReadsAndOpensFiles(t *testing.T) {
	provider := NewFileProvider(fstest.MapFS{
		"templates/page.gohtml": {Data: []byte("page template")},
		"css/site.css":          {Data: []byte("body {}")},
	})

	wantNames := []string{"css/site.css", "templates/page.gohtml"}
	names, err := provider.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(names) != len(wantNames) {
		t.Fatalf("List() = %v, want %v", names, wantNames)
	}
	for i := range wantNames {
		if names[i] != wantNames[i] {
			t.Errorf("List()[%d] = %q, want %q", i, names[i], wantNames[i])
		}
	}

	template, err := provider.ReadFile("templates/page.gohtml")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if got := string(template); got != "page template" {
		t.Errorf("ReadFile() = %q, want %q", got, "page template")
	}
	if _, err := provider.ReadFile("missing"); err == nil {
		t.Error("ReadFile(missing) error = nil, want an error")
	}

	file, err := provider.Open("css/site.css")
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if got := string(contents); got != "body {}" {
		t.Errorf("opened contents = %q, want %q", got, "body {}")
	}
}
