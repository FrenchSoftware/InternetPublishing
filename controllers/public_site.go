package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hyperstitieux/template/database/repositories"
	githubpkg "github.com/hyperstitieux/template/github"
	"github.com/hyperstitieux/template/markdown"
	"github.com/hyperstitieux/template/pages"
)

type PublicSiteController struct {
	sites *repositories.SitesRepository
}

func NewPublicSiteController(sites *repositories.SitesRepository) *PublicSiteController {
	return &PublicSiteController{sites: sites}
}

func (c *PublicSiteController) Render(w http.ResponseWriter, r *http.Request) error {
	// Extract slug from subdomain
	host := r.Host
	parts := strings.Split(host, ".")

	// Check if this is a subdomain request (e.g., slug.internetpublishing.co)
	if len(parts) < 3 {
		http.Error(w, "Invalid subdomain", http.StatusBadRequest)
		return nil
	}

	slug := parts[0]

	// Get site from database
	site, err := (*c.sites).GetBySlug(slug)
	if err != nil {
		return err
	}
	if site == nil {
		http.Error(w, "Site not found", http.StatusNotFound)
		return nil
	}

	// Get requested path (default to README.md)
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "README.md"
	}

	// Ensure path ends with .md
	if !strings.HasSuffix(path, ".md") {
		path = path + ".md"
	}

	// Prepend subdirectory if specified
	if site.Subdirectory != "" {
		path = strings.TrimSuffix(site.Subdirectory, "/") + "/" + path
	}

	// Fetch file from GitHub
	content, err := githubpkg.FetchRawFile(site.GithubRepo, site.GithubBranch, path)
	if err != nil {
		// Try index.md if README.md doesn't exist
		if strings.HasSuffix(path, "README.md") {
			indexPath := strings.Replace(path, "README.md", "index.md", 1)
			content, err = githubpkg.FetchRawFile(site.GithubRepo, site.GithubBranch, indexPath)
			if err != nil {
				http.Error(w, fmt.Sprintf("File not found: %s", err), http.StatusNotFound)
				return nil
			}
		} else {
			http.Error(w, fmt.Sprintf("File not found: %s", err), http.StatusNotFound)
			return nil
		}
	}

	// Render markdown to HTML
	htmlContent, err := markdown.RenderMarkdown(content)
	if err != nil {
		return fmt.Errorf("failed to render markdown: %w", err)
	}

	return pages.PublicSite(w, r, site, htmlContent)
}
