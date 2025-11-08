package pages

import (
	"fmt"
	"net/http"

	"github.com/frenchsoftware/libhtml/attr"
	"github.com/frenchsoftware/libhtml/html"
	"github.com/frenchsoftware/libvalidator/validator"
	"github.com/hyperstitieux/template/auth"
	"github.com/hyperstitieux/template/database/models"
	"github.com/hyperstitieux/template/views"
	"github.com/hyperstitieux/template/views/components/ui"
	"github.com/hyperstitieux/template/views/layouts"
)

func Home(w http.ResponseWriter, r *http.Request) error {
	// Get authenticated user from context (nil if not authenticated)
	user := auth.GetCurrentUser(r)

	// Build page
	page := layouts.Base(user, r, "Internet Publishing") // Page content goes here

	// Render page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return page.Render(w)
}

func Settings(w http.ResponseWriter, r *http.Request) error {
	return SettingsWithErrors(w, r, nil)
}

func SettingsWithErrors(w http.ResponseWriter, r *http.Request, errs validator.ValidationErrors) error {
	// Get authenticated user from context (required for settings page)
	user := views.GetUser(r)
	if user == nil {
		// Redirect to OAuth with return URL
		http.Redirect(w, r, "/auth/google?redirect=/settings", http.StatusTemporaryRedirect)
		return nil
	}

	// Build page
	page := layouts.Base(user, r, "Settings - Internet Publishing",
		// Main content container
		html.Div(
			attr.Class("max-w-4xl mx-auto px-8 py-8"),

			// Page header
			html.Div(
				attr.Class("mb-8"),
				html.H1(
					attr.Class("text-3xl font-semibold mb-2"),
					html.Text("Settings"),
				),
				html.P(
					attr.Class("text-muted-foreground"),
					html.Text("Manage your account settings and preferences"),
				),
			),

			// Settings cards container
			html.Div(
				attr.Class("flex flex-col gap-6"),

				// General settings form
				html.Form(
					attr.Id("profile-form"),
					attr.Action("/settings/update-profile"),
					attr.Method("POST"),

					ui.Card(
						ui.CardHeader(ui.CardHeaderProps{
							Title:       "General",
							Description: "Update your personal information",
						}),
						ui.CardSection(
							// Email field (readonly)
							html.Div(
								attr.Class("flex flex-col gap-2"),
								html.Label(
									attr.For("email"),
									attr.Class("text-sm font-medium"),
									html.Text("Email"),
								),
								html.Input(
									attr.Type("email"),
									attr.Id("email"),
									attr.Name("email"),
									attr.Value(user.Email),
									attr.Readonly("true"),
									attr.Class("input bg-muted text-muted-foreground cursor-not-allowed"),
								),
								html.P(
									attr.Class("text-xs text-muted-foreground"),
									html.Text("Your email address cannot be changed"),
								),
							),

							// Name field
							html.Div(
								attr.Class("flex flex-col gap-2"),
								html.Label(
									attr.For("name"),
									attr.Class("text-sm font-medium"),
									html.Text("Name"),
								),
								html.Input(
									attr.Type("text"),
									attr.Id("name"),
									attr.Name("name"),
									attr.Value(user.Name),
									attr.Required("true"),
									attr.Maxlength("80"),
									attr.ClassIfElse(errs != nil && errs.Has("name"), "input border-destructive focus:ring-destructive", "input"),
								),
								html.If(errs != nil && errs.Has("name"),
									html.P(
										attr.Class("text-xs text-destructive"),
										html.Text(errs.Get("name")),
									),
								),
							),
						),
						ui.CardFooter(
							html.Button(
								attr.Type("submit"),
								attr.Class("btn-primary"),
								html.Text("Save changes"),
							),
						),
					),
				),

				// Danger zone card
				ui.Card(
					ui.CardHeader(ui.CardHeaderProps{
						Title:       "Danger Zone",
						Description: "Irreversible and destructive actions",
					}),
					ui.CardSection(
						html.Div(
							attr.Class("border border-destructive rounded-lg p-4"),
							html.Div(
								attr.Class("flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4"),
								html.Div(
									html.H3(
										attr.Class("text-sm font-medium"),
										html.Text("Delete Account"),
									),
									html.P(
										attr.Class("text-sm text-muted-foreground"),
										html.Text("Permanently delete your account and all associated data"),
									),
								),
								html.Button(
									attr.Type("button"),
									attr.Class("btn-destructive"),
									html.Attr("onclick", "document.getElementById('delete-account-dialog').showModal()"),
									html.Text("Delete Account"),
								),
							),
						),
					),
				),
			),
		),

		// Delete account confirmation dialog with form in footer
		ui.AlertDialog(ui.AlertDialogProps{
			ID:          "delete-account-dialog",
			Title:       "Are you absolutely sure?",
			Description: "This action cannot be undone. This will permanently delete your account and remove all your data from our servers.",
			Footer: html.Div(
				attr.Class("flex gap-2"),
				html.Button(
					attr.Type("button"),
					attr.Class("btn-outline"),
					html.Attr("onclick", "this.closest('dialog').close()"),
					html.Text("Cancel"),
				),
				html.Form(
					attr.Action("/settings/delete-account"),
					attr.Method("POST"),
					attr.Class("inline"),
					html.Button(
						attr.Type("submit"),
						attr.Class("btn-destructive"),
						html.Text("Delete Account"),
					),
				),
			),
		}),
	)

	// Render page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return page.Render(w)
}

func Sites(w http.ResponseWriter, r *http.Request, sites []*models.Site) error {
	user := views.GetUser(r)

	// Build page
	page := layouts.Base(user, r, "My Sites - Internet Publishing",
		html.Div(
			attr.Class("max-w-6xl mx-auto px-8 py-8"),

			// Page header
			html.Div(
				attr.Class("flex items-center justify-between mb-8"),
				html.Div(
					html.H1(
						attr.Class("text-3xl font-semibold mb-2"),
						html.Text("My Sites"),
					),
					html.P(
						attr.Class("text-muted-foreground"),
						html.Text("Manage your published documentation sites"),
					),
				),
				html.A(
					attr.Href("/sites/new"),
					attr.Class("btn-primary"),
					html.Text("+ New Site"),
				),
			),

			// Sites grid or empty state
			html.If(len(sites) == 0,
				// Empty state
				html.Div(
					attr.Class("text-center py-16"),
					html.H2(
						attr.Class("text-xl font-medium mb-2"),
						html.Text("No sites yet"),
					),
					html.P(
						attr.Class("text-muted-foreground mb-6"),
						html.Text("Get started by creating your first documentation site"),
					),
					html.A(
						attr.Href("/sites/new"),
						attr.Class("btn-primary"),
						html.Text("Create Your First Site"),
					),
				),
			),

			// Sites grid
			html.If(len(sites) > 0,
				func() html.Node {
					var siteCards []any
					siteCards = append(siteCards, attr.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"))
					for _, site := range sites {
						siteCards = append(siteCards, ui.Card(
							ui.CardHeader(ui.CardHeaderProps{
								Title:       site.Slug,
								Description: site.GithubRepo,
							}),
							ui.CardSection(
								html.Div(
									attr.Class("flex flex-col gap-2 text-sm"),
									html.Div(
										attr.Class("flex items-center gap-2 text-muted-foreground"),
										html.Span(
											html.Text("Branch: "),
										),
										html.Span(
											attr.Class("font-mono text-xs bg-muted px-2 py-1 rounded"),
											html.Text(site.GithubBranch),
										),
									),
									html.If(site.Subdirectory != "",
										html.Div(
											attr.Class("flex items-center gap-2 text-muted-foreground"),
											html.Span(
												html.Text("Path: "),
											),
											html.Span(
												attr.Class("font-mono text-xs bg-muted px-2 py-1 rounded"),
												html.Text(site.Subdirectory),
											),
										),
									),
									html.Div(
										attr.Class("text-muted-foreground"),
										html.Text(fmt.Sprintf("Created: %s", site.CreatedAt.Format("Jan 2, 2006"))),
									),
								),
							),
							ui.CardFooter(
								html.Div(
									attr.Class("flex gap-2"),
									html.A(
										attr.Href(fmt.Sprintf("https://%s.internetpublishing.co", site.Slug)),
										attr.Target("_blank"),
										attr.Class("btn-primary text-sm"),
										html.Text("View Site →"),
									),
									html.Button(
										attr.Type("button"),
										attr.Class("btn-outline text-sm"),
										html.Attr("onclick", fmt.Sprintf("if(confirm('Delete %s?')) fetch('/sites/%d/delete', {method: 'POST'}).then(() => location.reload())", site.Slug, site.ID)),
										html.Text("Delete"),
									),
								),
							),
						))
					}
					return html.Div(siteCards...)
				}(),
			),
		),
	)

	// Render page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return page.Render(w)
}

func NewSite(w http.ResponseWriter, r *http.Request, errs validator.ValidationErrors) error {
	user := views.GetUser(r)

	// Build page
	page := layouts.Base(user, r, "New Site - Internet Publishing",
		html.Div(
			attr.Class("max-w-2xl mx-auto px-8 py-8"),

			// Page header
			html.Div(
				attr.Class("mb-8"),
				html.H1(
					attr.Class("text-3xl font-semibold mb-2"),
					html.Text("Create New Site"),
				),
				html.P(
					attr.Class("text-muted-foreground"),
					html.Text("Publish markdown documentation from your GitHub repository"),
				),
			),

			// Form
			html.Form(
				attr.Action("/sites/create"),
				attr.Method("POST"),

				ui.Card(
					ui.CardSection(
						// Slug field
						html.Div(
							attr.Class("flex flex-col gap-2"),
							html.Label(
								attr.For("slug"),
								attr.Class("text-sm font-medium"),
								html.Text("Site Slug"),
							),
							html.Div(
								attr.Class("flex items-center gap-2"),
								html.Input(
									attr.Type("text"),
									attr.Id("slug"),
									attr.Name("slug"),
									attr.Required("true"),
									attr.Pattern("[a-z0-9-]+"),
									attr.Placeholder("my-docs"),
									attr.ClassIfElse(errs != nil && errs.Has("slug"), "input border-destructive focus:ring-destructive flex-1", "input flex-1"),
								),
								html.Span(
									attr.Class("text-sm text-muted-foreground"),
									html.Text(".internetpublishing.co"),
								),
							),
							html.P(
								attr.Class("text-xs text-muted-foreground"),
								html.Text("Lowercase letters, numbers, and hyphens only"),
							),
							html.If(errs != nil && errs.Has("slug"),
								html.P(
									attr.Class("text-xs text-destructive"),
									html.Text(errs.Get("slug")),
								),
							),
						),

						// GitHub repo field
						html.Div(
							attr.Class("flex flex-col gap-2"),
							html.Label(
								attr.For("github_repo"),
								attr.Class("text-sm font-medium"),
								html.Text("GitHub Repository"),
							),
							html.Input(
								attr.Type("text"),
								attr.Id("github_repo"),
								attr.Name("github_repo"),
								attr.Required("true"),
								attr.Placeholder("username/repository"),
								attr.Pattern("[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+"),
								attr.ClassIfElse(errs != nil && errs.Has("github_repo"), "input border-destructive focus:ring-destructive", "input"),
							),
							html.P(
								attr.Class("text-xs text-muted-foreground"),
								html.Text("Must be a public repository (e.g., octocat/Hello-World)"),
							),
							html.If(errs != nil && errs.Has("github_repo"),
								html.P(
									attr.Class("text-xs text-destructive"),
									html.Text(errs.Get("github_repo")),
								),
							),
						),

						// Branch field
						html.Div(
							attr.Class("flex flex-col gap-2"),
							html.Label(
								attr.For("github_branch"),
								attr.Class("text-sm font-medium"),
								html.Text("Branch"),
							),
							html.Input(
								attr.Type("text"),
								attr.Id("github_branch"),
								attr.Name("github_branch"),
								attr.Value("main"),
								attr.Required("true"),
								attr.Class("input"),
							),
							html.P(
								attr.Class("text-xs text-muted-foreground"),
								html.Text("Branch to publish from (usually 'main' or 'master')"),
							),
						),

						// Subdirectory field (optional)
						html.Div(
							attr.Class("flex flex-col gap-2"),
							html.Label(
								attr.For("subdirectory"),
								attr.Class("text-sm font-medium"),
								html.Text("Subdirectory (optional)"),
							),
							html.Input(
								attr.Type("text"),
								attr.Id("subdirectory"),
								attr.Name("subdirectory"),
								attr.Placeholder("docs"),
								attr.Class("input"),
							),
							html.P(
								attr.Class("text-xs text-muted-foreground"),
								html.Text("Only publish files from this subdirectory (e.g., 'docs', 'content/posts')"),
							),
						),
					),

					ui.CardFooter(
						html.Div(
							attr.Class("flex gap-2"),
							html.A(
								attr.Href("/sites"),
								attr.Class("btn-outline"),
								html.Text("Cancel"),
							),
							html.Button(
								attr.Type("submit"),
								attr.Class("btn-primary"),
								html.Text("Create Site"),
							),
						),
					),
				),
			),
		),
	)

	// Render page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return page.Render(w)
}

func PublicSite(w http.ResponseWriter, r *http.Request, site *models.Site, htmlContent string) error {
	// Build simple public page
	page := html.Html(
		html.Head(
			html.Meta(attr.Charset("utf-8")),
			html.Meta(attr.Name("viewport"), attr.Content("width=device-width, initial-scale=1")),
			html.Title(html.Text(fmt.Sprintf("%s - %s", site.Slug, site.GithubRepo))),
			html.Link(attr.Rel("stylesheet"), attr.Href("https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css")),
			html.Style(
				html.Text(`
					body { padding: 2rem; }
					main { max-width: 800px; margin: 0 auto; }
					.header { border-bottom: 1px solid var(--contrast-lower); padding-bottom: 1rem; margin-bottom: 2rem; }
					.footer { margin-top: 3rem; padding-top: 2rem; border-top: 1px solid var(--contrast-lower); text-align: center; color: var(--muted-color); }
					pre { overflow-x: auto; }
					code { font-size: 0.9em; }
				`),
			),
		),
		html.Body(
			html.Header(
				attr.Class("header"),
				html.H1(html.Text(site.Slug)),
				html.P(
					html.Text("Published from "),
					html.A(
						attr.Href(fmt.Sprintf("https://github.com/%s", site.GithubRepo)),
						attr.Target("_blank"),
						html.Text(site.GithubRepo),
					),
				),
			),
			html.Main(
				html.Raw(htmlContent),
			),
			html.Footer(
				attr.Class("footer"),
				html.Small(
					html.Text("Powered by "),
					html.A(
						attr.Href("https://internetpublishing.co"),
						html.Text("Internet Publishing"),
					),
				),
			),
		),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return page.Render(w)
}
