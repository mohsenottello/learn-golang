package store

import (
	"context"
	"database/sql"
	"time"
)

// Resume is the header: slug, visibility, contact block, about text.
// Fields with an En twin are translatable; En empty means "fall back to the default".
type Resume struct {
	UserID      int64     `json:"user_id"`
	Slug        string    `json:"slug"`
	Published   bool      `json:"published"`
	DefaultLang string    `json:"default_lang"` // "de" | "en"
	Title       string    `json:"title"`
	TitleEn     string    `json:"title_en"`
	Subtitle    string    `json:"subtitle"`
	SubtitleEn  string    `json:"subtitle_en"`
	About       string    `json:"about"`
	AboutEn     string    `json:"about_en"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Location    string    `json:"location"`
	LinkedIn    string    `json:"linkedin"`
	HasPhoto    bool      `json:"has_photo"`
	SideBg			string    `json:"side_bg"`
	SideText		string    `json:"side_text"`
	MainBg			string    `json:"main_bg"`
	MainText		string    `json:"main_text"`
	Accent			string    `json:"accent"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdminResume struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Published bool      `json:"published"`
	UpdatedAt time.Time `json:"updated_at"`
}

const resumeCols = `user_id, slug, published, default_lang, title, title_en, subtitle, subtitle_en,
	about, about_en, email, phone, location, linkedin, photo IS NOT NULL, updated_at,
	side_bg, side_text, main_bg, main_text, accent`

func scanResume(row interface{ Scan(...any) error }) (*Resume, error) {
	var r Resume
	err := row.Scan(&r.UserID, &r.Slug, &r.Published, &r.DefaultLang, &r.Title, &r.TitleEn,
		&r.Subtitle, &r.SubtitleEn, &r.About, &r.AboutEn,
		&r.Email, &r.Phone, &r.Location, &r.LinkedIn, &r.HasPhoto, &r.UpdatedAt,
		&r.SideBg, &r.SideText, &r.MainBg, &r.MainText, &r.Accent)
	if err != nil {
		return nil, translate(err)
	}
	return &r, nil
}

func (s *Store) ResumeByUser(ctx context.Context, userID int64) (*Resume, error) {
	return scanResume(s.db.QueryRowContext(ctx,
		`SELECT `+resumeCols+` FROM resumes WHERE user_id = $1`, userID))
}

// ResumeBySlug returns a published resume, or ErrNotFound.
func (s *Store) ResumeBySlug(ctx context.Context, slug string) (*Resume, error) {
	return scanResume(s.db.QueryRowContext(ctx,
		`SELECT `+resumeCols+` FROM resumes WHERE slug = $1 AND published`, slug))
}

func (s *Store) CreateResume(ctx context.Context, userID int64, slug string) (*Resume, error) {
	return scanResume(s.db.QueryRowContext(ctx, `
		INSERT INTO resumes (user_id, slug) VALUES ($1, $2)
		RETURNING `+resumeCols, userID, slug))
}

func (s *Store) UpdateResume(ctx context.Context, r *Resume) (*Resume, error) {
	return scanResume(s.db.QueryRowContext(ctx, `
		UPDATE resumes SET slug=$2, published=$3, default_lang=$4,
			title=$5, title_en=$6, subtitle=$7, subtitle_en=$8, about=$9, about_en=$10,
			email=$11, phone=$12, location=$13, linkedin=$14,
			side_bg=$15, side_text=$16, main_bg=$17, main_text=$18, accent=$19,
			updated_at=now()
		WHERE user_id = $1
		RETURNING `+resumeCols,
		r.UserID, r.Slug, r.Published, r.DefaultLang,
		r.Title, r.TitleEn, r.Subtitle, r.SubtitleEn, r.About, r.AboutEn,
		r.Email, r.Phone, r.Location, r.LinkedIn, r.SideBg, r.SideText, r.MainBg, r.MainText, r.Accent))
}

// --- photo ---

func (s *Store) SetResumePhoto(ctx context.Context, userID int64, data []byte, contentType string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE resumes SET photo=$2, photo_type=$3, updated_at=now() WHERE user_id=$1`,
		userID, data, contentType)
	return err
}

func (s *Store) ResumePhotoBySlug(ctx context.Context, slug string) ([]byte, string, error) {
	var data []byte
	var ct sql.NullString
	err := s.db.QueryRowContext(ctx,
		`SELECT photo, photo_type FROM resumes WHERE slug = $1 AND published AND photo IS NOT NULL`, slug,
	).Scan(&data, &ct)
	if err != nil {
		return nil, "", translate(err)
	}
	return data, ct.String, nil
}

func (s *Store) ResumePhotoByUser(ctx context.Context, userID int64) ([]byte, string, error) {
	var data []byte
	var ct sql.NullString
	err := s.db.QueryRowContext(ctx,
		`SELECT photo, photo_type FROM resumes WHERE user_id = $1 AND photo IS NOT NULL`, userID,
	).Scan(&data, &ct)
	if err != nil {
		return nil, "", translate(err)
	}
	return data, ct.String, nil
}


func (s *Store) ListResumesForAdmin(ctx context.Context) ([]AdminResume, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT r.user_id, u.email, u.name, r.slug, r.published, r.updated_at
		FROM resumes r JOIN users u ON u.id = r.user_id ORDER BY r.updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AdminResume{}
	for rows.Next() {
		var x AdminResume
		if err := rows.Scan(&x.UserID, &x.Email, &x.Name, &x.Slug, &x.Published, &x.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (s *Store) SetResumePublished(ctx context.Context, userID int64, published bool) error {
	res, err := s.db.ExecContext(ctx, `UPDATE resumes SET published = $2 WHERE user_id = $1`, userID, published)
	if err != nil {
		return err
	}
	return rowsAffectedOrNotFound(res)
}
