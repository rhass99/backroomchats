package backroom

import (
	"google.golang.org/appengine"
	"net/http"
	"text/template"
)

const clientIdLen = 40

func init() {
	// Register our handlers with the http package.
	http.HandleFunc("/", root)
	http.HandleFunc("/post", post)
}

// rootTmpl is the main (and only) HTML template.
var rootTmpl = template.Must(template.New("root.html").ParseFiles("tmpl/root.html"))

// root is an HTTP handler that joins or creates a Room,
// creates a new Client, and writes the HTML response.
func root(w http.ResponseWriter, r *http.Request) {
	// Get the name from the request URL.
	name := r.URL.Path[1:]
	// If no valid name is provided, show an error.
	if !validName.MatchString(name) {
		http.Error(w, "Invalid room name", 404)
		return
	}

	c := appengine.NewContext(r)

	// Get or create the Room.
	room, err := getRoom(c, name)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Create a new Client, getting the channel token.
	token, err := room.AddClient(c, randId(clientIdLen))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Render the HTML template.
	data := struct{ Room, Token string }{room.Name, token}
	err = rootTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// post broadcasts a message to a specified Room.
func post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get the room.
	room, err := getRoom(c, r.FormValue("room"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Send the message to the clients in the room.
	err = room.Send(c, r.FormValue("msg"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
