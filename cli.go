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
	mdbook "github.com/ngyewch/mdbook-plugin"
)

type mode string

const (
	standalone mode = "standalone" // running as a standalone CLI tool
	plugin     mode = "plugin"     // running as a plugin for mdbook

	BOOK_TOML    = "book.toml"
	pluginName   = "sitemap-generator"
	domainOption = "domain"
)

type options struct {
	Domain    string
	Output    string
	Version   bool
	mode      mode
	pluginCtx *mdbook.RenderContext
}

func (o *options) SanitizeDomain() error {
	u, err := url.Parse(o.Domain)
	if err != nil {
		return fmt.Errorf("error: invalid value %q: %w", o.Domain, err)
	}
	if u.Scheme == "http" {
		o.Domain = strings.Replace(o.Domain, "http", "https", 1)
	} else if u.Scheme != "https" {
		o.Domain = "https://" + o.Domain
	}
	return nil
}

// representation of the book configuration (i.e., book.toml)
type bookConfig struct {
	Book struct {
		Src string
	}
}

func main() {
	opts, err := getOptions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if opts.mode == standalone && opts.Version {
		fmt.Println(version.Get())
		os.Exit(0)
	}

	if opts.Domain == "" {
		if opts.mode == standalone {
			flag.Usage()
		} else {
			fmt.Println("Domain is required. Please provide the domain in book.toml")
		}
		os.Exit(1)
	}

	if err := opts.SanitizeDomain(); err != nil {
		fmt.Printf("Error sanitizing domain: %v\n", err)
		os.Exit(1)
	}

	paths, err := getPaths(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sitemap := NewSitemap(getUrls(opts.Domain, paths))
	out, err := sitemap.ToXML()
	if err != nil {
		fmt.Printf("Error generating sitemap: %v\n", err)
		os.Exit(1)
	}
	switch {
	case opts.Output == "":
		fmt.Println(string(out))
	case opts.Output != "":
		if err := os.WriteFile(opts.Output, out, 0644); err != nil {
			fmt.Printf("Error writing sitemap content to %q: %v\n", opts.Output, err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

func getOptions() (*options, error) {
	var (
		opts    *options
		domain  = flag.String(domainOption, "", "Domain of the target (e.g., 'example.com') [required]")
		output  = flag.String("output", "", "File to write the output to, defaults to stdout [optional]")
		version = flag.Bool("version", false, "Print the version of the CLI [optional]")
	)
	flag.Parse()

	if *domain == "" && *output == "" && !*version {
		// fallback to plugin mode if no flags are provided
		ctx, err := mdbook.ParseRenderContextFromReader(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("error parsing stdin: %w", err)
		}
		domain, err := getDomainFromPluginContext(ctx)
		if err != nil {
			return nil, err
		}
		opts = &options{
			mode:      plugin,
			Domain:    domain,
			Output:    path.Join(ctx.Destination, "sitemap.xml"),
			pluginCtx: ctx,
		}
	} else {
		opts = &options{
			Domain:  *domain,
			Output:  *output,
			Version: *version,
			mode:    standalone,
		}
	}
	return opts, nil
}

func getDomainFromPluginContext(ctx *mdbook.RenderContext) (string, error) {
	if ctx == nil || ctx.Config == nil || ctx.Config.Output == nil {
		return "", fmt.Errorf("missing 'Config.Output' in plugin context")
	}
	for k, v := range ctx.Config.Output {
		if k != pluginName {
			continue
		}
		pluginCfg := v.(map[string]any)
		for kk, vv := range pluginCfg {
			if kk == domainOption {
				if domain, ok := vv.(string); ok {
					return domain, nil
				}
				return "", fmt.Errorf("'domain' value must be a string")
			}
		}
	}
	return "", fmt.Errorf("'domain' is not defined in the `[output.%s]` section of the book.toml", pluginName)
}

func getPaths(options *options) ([]string, error) {
	if options.mode == standalone {
		return getPathsFromStandalone()
	} else if options.mode == plugin {
		return getPathsFromPluginContext(options.pluginCtx)
	}
	return nil, fmt.Errorf("unsupported mode: %q", options.mode)
}

func getPathsFromStandalone() ([]string, error) {
	targetDir := "src"
	if _, err := os.Stat(getbookConfigPath()); err == nil {
		cfg := new(bookConfig)
		if _, err := toml.DecodeFile(getbookConfigPath(), cfg); err != nil {
			return nil, fmt.Errorf("error reading %q: %v\n", BOOK_TOML, err)
		}
		targetDir = cfg.Book.Src
	}
	paths, err := findPaths(targetDir, "")
	if err != nil {
		return nil, fmt.Errorf("error reading %q directory: %v\n", targetDir, errors.Unwrap(err))
	}
	return paths, nil
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

func getPathsFromPluginContext(ctx *mdbook.RenderContext) ([]string, error) {
	if ctx == nil || ctx.Book == nil || ctx.Book.Sections == nil {
		return nil, fmt.Errorf("invalid plugin context: missing Book or Sections")
	}
	paths := findPathsFromPluginContext(ctx.Book.Sections)
	return paths, nil
}

func findPathsFromPluginContext(items []*mdbook.BookItem) []string {
	paths := []string{}
	for _, item := range items {
		paths = append(paths, item.Chapter.SourcePath)
		if len(item.Chapter.SubItems) > 0 {
			subPaths := findPathsFromPluginContext(item.Chapter.SubItems)
			paths = append(paths, subPaths...)
		}
	}
	return paths
}

func getUrls(domain string, paths []string) []string {
	urls := []string{domain}
	for _, p := range paths {
		filename := strings.Replace(p, ".md", ".html", 1)
		urls = append(urls, fmt.Sprintf("%s/%s", domain, filename))
	}
	return urls
}

func getbookConfigPath() string {
	wd, _ := os.Getwd()
	return path.Join(wd, BOOK_TOML)
}
