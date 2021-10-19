package gogemhandler

import (
	"errors"
	"fmt"
	"net/http"

	api "github.com/Jackd4w/GoGEM-WikiAPI"
)

type Handler struct {
	Session         *http.Client
	year            int
	teamname        string
	offset          string
	alreadyUploaded map[string]bool
	loginURL        string
	logoutURL       string
	prefixURL       string
	timeout         int
}

//TODO move blacklist here, this is the place where all files pass, also this persists for the whole runtime

/*
  My Approach at a Session handler, trying to avoide to make obfuscated API calls, and trying to minimize requests to iGEM Servers by holding onto the Session for the runtime of the program.
  Also reduces the information that needs to be passed around.

  Save Session for the runtime of the program. Minimize server requests by keeping track of uploaded files during session
*/
func NewHandler(year, timeout int, username, password, teamname, offset, loginURL, logoutURL, prefixURL string) (*Handler, error) {
	handler := new(Handler)

	handler.loginURL = loginURL
	handler.logoutURL = logoutURL
	handler.prefixURL = fmt.Sprintf(prefixURL, year)
	handler.year = year
	handler.teamname = teamname
	handler.offset = offset
	handler.alreadyUploaded = make(map[string]bool)
	handler.timeout = timeout

	session, err := handler.Login(username, password)
	if err != nil {
		return nil, err
	}

	handler.Session = session

	return handler, nil
}

/*
	Wrappes the Login function in the API package
*/
func (h Handler) Login(username, password string) (*http.Client, error) {
	session, err := api.Login(username, password, h.loginURL, h.timeout)
	if err != nil {
		return nil, err
	}
	return session, nil
}

/*
	Wrappes the Logout function in the API package
*/
func (h Handler) Logout() error {
	return api.Logout(h.Session, h.logoutURL)
}

/*
	Wrappes the Upload function in the API package, also performs a check if a session is open and provides necessary metadata
*/
func (h Handler) Upload(filepath, offset string, force bool) (string, error) {
	if !h.loggedIn() {
		return "", errors.New("notLoggedIn")
	}
	url, err := api.Upload(h.Session, h.year, h.teamname, filepath, h.offset+offset, false, force)
	return url, err
}

func (h Handler) Redirect(source, target string) error {
	_, err := api.Redirect(h.Session, h.year, h.teamname, source, target)
	return err
}

/*
	Wrappes the Upload function in the API package, also performs a check if a session is open and provides necessary metadata.
	The force parameter is used to force the upload of a file that has already been uploaded.
	There is a local check if a file has been uploaded during this session, as this method is called on a per file basis and there will be redundant requests
	Returns the full url to the uploaded file.
*/
func (h Handler) UploadFile(filepath string, force bool) (string, error) {
	if !h.loggedIn() {
		return "", errors.New("notLoggedIn")
	}

	if h.alreadyUploaded[filepath] { // Check if file has already been uploaded in this session
		return "", errors.New("alreadyUploadedInThisSession")
	}

	// println("Uploading file: " + filepath) // Debugging

	url, err := api.Upload(h.Session, h.year, h.teamname, filepath, h.offset, true, force)

	if err == nil {
		h.alreadyUploaded[filepath] = true
	}

	return url, err
}

// Simple Wrapper
func (h Handler) GetFileUrl(url string) string {
	res_url, err := api.GetFileUrl(url, h.Session)
	if err != nil {
		return ""
	}
	return res_url
}

/*
Query all Pages from the specified prefix url that have the specified teamname and offset
*/
func (h Handler) GetAllPages() ([]string, error) {
	return api.QueryPages(h.prefixURL, h.teamname, h.offset, h.Session)
}

/* //TODO correct: there is a tag that gets added, so the checker can recognize deleted pages (?) Look into API
Overwrite the specified pageurl with an empty string, effectively deleting the page (also marking it for eventuell cleanup processes from the hoster side due to it having no user content)
*/
func (h Handler) DeletePage(pageurl string) error {
	return api.DeletePage(pageurl, h.year, h.Session)
}

/*
------------------------------------------------------------------------------
Internal Functions
------------------------------------------------------------------------------
*/

func (h Handler) loggedIn() bool {
	return h.Session != nil
}
