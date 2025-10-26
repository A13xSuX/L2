package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/net/html"
)

type Config struct {
	Url    string
	Depth  int
	Output string
}

type Crawler struct {
	visited map[string]bool
	queue   []DownloadTask
	baseURL *url.URL
	config  Config
	mu      sync.Mutex
}

type DownloadTask struct {
	URL   *url.URL
	Depth int
}

func NewCrawler(config Config) (*Crawler, error) {
	baseUrl, err := url.Parse(config.Url)
	if err != nil {
		return nil, err
	}

	initialTask := DownloadTask{
		URL:   baseUrl,
		Depth: config.Depth,
	}

	crawler := &Crawler{
		visited: make(map[string]bool),
		queue:   []DownloadTask{initialTask},
		baseURL: baseUrl,
		config:  config,
	}

	return crawler, err
}

func (c *Crawler) Start() error {
	for len(c.queue) > 0 {
		task := c.queue[0]
		c.queue = c.queue[1:]
		urlStr := task.URL.String()
		if c.visited[urlStr] {
			continue
		}
		c.visited[urlStr] = true

		data, err := c.downloadUrl(task.URL)
		if err != nil {
			fmt.Printf("Error downloading %s: %s\n", task.URL, err)
			continue
		}

		err = c.saveFile(task.URL, data)
		if err != nil {
			fmt.Printf("Error saving %s: %s\n", task.URL, err)
			continue
		}

		if task.Depth > 0 {
			links := extractLinks(data)
			c.addLinksToQueue(links, task.Depth-1)
		}
	}
	return nil
}

func (c *Crawler) downloadUrl(url *url.URL) ([]byte, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		fmt.Printf("Error downloading %s: %s\n", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("Error downloading %s: %d\n", url, resp.StatusCode)
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func (c *Crawler) saveFile(url *url.URL, data []byte) error {
	var filename string

	if url.Path == "" || url.Path == "/" {
		filename = "index.html"
	} else if url.Path[len(url.Path)-1] == '/' {
		filename = url.Path + "index.html"
	} else {
		filename = url.Path
	}
	if len(filename) > 0 && filename[0] == '/' {
		filename = filename[1:]
	}
	fullPath := filepath.Join(c.config.Output, filename)

	os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)

	return os.WriteFile(fullPath, data, 0644)
}

func (c *Crawler) addLinksToQueue(links []string, depth int) {
	for _, link := range links {
		if link == "" || link[0] == '#' {
			continue
		}

		absoluteURL, err := c.baseURL.Parse(link)
		if err != nil {
			continue
		}

		if absoluteURL.Host != c.baseURL.Host {
			continue
		}

		if absoluteURL.Scheme != "http" && absoluteURL.Scheme != "https" {
			continue
		}

		urlStr := absoluteURL.String()
		if c.visited[urlStr] {
			continue
		}

		newTask := DownloadTask{
			URL:   absoluteURL,
			Depth: depth,
		}
		c.queue = append(c.queue, newTask)

		fmt.Printf("Added to queue: %s (depth: %d)\n", urlStr, depth)
	}
}

func extractLinks(data []byte) []string {
	urls := []string{}

	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return urls
	}
	visitNode(doc, &urls)
	return urls
}

func getAttr(n *html.Node, attrName string) string {
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func visitNode(n *html.Node, urls *[]string) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			if href := getAttr(n, "href"); len(href) > 0 {
				*urls = append(*urls, href)
			}
		case "img":
			if src := getAttr(n, "src"); len(src) > 0 {
				*urls = append(*urls, src)
			}
		case "script":
			if src := getAttr(n, "src"); len(src) > 0 {
				*urls = append(*urls, src)
			}
		case "link":
			if href := getAttr(n, "href"); len(href) > 0 {
				*urls = append(*urls, href)
			}
		case "source":
			if src := getAttr(n, "src"); len(src) > 0 {
				*urls = append(*urls, src)
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		visitNode(child, urls)
	}
}

func main() {
	var config Config
	flag.StringVar(&config.Url, "url", "", "url")
	flag.IntVar(&config.Depth, "depth", 1, "depth")
	flag.StringVar(&config.Output, "out", ".", "output file")

	flag.Parse()

	if config.Url == "" {
		fmt.Println("url is empty")
		os.Exit(1)
	}

	// Создаем выходную директорию
	if config.Output != "" {
		err := os.MkdirAll(config.Output, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir err:%s\n", err)
		}
	}

	crawler, err := NewCrawler(config)
	if err != nil {
		fmt.Printf("Create crawler error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting crawl of %s with depth %d\n", config.Url, config.Depth)
	err = crawler.Start()
	if err != nil {
		fmt.Printf("Crawl error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Crawl completed!")
}
