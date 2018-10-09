package providers

import (
	"io/ioutil"
	"log"

	"github.com/bashmohandes/go-askme/shared"
	"github.com/bashmohandes/go-askme/web/askme"
	"github.com/fsnotify/fsnotify"
	"github.com/gobuffalo/packr"
)

type packrFileProvider struct {
	packr.Box
	watcher *fsnotify.Watcher
}

// NewFileProvider creates new packr based file provider
func NewFileProvider(config *askme.Config) shared.FileProvider {
	result := &packrFileProvider{
		Box: packr.NewBox(config.Assets),
	}

	return result
}

func (p *packrFileProvider) Watch() {
	if p.watcher != nil {
		panic("Watcher is already configured for this instance")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	p.watcher = watcher

	p.Walk(p.watchDir)

	go func(fp *packrFileProvider) {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				if event.Op == fsnotify.Write {
					log.Printf("File updated! %v\n", event.Name)
					content, err := ioutil.ReadFile(event.Name)
					if err != nil {
						log.Fatal(err)
					}
					fp.AddBytes(event.Name, content)
				}

				// watch for errors
			case err := <-watcher.Errors:
				log.Println("ERROR", err)
			}
		}
	}(p)
}

func (p *packrFileProvider) Close() {
	p.watcher.Close()
}

func (p *packrFileProvider) watchDir(path string, fi packr.File) error {
	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	return p.watcher.Add(path)
}
