package main

//Import relevant modules to assist in html parsing and crawling
import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

//define variabled to hold file name
var (
	fileName string
)

//Check if file name is already created in the Downloads Folders
//and prevent parsing further crawling of the already crawled url
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//From the Url, create file names from the url by replacing "/"
//with uderscore to ebale proposer naming of files
func create_file_from_url(url string, resp *http.Response) string {
	fileName = strings.Replace(url, "/", "-", -1)
	fileExt := ".html"
	path := "./downloads/"
	// fmt.Println("File Name %s  %s ", url, path+fileName+fileExt)
	f := path + fileName + fileExt
	if fileExists(f) {
		fmt.Printf("Already Crawled -  ")
	} else {
		create_file_and_download_content(f, resp)

	}
	return f

}

//Collect all links from response body and return it as an array
//of strings to the calling function
func getLinks(body io.Reader) []string {
	var links []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		//Using Switch Control, ensure that Links are crawled
		//successfully and contents of the links downloaded and stored in file
		switch tt {
		case html.ErrorToken:
			//todo: links list shoudn't contain duplicates
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}

				}
			}

		}
	}
}

//Download contents of the url and store that same in the file inside download folders using the OS Create modules
func create_file_and_download_content(fileName string, resp *http.Response) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d", fileName, size)

}

func get_url_contents(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		// log.Fatal(err)
	}
	return resp
}

func main() {

	//Default command line args
	url := flag.String("url", "https://vorozhko.net/", "Url TO parse")
	flag.Parse()

	resp := get_url_contents(*url)
	create_file_from_url(*url, resp)

	// fmt.Println(fileName)

	for _, v := range getLinks(resp.Body) {
		if strings.Contains(v, *url) {
			fmt.Println(v)

			resp := get_url_contents(v)
			// fmt.Println(resp)

			create_file_from_url(v, resp)
		}

	}
}
