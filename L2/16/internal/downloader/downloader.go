package downloader

import (
	"16/internal/parser"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Config struct {
	URL            string
	RecursionDepth int
	OutputDir      string
}

type Downloader struct {
	config       *Config
	visitedURLs  map[string]bool
	visitedMutex sync.RWMutex
	baseURL      *url.URL
	baseDomain   string
	client       *http.Client
	parser       *parser.HTMLParser
}

func NewDownloader(config *Config) (*Downloader, error) {
	if config.OutputDir == "" {
		config.OutputDir = "content"
	}

	parser, err := parser.NewHTMLParser(config.URL)
	if err != nil {
		return nil, err
	}

	return &Downloader{
		config:      config,
		visitedURLs: make(map[string]bool),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		parser: parser,
	}, nil
}

func DownloadPage(URL string, recursionDepth int) error {
	config := &Config{
		URL:            URL,
		RecursionDepth: recursionDepth,
		OutputDir:      "content",
	}

	downloader, err := NewDownloader(config)
	if err != nil {
		return err
	}

	return downloader.Download()
}

func (d *Downloader) Download() error {
	parsedURL, err := url.Parse(d.config.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	d.baseURL = parsedURL
	d.baseDomain = parsedURL.Hostname()

	if err := os.MkdirAll(d.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return d.downloadRecursive(d.config.URL, 0)
}

func (d *Downloader) downloadRecursive(URL string, depth int) error {

	if depth > d.config.RecursionDepth {
		return nil
	}

	if d.isVisited(URL) {
		return nil
	}
	d.markVisited(URL)

	fmt.Printf("Downloading: %s (depth: %d)\n", URL, depth)

	content, contentType, err := d.downloadResource(URL)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", URL, err)
	}

	localPath, err := d.urlToLocalPath(URL)
	if err != nil {
		return err
	}

	if strings.Contains(contentType, "text/html") {
		result, err := d.parser.ProcessHTML(content, URL)
		if err != nil {
			return fmt.Errorf("failed to process HTML %s: %w", URL, err)
		}
		content = result.ModifiedHTML

		for _, resourceURL := range result.ResourceURLs {
			if d.shouldDownload(resourceURL) {
				err := d.downloadRecursive(resourceURL, depth+1)
				if err != nil {
					fmt.Printf("Warning: failed to download resource %s: %v\n", resourceURL, err)
				}
			}
		}
	}

	if err := d.saveFile(localPath, content); err != nil {
		return err
	}

	return nil
}

func (d *Downloader) shouldDownload(URL string) bool {
	parsed, err := url.Parse(URL)
	if err != nil {
		return false
	}

	if parsed.Hostname() != d.baseDomain {
		return false
	}

	if d.isVisited(URL) {
		return false
	}

	return true
}

func (d *Downloader) isVisited(URL string) bool {
	d.visitedMutex.RLock()
	defer d.visitedMutex.RUnlock()
	return d.visitedURLs[URL]
}

func (d *Downloader) markVisited(URL string) {
	d.visitedMutex.Lock()
	defer d.visitedMutex.Unlock()
	d.visitedURLs[URL] = true
}

func (d *Downloader) downloadResource(URL string) ([]byte, string, error) {
	resp, err := d.client.Get(URL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return content, resp.Header.Get("Content-Type"), nil
}

func (d *Downloader) saveFile(path string, content []byte) error {

	fullPath := filepath.Join(d.config.OutputDir, path)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, content, 0644)
}

func (d *Downloader) urlToLocalPath(URL string) (string, error) {
	parsed, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	localPath := filepath.Join(parsed.Host, parsed.Path)

	if strings.HasSuffix(localPath, "/") || localPath == parsed.Host {
		localPath = filepath.Join(localPath, "index.html")
	}

	ext := filepath.Ext(localPath)
	if ext == "" {
		localPath += ".html"
	}

	return localPath, nil
}