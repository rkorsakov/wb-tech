package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type HTMLParser struct {
	baseURL    *url.URL
	baseDomain string
}

func NewHTMLParser(baseURL string) (*HTMLParser, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &HTMLParser{
		baseURL:    parsedURL,
		baseDomain: parsedURL.Hostname(),
	}, nil
}

type ProcessResult struct {
	ModifiedHTML []byte
	ResourceURLs []string
}

func (p *HTMLParser) ProcessHTML(content []byte, baseURL string) (*ProcessResult, error) {
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	var resourceURLs []string
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	processAttribute := func(node *html.Node, attrName string) {
		for i, attr := range node.Attr {
			if attr.Key == attrName {
				absoluteURL, err := p.makeAbsoluteURL(attr.Val, base)
				if err != nil {
					continue
				}

				if p.shouldDownload(absoluteURL) {
					resourceURLs = append(resourceURLs, absoluteURL)
				}

				localPath, err := p.urlToLocalPath(absoluteURL)
				if err == nil {
					node.Attr[i].Val = localPath
				}
			}
		}
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a", "link":
				processAttribute(n, "href")
			case "img", "script":
				processAttribute(n, "src")
			case "iframe", "source", "track":
				processAttribute(n, "src")
			case "form":
				processAttribute(n, "action")
			case "meta":
				if attrValue(n, "property") == "og:image" || attrValue(n, "name") == "twitter:image" {
					processAttribute(n, "content")
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, err
	}

	return &ProcessResult{
		ModifiedHTML: buf.Bytes(),
		ResourceURLs: resourceURLs,
	}, nil
}

func (p *HTMLParser) makeAbsoluteURL(relativeURL string, base *url.URL) (string, error) {
	if relativeURL == "" {
		return "", fmt.Errorf("empty URL")
	}

	parsed, err := url.Parse(relativeURL)
	if err != nil {
		return "", err
	}

	if parsed.Scheme == "javascript" || parsed.Fragment != "" && parsed.Path == "" {
		return "", fmt.Errorf("skip anchor or javascript")
	}

	parsed.Fragment = ""

	absoluteURL := base.ResolveReference(parsed).String()
	return absoluteURL, nil
}

func (p *HTMLParser) shouldDownload(URL string) bool {
	parsed, err := url.Parse(URL)
	if err != nil {
		return false
	}

	if parsed.Hostname() != p.baseDomain {
		return false
	}

	return true
}

func (p *HTMLParser) urlToLocalPath(URL string) (string, error) {
	parsed, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	localPath := filepath.Join(parsed.Host, parsed.Path)

	if strings.HasSuffix(localPath, "/") {
		localPath = filepath.Join(localPath, "index.html")
	}

	ext := filepath.Ext(localPath)
	if ext == "" {
		localPath += ".html"
	}

	return localPath, nil
}

func attrValue(node *html.Node, attrName string) string {
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}