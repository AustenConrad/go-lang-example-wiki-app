package examplewiki

import (
	/* Switched to a closure to load and validate the page's title instead. See 'makeHandler'
	"errors"
	*/
	"appengine"
	"appengine/user"
	"fmt"
	//"html/template"
	//"io/ioutil"
	"net/http"
	//"regexp"
)

// Base URL is a 'Hello World' that uses your google account name.
func adminHandler(res http.ResponseWriter, req *http.Request) {

	// User's Login with their gmail accounts.
	c := appengine.NewContext(req)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, req.URL.String())
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusFound)
		return
	}

	fmt.Fprintf(res, "Hello %v, you are an admin for this site :)", u)

}
