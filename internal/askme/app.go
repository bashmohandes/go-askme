package askme

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bashmohandes/go-askme/internal/askme/controllers"
)

//Start method starts the AskMe App
func Start() {
	b := controllers.Blog()
	fmt.Println("Listening on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", b))
}
