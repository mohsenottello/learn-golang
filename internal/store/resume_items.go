package store

// Resume sections. Every row is owned by a user and ordered by position; every
// write is scoped by user_id so one user can never touch another's rows.
// Fields ending in En are translations; empty means "use the default language".

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrInvalid = errors.New("invalid value")

type Software struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

type Language struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	NameEn   string `json:"name_en"`
	Level    string `json:"level"`
	LevelEn  string `json:"level_en"`
	Position int    `json:"position"`
}

type Education struct {
	ID       int64  `json:"id"`
	Degree   string `json:"degree"`
	DegreeEn string `json:"degree_en"`
	School   string `json:"school"`
	SchoolEn string `json:"school_en"`
	DateFrom string `json:"date_from"` // "YYYY-MM" or ""
	DateTo   string `json:"date_to"`   // "YYYY-MM" or "" (= present)
	Position int    `json:"position"`
}

type Bullet struct {
	ID       int64  `json:"id"`
	Text     string `json:"text"`
	TextEn   string `json:"text_en"`
	Position int    `json:"position"`
}

type Experience struct {
	ID         int64    `json:"id"`
	Company    string   `json:"company"`
	Link       string   `json:"link"`
	Role       string   `json:"role"`
	RoleEn     string   `json:"role_en"`
	Subtitle   string   `json:"subtitle"`
	SubtitleEn string   `json:"subtitle_en"`
	DateFrom   string   `json:"date_from"`
	DateTo     string   `json:"date_to"`
	Location   string   `json:"location"`
	LocationEn string   `json:"location_en"`
	Bullets    []Bullet `json:"bullets"`
	Position   int      `json:"position"`
}

// ---- helpers ----

func monthToDB(s string) (any, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01", s)
	if err != nil {
		return nil, fmt.Errorf("%w: bad month %q", ErrInvalid, s)
	}
	return t, nil
}

func monthFromDB(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format("2006-01")
}

var sectionTables = map[string]string{
	"software":   "resume_software",
	"languages":  "resume_languages",
	"education":  "resume_education",
	"experience": "resume_experience",
}

func (s *Store) nextPosition(ctx context.Context, table, where string, arg int64) (int, error) {
	var n sql.NullInt64
	err := s.db.QueryRowContext(ctx,
		`SELECT max(position) FROM `+table+` WHERE `+where+` = $1`, arg).Scan(&n)
	return int(n.Int64) + 1, err
}

func rowsAffectedOrNotFound(res sql.Result) error {
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) deleteItem(ctx context.Context, table string, userID, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM `+table+` WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	return rowsAffectedOrNotFound(res)
}

// Reorder sets positions for a section in the given id order (ids not owned by the user are ignored).
func (s *Store) Reorder(ctx context.Context, section string, userID int64, ids []int64) error {
	table, ok := sectionTables[section]
	if !ok {
		return ErrInvalid
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, id := range ids {
		if _, err := tx.ExecContext(ctx,
			`UPDATE `+table+` SET position = $1 WHERE id = $2 AND user_id = $3`, i, id, userID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// ---------- software ----------

func (s *Store) ListSoftware(ctx context.Context, userID int64) ([]Software, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, position FROM resume_software WHERE user_id = $1 ORDER BY position, id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Software{}
	for rows.Next() {
		var x Software
		if err := rows.Scan(&x.ID, &x.Name, &x.Position); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (s *Store) CreateSoftware(ctx context.Context, userID int64, x Software) (*Software, error) {
	pos, err := s.nextPosition(ctx, "resume_software", "user_id", userID)
	if err != nil {
		return nil, err
	}
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO resume_software (user_id, name, position) VALUES ($1, $2, $3)
		RETURNING id, position`, userID, x.Name, pos).Scan(&x.ID, &x.Position)
	return &x, err
}

func (s *Store) UpdateSoftware(ctx context.Context, userID int64, x Software) (*Software, error) {
	res, err := s.db.ExecContext(ctx,
		`UPDATE resume_software SET name = $3 WHERE id = $1 AND user_id = $2`, x.ID, userID, x.Name)
	if err != nil {
		return nil, err
	}
	return &x, rowsAffectedOrNotFound(res)
}

func (s *Store) DeleteSoftware(ctx context.Context, userID, id int64) error {
	return s.deleteItem(ctx, "resume_software", userID, id)
}

// ---------- languages ----------

func (s *Store) ListLanguages(ctx context.Context, userID int64) ([]Language, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, name_en, level, level_en, position
		FROM resume_languages WHERE user_id = $1 ORDER BY position, id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Language{}
	for rows.Next() {
		var x Language
		if err := rows.Scan(&x.ID, &x.Name, &x.NameEn, &x.Level, &x.LevelEn, &x.Position); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (s *Store) CreateLanguage(ctx context.Context, userID int64, x Language) (*Language, error) {
	pos, err := s.nextPosition(ctx, "resume_languages", "user_id", userID)
	if err != nil {
		return nil, err
	}
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO resume_languages (user_id, name, name_en, level, level_en, position)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, position`,
		userID, x.Name, x.NameEn, x.Level, x.LevelEn, pos).Scan(&x.ID, &x.Position)
	return &x, err
}

func (s *Store) UpdateLanguage(ctx context.Context, userID int64, x Language) (*Language, error) {
	res, err := s.db.ExecContext(ctx, `
		UPDATE resume_languages SET name=$3, name_en=$4, level=$5, level_en=$6
		WHERE id = $1 AND user_id = $2`, x.ID, userID, x.Name, x.NameEn, x.Level, x.LevelEn)
	if err != nil {
		return nil, err
	}
	return &x, rowsAffectedOrNotFound(res)
}

func (s *Store) DeleteLanguage(ctx context.Context, userID, id int64) error {
	return s.deleteItem(ctx, "resume_languages", userID, id)
}

// ---------- education ----------

func (s *Store) ListEducation(ctx context.Context, userID int64) ([]Education, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, degree, degree_en, school, school_en, date_from, date_to, position
		FROM resume_education WHERE user_id = $1 ORDER BY position, id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Education{}
	for rows.Next() {
		var x Education
		var from, to sql.NullTime
		if err := rows.Scan(&x.ID, &x.Degree, &x.DegreeEn, &x.School, &x.SchoolEn, &from, &to, &x.Position); err != nil {
			return nil, err
		}
		x.DateFrom, x.DateTo = monthFromDB(from), monthFromDB(to)
		out = append(out, x)
	}
	return out, rows.Err()
}

func (s *Store) CreateEducation(ctx context.Context, userID int64, x Education) (*Education, error) {
	from, err := monthToDB(x.DateFrom)
	if err != nil {
		return nil, err
	}
	to, err := monthToDB(x.DateTo)
	if err != nil {
		return nil, err
	}
	pos, err := s.nextPosition(ctx, "resume_education", "user_id", userID)
	if err != nil {
		return nil, err
	}
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO resume_education (user_id, degree, degree_en, school, school_en, date_from, date_to, position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, position`,
		userID, x.Degree, x.DegreeEn, x.School, x.SchoolEn, from, to, pos).Scan(&x.ID, &x.Position)
	return &x, err
}

func (s *Store) UpdateEducation(ctx context.Context, userID int64, x Education) (*Education, error) {
	from, err := monthToDB(x.DateFrom)
	if err != nil {
		return nil, err
	}
	to, err := monthToDB(x.DateTo)
	if err != nil {
		return nil, err
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE resume_education SET degree=$3, degree_en=$4, school=$5, school_en=$6, date_from=$7, date_to=$8
		WHERE id = $1 AND user_id = $2`, x.ID, userID, x.Degree, x.DegreeEn, x.School, x.SchoolEn, from, to)
	if err != nil {
		return nil, err
	}
	return &x, rowsAffectedOrNotFound(res)
}

func (s *Store) DeleteEducation(ctx context.Context, userID, id int64) error {
	return s.deleteItem(ctx, "resume_education", userID, id)
}

// ---------- experience ----------

func (s *Store) ListExperience(ctx context.Context, userID int64) ([]Experience, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, company, link, role, role_en, subtitle, subtitle_en, date_from, date_to, location, location_en, position
		FROM resume_experience WHERE user_id = $1 ORDER BY position, id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Experience{}
	index := map[int64]int{}
	for rows.Next() {
		var x Experience
		var from, to sql.NullTime
		if err := rows.Scan(&x.ID, &x.Company, &x.Link, &x.Role, &x.RoleEn, &x.Subtitle, &x.SubtitleEn,
			&from, &to, &x.Location, &x.LocationEn, &x.Position); err != nil {
			return nil, err
		}
		x.DateFrom, x.DateTo = monthFromDB(from), monthFromDB(to)
		x.Bullets = []Bullet{}
		index[x.ID] = len(out)
		out = append(out, x)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	brows, err := s.db.QueryContext(ctx, `
		SELECT b.id, b.experience_id, b.text, b.text_en, b.position
		FROM resume_experience_bullets b
		JOIN resume_experience e ON e.id = b.experience_id
		WHERE e.user_id = $1 ORDER BY b.experience_id, b.position, b.id`, userID)
	if err != nil {
		return nil, err
	}
	defer brows.Close()
	for brows.Next() {
		var b Bullet
		var expID int64
		if err := brows.Scan(&b.ID, &expID, &b.Text, &b.TextEn, &b.Position); err != nil {
			return nil, err
		}
		if i, ok := index[expID]; ok {
			out[i].Bullets = append(out[i].Bullets, b)
		}
	}
	return out, brows.Err()
}

func (s *Store) CreateExperience(ctx context.Context, userID int64, x Experience) (*Experience, error) {
	from, err := monthToDB(x.DateFrom)
	if err != nil {
		return nil, err
	}
	to, err := monthToDB(x.DateTo)
	if err != nil {
		return nil, err
	}
	pos, err := s.nextPosition(ctx, "resume_experience", "user_id", userID)
	if err != nil {
		return nil, err
	}
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO resume_experience (user_id, company, link, role, role_en, subtitle, subtitle_en, date_from, date_to, location, location_en, position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, position`,
		userID, x.Company, x.Link, x.Role, x.RoleEn, x.Subtitle, x.SubtitleEn, from, to, x.Location, x.LocationEn, pos,
	).Scan(&x.ID, &x.Position)
	x.Bullets = []Bullet{}
	return &x, err
}

func (s *Store) UpdateExperience(ctx context.Context, userID int64, x Experience) (*Experience, error) {
	from, err := monthToDB(x.DateFrom)
	if err != nil {
		return nil, err
	}
	to, err := monthToDB(x.DateTo)
	if err != nil {
		return nil, err
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE resume_experience
		SET company=$3, link=$4, role=$5, role_en=$6, subtitle=$7, subtitle_en=$8, date_from=$9, date_to=$10, location=$11, location_en=$12
		WHERE id = $1 AND user_id = $2`,
		x.ID, userID, x.Company, x.Link, x.Role, x.RoleEn, x.Subtitle, x.SubtitleEn, from, to, x.Location, x.LocationEn)
	if err != nil {
		return nil, err
	}
	return &x, rowsAffectedOrNotFound(res)
}

func (s *Store) DeleteExperience(ctx context.Context, userID, id int64) error {
	return s.deleteItem(ctx, "resume_experience", userID, id)
}

// ---------- bullets (owned via their experience) ----------

// ownsExperience returns ErrNotFound unless the experience belongs to the user.
func (s *Store) ownsExperience(ctx context.Context, userID, expID int64) error {
	var n int
	err := s.db.QueryRowContext(ctx,
		`SELECT count(*) FROM resume_experience WHERE id = $1 AND user_id = $2`, expID, userID).Scan(&n)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) CreateBullet(ctx context.Context, userID, expID int64, b Bullet) (*Bullet, error) {
	if err := s.ownsExperience(ctx, userID, expID); err != nil {
		return nil, err
	}
	pos, err := s.nextPosition(ctx, "resume_experience_bullets", "experience_id", expID)
	if err != nil {
		return nil, err
	}
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO resume_experience_bullets (experience_id, text, text_en, position)
		VALUES ($1, $2, $3, $4) RETURNING id, position`, expID, b.Text, b.TextEn, pos).Scan(&b.ID, &b.Position)
	return &b, err
}

func (s *Store) UpdateBullet(ctx context.Context, userID, expID int64, b Bullet) (*Bullet, error) {
	if err := s.ownsExperience(ctx, userID, expID); err != nil {
		return nil, err
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE resume_experience_bullets SET text=$3, text_en=$4
		WHERE id = $1 AND experience_id = $2`, b.ID, expID, b.Text, b.TextEn)
	if err != nil {
		return nil, err
	}
	return &b, rowsAffectedOrNotFound(res)
}

func (s *Store) DeleteBullet(ctx context.Context, userID, expID, id int64) error {
	if err := s.ownsExperience(ctx, userID, expID); err != nil {
		return err
	}
	res, err := s.db.ExecContext(ctx,
		`DELETE FROM resume_experience_bullets WHERE id = $1 AND experience_id = $2`, id, expID)
	if err != nil {
		return err
	}
	return rowsAffectedOrNotFound(res)
}

func (s *Store) ReorderBullets(ctx context.Context, userID, expID int64, ids []int64) error {
	if err := s.ownsExperience(ctx, userID, expID); err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, id := range ids {
		if _, err := tx.ExecContext(ctx,
			`UPDATE resume_experience_bullets SET position = $1 WHERE id = $2 AND experience_id = $3`, i, id, expID); err != nil {
			return err
		}
	}
	return tx.Commit()
}
