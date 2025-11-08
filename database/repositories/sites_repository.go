package repositories

import (
	"database/sql"

	"github.com/hyperstitieux/template/database/models"
)

type SitesRepository interface {
	Create(userID int, slug, githubRepo, githubBranch, subdirectory string) (*models.Site, error)
	GetByID(id int) (*models.Site, error)
	GetBySlug(slug string) (*models.Site, error)
	GetByUserID(userID int) ([]*models.Site, error)
	Delete(id int) error
}

type sitesRepository struct {
	db *sql.DB
}

func NewSitesRepository(db *sql.DB) SitesRepository {
	return &sitesRepository{db: db}
}

func (r *sitesRepository) Create(userID int, slug, githubRepo, githubBranch, subdirectory string) (*models.Site, error) {
	query := `
		INSERT INTO sites (user_id, slug, github_repo, github_branch, subdirectory)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, userID, slug, githubRepo, githubBranch, subdirectory)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(int(id))
}

func (r *sitesRepository) GetByID(id int) (*models.Site, error) {
	query := `
		SELECT id, user_id, slug, github_repo, github_branch, subdirectory, created_at
		FROM sites
		WHERE id = ?
	`
	site := &models.Site{}
	err := r.db.QueryRow(query, id).Scan(
		&site.ID,
		&site.UserID,
		&site.Slug,
		&site.GithubRepo,
		&site.GithubBranch,
		&site.Subdirectory,
		&site.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *sitesRepository) GetBySlug(slug string) (*models.Site, error) {
	query := `
		SELECT id, user_id, slug, github_repo, github_branch, subdirectory, created_at
		FROM sites
		WHERE slug = ?
	`
	site := &models.Site{}
	err := r.db.QueryRow(query, slug).Scan(
		&site.ID,
		&site.UserID,
		&site.Slug,
		&site.GithubRepo,
		&site.GithubBranch,
		&site.Subdirectory,
		&site.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *sitesRepository) GetByUserID(userID int) ([]*models.Site, error) {
	query := `
		SELECT id, user_id, slug, github_repo, github_branch, subdirectory, created_at
		FROM sites
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sites := []*models.Site{}
	for rows.Next() {
		site := &models.Site{}
		err := rows.Scan(
			&site.ID,
			&site.UserID,
			&site.Slug,
			&site.GithubRepo,
			&site.GithubBranch,
			&site.Subdirectory,
			&site.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}

	return sites, rows.Err()
}

func (r *sitesRepository) Delete(id int) error {
	query := `DELETE FROM sites WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
