// Roughly based on the Go lang tutorial for web services at: http://golang.org/doc/articles/wiki/
package examplewiki

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

// Main wiki controller.
func wiki() {

	/* Switched to a closure to load and validate the page's title instead. See 'makeHandler'
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	*/
	// Updated way to call the various route handlers using a wrapping closure function.
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	/* Turn this off when App is being started by the Google App Engine as it's not needed.
	http.ListenAndServe(":8080", nil)
	*/
}

// Define the data structure of a page.
// NOTE: 'Body' is a byte slice because that's what the 'io' package expect.
type Page struct {
	Title string
	Body  []byte
}

// Saves the page from memory to disk. Filename is the title of the page.
func (p *Page) save() error {
	filename := "wiki/data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Loads the page from the Google App Engine datastore into memory. Uses the page's title to reference the page.
func loadPage(title string, req *http.Request) (*Page, error) {

	// Retrieve the page from the Google App Engine datastore.
	datastoreContext := appengine.NewContext(req)

	// Contruct a query for the requested page. Will return all revisions of the page with that title.
	query := datastore.NewQuery("page").Filter("Title =", title)

	// Initialize an empty page array to hold the results of the query.
	var pagesFromDatastore []*Page

	// Run the query.
	_, err := query.GetAll(datastoreContext, &pagesFromDatastore)

	// Handle query errors.
	if err != nil {
		return nil, err
	}
	// If no existing pages match the query then return an error.
	if pagesFromDatastore == nil {
		err := errors.New("Page does not exist.")
		return nil, err
	}

	// Output to the consule how many revisions of the page have been made.
	fmt.Printf("\nNumber of revisions: %s", len(pagesFromDatastore))

	// Iterate through the array of Pages returned by the search. Output to the consule 
	// the values of the previous versions of the page.
	for x := range pagesFromDatastore {
		fmt.Printf("\nVersion %s = Title: \"%s\" - Body: \"%s\" \n", x, pagesFromDatastore[x].Title, pagesFromDatastore[x].Body)
	}

	// Return the latest version to the caller so that it can be rendered.
	return &Page{Title: pagesFromDatastore[len(pagesFromDatastore)-1].Title, Body: pagesFromDatastore[len(pagesFromDatastore)-1].Body}, nil

}

////// HTTP Logic //////

// Set the paths that the handlers are working at.
const lenPath = len("/view/")

// Cache all of the HTML files in the templates directory so that we only have to hit disk once.
var cached_templates = template.Must(template.ParseGlob("wiki/templates/*.html"))

// Set the valid title rules so that user's cannot request just any file they like from the server. #Security
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

// Wrapper function that takes a function like our (view, edit, and save functions) and returns a 
// function of type http.HandlerFunc (suitable to be passed to the function http.HandleFunc).
//
// The returned function is called a closure because it encloses values defined outside of it. 
// In this case, the variable fn (the single argument to makeHandler) is enclosed by the closure. 
// The variable fn will be one of our save, edit, or view handlers.
//
// This allows us to validate function input parameters, such as the 'title' string before calling the function.
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	// Returned the loaded function with it's inputs validated.
	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve the page name the user requested. 
		title := req.URL.Path[lenPath:]

		// For security reasons, make certain that the title is valid.
		if !titleValidator.MatchString(title) {
			http.NotFound(res, req)
			return
		}

		// Call the originally requested function, passing it the validated title.
		fn(res, req, title)
	}
}

/* Switched to a closure to load and validate the page's title instead. See 'makeHandler'
// Retrieve the page's name and check that it is a valid name.
func getTitle(res http.ResponseWriter, req *http.Request) (title string, err error) {

	// Retrieve the page name the user requested.
	title = req.URL.Path[lenPath:]

	// For security reasons, make certain that the title is valid.
	if !titleValidator.MatchString(title) {
		http.NotFound(res, req)
		err = errors.New("\"" + title + "\"" + " is an invalid page title.")
	}

	// Return the page name.
	return
}
*/

// HTTP Handler for viewing wiki pages.
func viewHandler(res http.ResponseWriter, req *http.Request, title string) {

	/* Switched to a closure to load and validate the page's title instead. See 'makeHandler'
	// Retrieve the page name the user requested.
	title, err := getTitle(res, req)

	// Handle bad data requests.
	if err != nil {
		return
	}
	*/

	// Load the page from disk.
	p, err := loadPage(title, req)

	// If the page does not exist then redirect the user to a blank page for adding a new page to the wiki.
	if err != nil {
		http.Redirect(res, req, "/edit/"+title, http.StatusFound)
		return
	}

	// If the page exists then render using a html template.
	renderTemplate(res, "view", p)
}

// HTTP Handler for editing wiki pages.
func editHandler(res http.ResponseWriter, req *http.Request, title string) {

	/* Switched to a closure to load and validate the page's title instead. See 'makeHandler'
	// Retrieve the page name the user requested.
	title, err := getTitle(res, req)

	// Handle bad data requests.
	if err != nil {
		return
	}
	*/

	// Retrieve the requested page.
	p, err := loadPage(title, req)

	// If the page does not yet exist, then let the user create it.
	if err != nil {
		p = &Page{Title: title}
	}

	// Use a html template.
	renderTemplate(res, "edit", p)
}

// HTTP handler for saving wiki pages edits / creations.
func saveHandler(res http.ResponseWriter, req *http.Request, title string) {

	/* Switched to a closure to load and validate the page's title instead. See 'makeHandler'
	// Retrieve the page name the user requested.
	title, err := getTitle(res, req)

	// Handle bad data requests.
	if err != nil {
		return
	}
	*/

	// Retrive the submitted form's values.
	body := req.FormValue("body")

	// Construct the updated page.
	// NOTE: 'Body' is a byte slice because that's what the 'io' package expect.
	p := &Page{
		Title: title,
		Body:  []byte(body)}

	/* Switched from disk to datastore for storing the content.
	// Save the page to disk.
	err := p.save()
	*/

	// Save the page as the latest revision to the Google App Engine datastore.
	datastoreContext := appengine.NewContext(req)
	key, err := datastore.Put(datastoreContext, datastore.NewIncompleteKey(datastoreContext, "page", nil), p)

	// If there was an error saving the page to disk, then return HTTP Internal Server Error error.
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify that it stored by retrieving it's value from the datastore and then inform the consule.
	var p2 Page
	err = datastore.Get(datastoreContext, key, &p2)
	if err == nil {
		fmt.Printf("Verified Storage of Page: "+p2.Title+" -- %s", p2.Body)
	}

	// Redirect to the updated page.
	http.Redirect(res, req, "/view/"+title, http.StatusFound)
}

// Extract the rendering of templates code to reduce boiler-plate for the various route handlers.
func renderTemplate(res http.ResponseWriter, template_name string, p *Page) {

	// Retrieve the specified template from the 'cached_templates' and pipe the Page to the client. 
	err := cached_templates.ExecuteTemplate(res, template_name+".html", p)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	/* Commented out because we switched to a template caching method.

	// Use Go's html/template engine.
	t, err := template.ParseFiles("wiki/templates/" + template_name + ".html")

	// If there is an error in the Page we're trying to pipe back to the client, then return HTTP Internal Server Error error.
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Pipe the template to the user.
	err = t.Execute(res, p)

	// If there was an error rendering the template, then return HTTP Internal Server Error error.
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	*/

}
