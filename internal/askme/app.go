package askme

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/packr"

	"github.com/bashmohandes/go-askme/internal/askme/controllers"
)

//Start method starts the AskMe App
func Start() {
	box := packr.NewBox("./templates")
	b := controllers.Blog(box)
	fmt.Println("Listening on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", b))
}
