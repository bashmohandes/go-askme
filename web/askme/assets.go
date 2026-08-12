package askme

import (
	"embed"
	"io/fs"

	"github.com/bashmohandes/go-askme/web/framework"
)

//go:embed public
var publicFiles embed.FS

// NewPublicFiles returns the application's embedded templates and static assets.
func NewPublicFiles() (framework.FileProvider, error) {
	root, err := fs.Sub(publicFiles, "public")
	if err != nil {
		return nil, err
	}
	return framework.NewFileProvider(root), nil
}
