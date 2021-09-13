package gogemgostatic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/gocolly/colly"
)

//TODO If you use local fonts, place them manually in the fonts folder in the created project
/*

	Download all files from the given url and save them to the given path.
	Tested with WordPress, should also work with other websites.
	At the moment only links on html pages are regarded.
	This can lead to problems if files are only included via a css file (aka fonts) or via js files

*/
func GoStatic(url, path string) (string,error) {
	url = sanitize_url(url)
	project_path, err := createProject(path, url)
	if err != nil {
		return "", err
	}

	pages, remove, err := crawlDomain(url)
	if err != nil {
		return "", err
	}

	err = createFileLinks(pages, url)
	if err != nil {
		return "", err
	}

	err = fetchPages(pages, remove, project_path)
	if err != nil {
		return "", err
	}

	return project_path, nil
}

/*

	Creates uniform folder structure for the project.
	After creation it's empty.

*/
func createProject(path, url string) (string, error) {

	domain, err := urlToDomain(url)
	if err != nil {
		return "", err
	}

	if path == "" {
		if path, err = os.Getwd(); err != nil { // Set path to current working directory
			return "", err
		}
	}
	path = path + "/" + domain
	pathCSS := path + "/css"
	pathJS := path + "/js"
	pathAssets := path + "/assets"
	pathFonts := path + "/fonts"

	if err = makeDir(pathCSS); err != nil {
		return "", err
	}
	if err = makeDir(pathJS); err != nil {
		return "", err
	}
	if err = makeDir(pathAssets); err != nil {
		return "", err
	}
	if err = makeDir(pathFonts); err != nil {
		return "", err
	}
	return path, nil
}

/*

	Crawl the domain and create a map of all pages.
	Using colly.

*/
func crawlDomain(url string) (pages, remove map[string]string, err error) { // Crawl domain
	pages = make(map[string]string) // Map of all found page links to file/type
	remove = make(map[string]string) // Map of all links that need to be removed

	domain, err := urlToDomain(url)
	if err != nil {
		return nil, nil, err
	}

	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) { // Register callback functions for all types of links
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("link[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("video[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("audio[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("header[style]", func(e *colly.HTMLElement) { // Fix for featured images
		link := e.Attr("style")
		urlRegex := regexp.MustCompile(`.*url\((.*?)\).*;`)
		replace := `${1}`
		link = urlRegex.ReplaceAllString(link, replace)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnResponse(func(r *colly.Response) { // When we get a response create pages list
		if r.StatusCode != 200 {
			println(fmt.Sprint(r.StatusCode) + " " + r.Request.URL.String())
		}
		filetype := r.Headers.Get("Content-Type")
		if filetype == "image/svg+xml" || (!strings.Contains(filetype, "json") && !strings.Contains(filetype, "xml")) { // Skipping JSON and XML files, as JSON is the response from the WP REST API, and XML the response of the legacy XML-RPC API
			pages[r.Request.URL.String()] = filetype
		} else {
			remove[r.Request.URL.String()] = ""
		}
	})

	c.Visit(url) // Start Crawling from the given URL

	return pages, remove, nil

}
/*

	Deconstruct given url to relative path, while deligating by filetype to different subfolders.

*/
func createFileLinks(pages map[string]string, url string) error {

	domain, err := urlToDomain(url)
	if err != nil {
		return err
	}

	for link, filetype := range pages {
		if strings.Contains(filetype, "text/html") {
			fragments := strings.Split(link, "/")
			fragments = delete_empty(fragments)
			filename := fragments[len(fragments)-1] + ".html"
			if strings.Contains(filename, domain) {
				filename = "index.html"
			}
			pages[link] = "./" + filename
		} else if strings.Contains(filetype, "text/css") {
			fragments := strings.Split(link, "/")
			fragments = delete_empty(fragments)
			filename := fragments[len(fragments)-1]
			filename = strings.Split(filename, "?")[0]
			pages[link] = "./css/" + filename
		} else if strings.Contains(filetype, "javascript") {
			fragments := strings.Split(link, "/")
			fragments = delete_empty(fragments)
			filename := fragments[len(fragments)-1]
			filename = strings.Split(filename, "?")[0]
			pages[link] = "./js/" + filename
		} else {
			fragments := strings.Split(link, "/")
			fragments = delete_empty(fragments)
			filename := fragments[len(fragments)-1]
			pages[link] = "./assets/" + filename
		}
	}

	return nil

}
/*
	Fetch all pages from the given list of pages and create the files
	Reomve all URLs from the files specified in the remove list
	Replace all absolut links given in the pages list with relative links, which are the second parameter in the pages list
*/
func fetchPages(pages, remove map[string]string, path string) error {

	for link, rel_link := range pages {
		resp, err := http.Get(link)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp_body := string(body)

		ordered_key_list_pages := orderMapKeys(pages)
		ordered_key_list_remove :=  orderMapKeys(remove)

		for _,key := range ordered_key_list_remove {
			resp_body = strings.ReplaceAll(resp_body, key, "")
		}

		for _,key := range ordered_key_list_pages {
			rep_rel_link := pages[key]

			if strings.Contains(rel_link, "css") || strings.Contains(rel_link, "js") || strings.Contains(rel_link, "assets") {
				rep_rel_link = strings.ReplaceAll(rep_rel_link, "./", "./../")
				resp_body = strings.ReplaceAll(resp_body, key, rep_rel_link)
			} else {
				resp_body = strings.ReplaceAll(resp_body, key, rep_rel_link)
			}
		}

		resp_body = strings.ReplaceAll(resp_body, "href=\"/#", "href=\"#") // Fix links to anchors

		file_path := path + strings.Replace(rel_link, "./", "/", -1)

		err = ioutil.WriteFile(file_path, []byte(resp_body), 0644)
		if err != nil {
			return err
		}
	}
	return nil

}

/*
	Try to remove query information and anchors from url
*/
func sanitize_url(url string) string {
	if strings.Contains(url, "?") {
		url = strings.Split(url, "?")[0]
	}
	if strings.Contains(url, "#") {
		url = strings.Split(url, "#")[0]
	}
	if strings.Contains(url, "&") {
		url = strings.Split(url, "&")[0]
	}
	if strings.Contains(url, "=") {
		url = strings.Split(url, "=")[0]
	}
	if len(delete_empty(strings.Split(url, "/"))) == 2 {
		if url[len(url)-1] != '/' {
			url = url + "/"
		}
	}
	return url
}

/*

	Try to convert a URL to a domain name.

*/
func urlToDomain(URL string) (string, error) {
	split := strings.Split(URL, ":")
	if len(split) == 1 || !strings.Contains(split[0], "http") {
		return "", errors.New("no protocol information found")
	}
	if len(split) > 2 {
		return "", errors.New("please remove port information")
	}
	domain := split[1]
	split = strings.Split(domain, "/")
	if len(split) < 3 {
		return "", errors.New("malformed url")
	}
	return split[2], nil
}

/*
	Create Folder and all necessary parents
*/
func makeDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) { // Create project directory
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

/*
	delete empty strings from a slice
*/
func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

/*
	Order the keys of a map by length, the first one is the longest
*/
func orderMapKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
		sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	return keys
}