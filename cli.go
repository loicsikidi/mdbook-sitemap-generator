package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/imjasonh/version"
)

const (
	BOOK_TOML = "book.toml"
)

type config struct {
	Book struct {
		Src string
	}
}

var (
	domain   = flag.String("domain", "", "Domain of the target (e.g., https://example.com) [required]")
	output   = flag.String("output", "", "File to write the output to, defaults to stdout [optional]")
	_version = flag.Bool("version", false, "Print the version of the CLI [optional]")
)

func main() {
	flag.Parse()

	if *_version {
		fmt.Println(version.Get())
		os.Exit(0)
	}

	if *domain == "" {
		flag.Usage()
		return
	}

	u, err := url.Parse(*domain)
	if err != nil {
		fmt.Printf("Error parsing domain %q: %v\n", *domain, err)
		os.Exit(1)
	}

	if u.Scheme == "http" {
		*domain = strings.Replace(*domain, "http", "https", 1)
	} else if u.Scheme != "https" {
		*domain = "https://" + *domain
	}

	targetDir := "src"
	if _, err := os.Stat(getConfigPath()); err == nil {
		cfg := new(config)
		if _, err := toml.DecodeFile(getConfigPath(), cfg); err != nil {
			fmt.Printf("Error reading %q: %v\n", BOOK_TOML, err)
			os.Exit(1)
		}
		targetDir = cfg.Book.Src
	}
	paths, err := findPaths(targetDir, "")
	if err != nil {
		fmt.Printf("Error reading %q directory: %v\n", targetDir, errors.Unwrap(err))
		os.Exit(1)
	}

	sitemap := NewSitemap(getUrls(*domain, paths))
	out, err := sitemap.ToXML()
	if err != nil {
		fmt.Printf("Error generating sitemap: %v\n", err)
		os.Exit(1)
	}
	switch {
	case *output == "":
		fmt.Println(string(out))
	case *output != "":
		if err := os.WriteFile(*output, out, 0644); err != nil {
			fmt.Printf("Error writing sitemap content to %q: %v\n", *output, err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

func findPaths(directory, currentPath string) ([]string, error) {
	paths := []string{}

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			subPaths, err := findPaths(path.Join(directory, file.Name()), currentPath)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subPaths...)
		} else if strings.HasSuffix(file.Name(), ".md") && file.Name() != "SUMMARY.md" {
			paths = append(paths, path.Join(currentPath, file.Name()))
		}
	}
	return paths, nil
}

func getUrls(domain string, paths []string) []string {
	urls := []string{domain}
	for _, p := range paths {
		filename := strings.Replace(p, ".md", ".html", 1)
		urls = append(urls, fmt.Sprintf("%s/%s", domain, filename))
	}
	return urls
}

func getConfigPath() string {
	wd, _ := os.Getwd()
	return path.Join(wd, BOOK_TOML)
}
