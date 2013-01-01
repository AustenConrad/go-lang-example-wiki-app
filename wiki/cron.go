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

// Manage the routes for the various cron tasks.
func cron() {

	http.HandleFunc("/cron/heartbeat/cloudant", cloudantHandler)
	http.HandleFunc("/cron/heartbeat/rackspace", rackspaceHandler)

}

// Checks that we have a connection to cloudant.
func cloudantHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(res, "TODO: Check that we have a network connetion to cloudant. Alert if we don't. Perhaps log this data in a google database for charting later.")

}

// Checks that we have a connection to cloudant.
func rackspaceHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(res, "TODO: Check that we have a network connetion to the Rackspace Cloud API. Alert if we don't. Perhaps log this data in a google database for charting later.")

}
