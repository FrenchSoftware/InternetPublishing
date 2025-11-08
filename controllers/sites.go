package controllers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/frenchsoftware/libvalidator/validator"
	"github.com/gorilla/mux"
	"github.com/hyperstitieux/template/database/repositories"
	"github.com/hyperstitieux/template/pages"
	"github.com/hyperstitieux/template/views"
)

type SitesController struct {
	sites *repositories.SitesRepository
}

func NewSitesController(sites *repositories.SitesRepository) *SitesController {
	return &SitesController{sites: sites}
}

var slugPattern = regexp.MustCompile(`^[a-z0-9-]+$`)
var repoPattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+$`)

func (c *SitesController) List(w http.ResponseWriter, r *http.Request) error {
	user := views.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/auth/google?redirect=/sites", http.StatusTemporaryRedirect)
		return nil
	}

	sites, err := (*c.sites).GetByUserID(int(user.ID))
	if err != nil {
		return err
	}

	return pages.Sites(w, r, sites)
}

func (c *SitesController) New(w http.ResponseWriter, r *http.Request) error {
	return pages.NewSite(w, r, nil)
}

func (c *SitesController) Create(w http.ResponseWriter, r *http.Request) error {
	user := views.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/auth/google?redirect=/sites/new", http.StatusTemporaryRedirect)
		return nil
	}

	// Validate
	v := validator.New(
		validator.Field("slug").Required().MinLength(1).MaxLength(50),
		validator.Field("github_repo").Required().MinLength(1),
		validator.Field("github_branch").Required().MinLength(1),
	)

	ok, errs := v.Validate(r)
	if !ok {
		return pages.NewSite(w, r, errs)
	}

	// Get form values
	slug := r.FormValue("slug")
	githubRepo := r.FormValue("github_repo")
	githubBranch := r.FormValue("github_branch")
	subdirectory := r.FormValue("subdirectory") // Optional

	// Additional validation
	additionalErrs := make(validator.ValidationErrors)

	if !slugPattern.MatchString(slug) {
		additionalErrs.Add("slug", "Slug must contain only lowercase letters, numbers, and hyphens")
	}

	if !repoPattern.MatchString(githubRepo) {
		additionalErrs.Add("github_repo", "Invalid repository format (use: username/repository)")
	}

	// Check if slug already exists
	existing, err := (*c.sites).GetBySlug(slug)
	if err != nil {
		return err
	}
	if existing != nil {
		additionalErrs.Add("slug", "This slug is already taken")
	}

	if !additionalErrs.IsEmpty() {
		return pages.NewSite(w, r, additionalErrs)
	}

	// Create site
	_, err = (*c.sites).Create(int(user.ID), slug, githubRepo, githubBranch, subdirectory)
	if err != nil {
		return err
	}

	// Redirect to sites list
	http.Redirect(w, r, "/sites", http.StatusSeeOther)
	return nil
}

func (c *SitesController) Delete(w http.ResponseWriter, r *http.Request) error {
	user := views.GetUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	// Get site ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid site ID", http.StatusBadRequest)
		return nil
	}

	// Get site to verify ownership
	site, err := (*c.sites).GetByID(id)
	if err != nil {
		return err
	}
	if site == nil {
		http.Error(w, "Site not found", http.StatusNotFound)
		return nil
	}
	if site.UserID != int(user.ID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return nil
	}

	// Delete site
	if err := (*c.sites).Delete(id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
