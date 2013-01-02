package examplewiki

import (
	/* Switched to a closure to load and validate the page's title instead. See 'makeHandler'
	"errors"
	*/
	//"appengine"
	//"appengine/user"
	"fmt"
	//"html/template"
	//"io/ioutil"
	"net/http"
	//"regexp"
)

// This is where Google App Engine sets up which handler lives at the root url.
func init() {

	// Immediately enter the main app.
	main()
}

func main() {

	// Setup handlers for non-wiki parts of the app.
	http.HandleFunc("/", homepage)
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/_ah/warmup", warmuphandler)
	http.HandleFunc("/jds", jds)

	// Setup wiki routes.
	wiki()

	// Setup cron routes.
	cron()

}

// Base URL is a 'Hello World' that uses your google account name.
func homepage(res http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(res, "Welcome Mate! Go to /view/[topic] to read, edit, or add that topic to the wiki. For example: view/squirrel, view/ruby. view/test")

}

// Test URL.
func jds(res http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(res, "What is love?")

}
