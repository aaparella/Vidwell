package users

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/storage"
	"github.com/gorilla/mux"
)

const USER_HTML = `
	<html>
		<div class="user">
			{{ .DisplayName }} {{ .Email }}
		</div>
	</html>
`

func ViewUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User

	if err := storage.DB.Find(&user, id); err != nil {
		fmt.Fprintf(w, "Could not find user with ID: %s", id)
		return
	}

	tmpl, _ := template.New("user").Parse(USER_HTML)
	tmpl.Execute(w, user)
}
