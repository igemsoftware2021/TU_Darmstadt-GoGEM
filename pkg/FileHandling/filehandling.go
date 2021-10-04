package gogemfilehandling

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	h "github.com/Jackd4w/goGEM/pkg/Handler"
)

/*
	Prepares pages for upload, and begins uploading the media files... (This dual purpose really gives me a headache, but iGEM randomizes the absolut url to the media files, and there is no other way than uploading to safely replace all links)
	Cuts out all empty link references (these get created when converting a wp site to a static one and are remnants of the WP APIs)
	Replaces the HTML DOCTYPE declaration with the standard iGEM template for the team (i.e. {{teamname}})
	Removes all srcsets, these are good for optimization but dramatically increase the difficulty of uploading images to igem
	Replace pageextensions: We can not easily upload JavaScript to the server and request it, because all our Files are just pages on the iGEM Wiki.

*/

var blacklist = make(map[string]string)

func PrepFilesForIGEM(teamname, root string, client *h.Handler) error {

	// Get all files in the root directory
	files, err := allFilesInDir(root)
	if err != nil {
		return err
	}

	for _, filepath := range files {
		println("Preparing file: " + filepath)
		file, err := os.Open(filepath) // Open file
		if err != nil {
			return err
		}
		defer file.Close()

		content, err := ioutil.ReadAll(file) // Read file
		if err != nil {
			return err
		}
		var newContent = string(content) // Create newContent string from content byte array

		if strings.Contains(filepath, ".css") { // If file is a css file
			continue
		}
		newContent = removeAllEmptyLinks(newContent)
		newContent = removeRemoveLinks(newContent)
		newContent = replaceDoctypeWithTemplate(newContent, teamname)
		newContent = removeSrcSet(newContent)
		newContent = removeInlineWP(newContent)
		newContent = replacePageExtensions(newContent)

		fileLinks := findAllFileLinks(newContent)
		fileAssociations, err := fileUpload(fileLinks, root, client) // Output from FileUpload method takes fileLinks as input
		if err != nil {
			return err
		}

		// DEBUG - Prints all file associations
		// println(" ")
		// println("--------------------------------------------------------------")
		// println(filepath)
		// println("--------------------------------------------------------------")
		// for key, value := range fileAssociations {
		// 	println(key + ": " + value)
		// }

		newContent = replaceAllFileLinks(newContent, fileAssociations)

		file, err = os.Create(file.Name())
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(newContent); err != nil {
			return err
		}
	}
	println("File Upload: Done")
	for _, filepath := range files {
		err := pageUpload(filepath, client)
		if err != nil {
			return err
		}
	}
	return nil
}

// Creates list of all files in a directory, and its respective subdirectories.
func allFilesInDir(path string) ([]string, error) {
	var files []string

	objects, err := os.ReadDir(path) // Read the root directory
	if err != nil {
		return nil, err
	}
	for _, dir := range objects { // For each entry in the root directory
		if dir.IsDir() { // If the entry is a directory
			subFiles, err := allFilesInDir(path + "/" + dir.Name()) // Call this function again, but with the subdirectory
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...) // Append all files in the subdirectory to the list
		} else {
			files = append(files, path+"/"+dir.Name()) // Append the file to the list
		}
	}
	return files, nil // Return the list of files
}

func findAllFileLinks(newContent string) []string {
	var fileLinks []string
	srcRegEx := regexp.MustCompile(`src=("|')(.*?)("|')`) // Regex to find all src attributes
	srcLinks := srcRegEx.FindAllString(newContent, -1)    // Find all src attributes

	for _, link := range srcLinks {
		link = srcRegEx.ReplaceAllString(link, `${2}`) // Replace src attribute with just the path
		if strings.Contains(link, "assets") {          // If link is not a css, js, json or html file, append it to fileLinks)
			fileLinks = append(fileLinks, link) // Append all links to fileLinks
		}
	}
	// Fix for featured images
	urlRegEx := regexp.MustCompile(`url\((.*?)\)`)     // Regex to find all url attributes
	urlLinks := urlRegEx.FindAllString(newContent, -1) // Find all url attributes
	for _, link := range urlLinks {
		link = urlRegEx.ReplaceAllString(link, `${1}`) // Replace url attribute with just the path
		if strings.Contains(link, "assets") {          // If link is not a css, js, json or html file, append it to fileLinks)
			fileLinks = append(fileLinks, link) // Append all links to fileLinks
		}
	}

	fileLinks = removeDuplicateStr(fileLinks)
	return fileLinks // Return new array
}

func removeAllEmptyLinks(newContent string) string {
	emptyHrefRegEx := regexp.MustCompile(`<.*?href="".*?\>`)
	newContent = emptyHrefRegEx.ReplaceAllString(newContent, "")
	return newContent
}

func removeSrcSet(newContent string) string {
	srcSetRegEx := regexp.MustCompile(`srcset=".*?"`)
	sizesRegEx := regexp.MustCompile(`sizes=".*?"`)

	newContent = srcSetRegEx.ReplaceAllString(newContent, "")
	newContent = sizesRegEx.ReplaceAllString(newContent, "")

	return newContent
}

func removeInlineWP(newContent string) string {
	styleRegEx := regexp.MustCompile(`(?s)<script>.*?</style>`)
	newContent = styleRegEx.ReplaceAllString(newContent, "")
	return newContent
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func removeRemoveLinks(newContent string) string{
	removeRegEx := regexp.MustCompile(`<a class="remove" .*?<\/a>`)
	newContent = removeRegEx.ReplaceAllString(newContent, "")
	return newContent
}

func ReplaceAllFileLinksDebug(newContent string, fileAssociations map[string]string) string {
	return replaceAllFileLinks(newContent, fileAssociations)
}

/*
	searches for src attributes that do not reference pages, returns list of these links
*/
func replaceAllFileLinks(newContent string, fileLinks map[string]string) string {

	for org, new := range fileLinks {
		newContent = strings.ReplaceAll(newContent, org, new)
	}
	return newContent

}

/*
	Replaces DOCTYPE with the iGEM Standardtemplate of the team
*/
func replaceDoctypeWithTemplate(newContent, teamname string) string {
	newContent = strings.Replace(newContent, "<!DOCTYPE html>", "{{"+teamname+"}}", 1)
	newContent = strings.Replace(newContent, "<!doctype html>", "{{"+teamname+"}}", 1)
	return newContent
}

/*

We can not easily upload JavaScript to the server and request it, because all our Files are just pages on the iGEM Wiki.
Therefor we have to request the raw HTML of the page from the server, as this is only the text we entered, without the iGEM additions (i.e. the mandatory iGEM Nav-Bar etc.).
But even if we do that, iGEM tries to prevent the unintended load of JS by checking the raw conent and recognizing if it is a script, preventing display if the content type does not match.

*/
func replacePageExtensions(newContent string) string {
	cssRegex := regexp.MustCompile(`((href|src)=("|').*?)(\.css)("|')`)         // Regex to find all relative referenced css files
	mincssRegex := regexp.MustCompile(`((href|src)=("|').*?)(\.min\.css)("|')`) // Regex to find all relative referenced css files
	jsRegex := regexp.MustCompile(`((src|href)=("|').*)(\.js)(\?.*?)?("|')`)    // Regex to find all relative referenced js files
	minjsRegex := regexp.MustCompile(`((src|href)=("|').*)(\.min\.js)(\?.*?)?("|')`)
	indexRegEx := regexp.MustCompile(`index\.html`) // Regex to find all href and src attributes that reference index.html
	htmlRegEx := regexp.MustCompile(`\.html`)

	cssReplace := `${1}?action=raw&ctype=text/css${3}`        // Replace all relative css paths with ?action=raw&ctype=text/css, requesting the raw file from the server with the right content type
	mincssReplace := `${1}-min?action=raw&ctype=text/css${3}` // Replace all relative css paths with ?action=raw&ctype=text/css, requesting the raw file from the server with the right content type
	jsReplace := `${1}?action=raw&ctype=text/javascript${3}`
	minjsReplace := `${1}-min?action=raw&ctype=text/javascript${3}`
	htmlReplace := ``

	newContent = mincssRegex.ReplaceAllString(newContent, mincssReplace) // Replace all '.css' in relative paths with ?action=raw&ctype=text/css, requesting the raw file from the server with the right content type
	newContent = cssRegex.ReplaceAllString(newContent, cssReplace)       // Replace all '.css' in relative paths with ?action=raw&ctype=text/css, requesting the raw file from the server with the right content type
	newContent = minjsRegex.ReplaceAllString(newContent, minjsReplace)   // Replace all '.min.js' in relative paths with ?action=raw&ctype=text/javascript, requesting the raw file from the server with the right content type
	newContent = jsRegex.ReplaceAllString(newContent, jsReplace)         // Replace all '.js' in relative paths with ?action=raw&ctype=text/javascript, requesting the raw file from the server with the right content type
	newContent = indexRegEx.ReplaceAllString(newContent, htmlReplace)    // Replace all href and src attributes that reference index.html with empty string
	newContent = htmlRegEx.ReplaceAllString(newContent, htmlReplace)     // Replace all '.html' in relative paths with empty string, so the raw file is requested from the server without the .html extension

	return newContent
}

func fileUpload(fileLinks []string, root string, client *h.Handler) (map[string]string, error) {
	result := make(map[string]string)
	local_blacklist := make(map[string]bool)

	for _, link := range fileLinks {
		path := root + link[1:] // Remove leading slash
		path = strings.ReplaceAll(path, "/", `\`)

		res_url := ""
		// fmt.Println(blacklist)
		if blacklist[path] != "" && !local_blacklist[path] {
			local_blacklist[path] = true
			result[link] = blacklist[path]
			continue
		}

		if !local_blacklist[path] {
			url, err := client.UploadFile(path, false)
			if err != nil {
				if err.Error() == "alreadyUploadedInThisSession" || err.Error() == "fileAlreadyUploaded" {
					local_blacklist[path] = true
					res_url = client.GetFileUrl(url)
					blacklist[path] = res_url
					result[link] = res_url
					continue
					// return nil, err
				} else {
					println("Error " + err.Error() + " uploading file: " + path)
				}
			}
			res_url = client.GetFileUrl(url)

		}
		println("Uploaded file: " + res_url)
		result[link] = res_url
		blacklist[path] = res_url
	}

	return result, nil
}

func pageUpload(filepath string, client *h.Handler) error {
	offset := ""
	filename := filepath[strings.LastIndex(filepath, "/")+1:]
	if isPage(filename) {
		if offset != "" {
			if strings.Contains(filename, ".css") {
				offset = "/css"
			} else if strings.Contains(filename, ".js") {
				offset = "/js"
			}
		} else {
			if strings.Contains(filename, ".css") {
				offset = "css"
			} else if strings.Contains(filename, ".js") {
				offset = "js"
			}
		}

		filepath = strings.ReplaceAll(filepath, "/", `\`)
		url, err := client.Upload(filepath, offset, false)
		if err != nil {
			if err.Error() == "alreadyUploadedInThisSession" || err.Error() == "fileAlreadyUploaded" {
				return nil
			}
			return err
		}
		println("Uploaded page: " + url)
	}
	return nil

}

func isPage(filepath string) bool {
	if strings.Contains(filepath, ".html") || strings.Contains(filepath, ".htm") || strings.Contains(filepath, ".css") || strings.Contains(filepath, ".js") {
		return true
	}
	return false
}
