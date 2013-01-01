// Roughly based on the Go lang tutorial for web services at: http://golang.org/doc/articles/wiki/
package main

import (
	/* Switched to a closure to load and validate the page's title instead. See 'makeHandler'
	"errors"
	*/
	//"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

func main() {

	/* Switched to a closure to load and validate the page's title instead. See 'makeHandler'
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	*/
	// Updated way to call the various route handlers using a wrapping closure function.
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.ListenAndServe(":8080", nil)
}

// Define the data structure of a page.
// NOTE: 'Body' is a byte slice because that's what the 'io' package expect.
type Page struct {
	Title string
	Body  []byte
}

// Saves the page from memory to disk. Filename is the title of the page.
func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Loads the page from disk into memory. Uses the page's title to reference the page's file name.
func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)

	// Check for any errors that occured while opening the specified file.
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

////// HTTP Logic //////

// Set the paths that the handlers are working at.
const lenPath = len("/view/")

// Set which templates we should cache so that we don't have to 'ParseFiles' them from disk on every page request.
var cached_templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))

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
	p, err := loadPage(title)

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
	p, err := loadPage(title)

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
	p := &Page{Title: title, Body: []byte(body)}

	// Save the page to disk.
	err := p.save()

	// If there was an error saving the page to disk, then return HTTP Internal Server Error error.
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
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
	t, err := template.ParseFiles(template_name + ".html")

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
