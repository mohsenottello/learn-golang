package httpapi

// User's own resume: /api/me/resume/...

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/mohsenottello/golang/internal/store"
)

var slugRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

// Paths the SPA uses; a user may not claim them as a slug.
var reservedSlugs = map[string]bool{
	"login": true, "register": true, "verify": true, "account": true, "admin": true,
	"forgot-password": true, "reset-password": true, "api": true, "healthz": true,
	"assets": true, "static": true, "index.html": true,
}

func validSlug(s string) bool {
	return len(s) >= 3 && len(s) <= 40 && slugRe.MatchString(s) && !reservedSlugs[s]
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.Trim(nonSlug.ReplaceAllString(s, "-"), "-")
	if len(s) > 40 {
		s = strings.Trim(s[:40], "-")
	}
	if len(s) < 3 {
		s = "user-" + s
	}
	return s
}

var validLangs = map[string]bool{"de": true, "en": true}

// FullResume is what the editor and the public page both consume.
type FullResume struct {
	Resume     *store.Resume      `json:"resume"`
	Name       string             `json:"name"`
	Lang       string             `json:"lang,omitempty"` // set on public (localized) responses
	Software   []store.Software   `json:"software"`
	Languages  []store.Language   `json:"languages"`
	Education  []store.Education  `json:"education"`
	Experience []store.Experience `json:"experience"`
}

func (h *handlers) loadFull(r *http.Request, res *store.Resume, name string) (*FullResume, error) {
	ctx := r.Context()
	sw, err := h.Store.ListSoftware(ctx, res.UserID)
	if err != nil {
		return nil, err
	}
	la, err := h.Store.ListLanguages(ctx, res.UserID)
	if err != nil {
		return nil, err
	}
	ed, err := h.Store.ListEducation(ctx, res.UserID)
	if err != nil {
		return nil, err
	}
	ex, err := h.Store.ListExperience(ctx, res.UserID)
	if err != nil {
		return nil, err
	}
	return &FullResume{Resume: res, Name: name, Software: sw, Languages: la, Education: ed, Experience: ex}, nil
}

// myResume returns the user's resume, creating an empty one with a slug from their name on first use.
func (h *handlers) myResume(r *http.Request) (*store.Resume, error) {
	u := currentUser(r)
	res, err := h.Store.ResumeByUser(r.Context(), u.ID)
	if !errors.Is(err, store.ErrNotFound) {
		return res, err
	}
	base := slugify(u.Name)
	for i := 0; ; i++ {
		slug := base
		if i > 0 {
			slug = base + "-" + strconv.Itoa(i+1)
		}
		if reservedSlugs[slug] {
			continue
		}
		res, err = h.Store.CreateResume(r.Context(), u.ID, slug)
		if !errors.Is(err, store.ErrDuplicate) {
			return res, err
		}
	}
}

func (h *handlers) ensureResume(w http.ResponseWriter, r *http.Request) bool {
	if _, err := h.myResume(r); err != nil {
		serverError(w, r, err)
		return false
	}
	return true
}

func (h *handlers) getMyResume(w http.ResponseWriter, r *http.Request) {
	res, err := h.myResume(r)
	if err != nil {
		serverError(w, r, err)
		return
	}
	full, err := h.loadFull(r, res, currentUser(r).Name)
	if err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, full)
}

func (h *handlers) updateMyResume(w http.ResponseWriter, r *http.Request) {
	var in store.Resume
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	in.Slug = strings.ToLower(strings.TrimSpace(in.Slug))
	if !validSlug(in.Slug) {
		writeError(w, http.StatusBadRequest, "Slug must be 3–40 characters: lowercase letters, digits and hyphens, and not a reserved word.")
		return
	}
	if !validLangs[in.DefaultLang] {
		in.DefaultLang = "de"
	}
	var hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

	for _, c := range []string{in.SideBg, in.SideText, in.MainBg, in.MainText, in.Accent} {
			if !hexColor.MatchString(c) {
					writeError(w, http.StatusBadRequest, "Colours must be hex like #7f7a72.")
					return
			}
	}
	in.UserID = currentUser(r).ID
	if !h.ensureResume(w, r) {
		return
	}
	res, err := h.Store.UpdateResume(r.Context(), &in)
	switch {
	case errors.Is(err, store.ErrDuplicate):
		writeError(w, http.StatusConflict, "That slug is already taken.")
	case err != nil:
		serverError(w, r, err)
	default:
		writeJSON(w, http.StatusOK, res)
	}
}

// --- photo ---

const maxPhotoBytes = 1 << 20

func (h *handlers) uploadMyPhoto(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxPhotoBytes+4096)
	if err := r.ParseMultipartForm(maxPhotoBytes); err != nil {
		writeError(w, http.StatusBadRequest, "Image must be smaller than 1 MB.")
		return
	}
	f, _, err := r.FormFile("photo")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing 'photo' field")
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read file")
		return
	}
	ct := http.DetectContentType(data)
	if ct != "image/jpeg" && ct != "image/png" && ct != "image/webp" {
		writeError(w, http.StatusBadRequest, "Only JPEG, PNG or WebP images are allowed.")
		return
	}
	if !h.ensureResume(w, r) {
		return
	}
	if err := h.Store.SetResumePhoto(r.Context(), currentUser(r).ID, data, ct); err != nil {
		serverError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) myPhoto(w http.ResponseWriter, r *http.Request) {
	data, ct, err := h.Store.ResumePhotoByUser(r.Context(), currentUser(r).ID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Cache-Control", "private, max-age=60")
	w.Write(data)
}

// --- section items ---

func pathID(r *http.Request, name string) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue(name), 10, 64)
	return id, err == nil
}

func (h *handlers) writeItemResult(w http.ResponseWriter, r *http.Request, v any, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, store.ErrInvalid):
		writeError(w, http.StatusBadRequest, "Dates must be in YYYY-MM format.")
	case err != nil:
		serverError(w, r, err)
	case v == nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		writeJSON(w, http.StatusOK, v)
	}
}

func (h *handlers) reorderSection(w http.ResponseWriter, r *http.Request) {
	var in struct {
		IDs []int64 `json:"ids"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	err := h.Store.Reorder(r.Context(), r.PathValue("section"), currentUser(r).ID, in.IDs)
	if errors.Is(err, store.ErrInvalid) {
		writeError(w, http.StatusBadRequest, "unknown section")
		return
	}
	h.writeItemResult(w, r, nil, err)
}

// software
func (h *handlers) createSoftware(w http.ResponseWriter, r *http.Request) {
	var in store.Software
	if err := readJSON(w, r, &in); err != nil || strings.TrimSpace(in.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if !h.ensureResume(w, r) {
		return
	}
	v, err := h.Store.CreateSoftware(r.Context(), currentUser(r).ID, in)
	h.writeItemResult(w, r, v, err)
}
func (h *handlers) updateSoftware(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	var in store.Software
	if err := readJSON(w, r, &in); !ok || err != nil || strings.TrimSpace(in.Name) == "" {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	in.ID = id
	v, err := h.Store.UpdateSoftware(r.Context(), currentUser(r).ID, in)
	h.writeItemResult(w, r, v, err)
}
func (h *handlers) deleteSoftware(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	h.writeItemResult(w, r, nil, h.Store.DeleteSoftware(r.Context(), currentUser(r).ID, id))
}

// languages
func (h *handlers) createLanguage(w http.ResponseWriter, r *http.Request) {
	var in store.Language
	if err := readJSON(w, r, &in); err != nil || strings.TrimSpace(in.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if !h.ensureResume(w, r) {
		return
	}
	v, err := h.Store.CreateLanguage(r.Context(), currentUser(r).ID, in)
	h.writeItemResult(w, r, v, err)
}
func (h *handlers) updateLanguage(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	var in store.Language
	if err := readJSON(w, r, &in); !ok || err != nil || strings.TrimSpace(in.Name) == "" {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	in.ID = id
	v, err := h.Store.UpdateLanguage(r.Context(), currentUser(r).ID, in)
	h.writeItemResult(w, r, v, err)
}
func (h *handlers) deleteLanguage(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	h.writeItemResult(w, r, nil, h.Store.DeleteLanguage(r.Context(), currentUser(r).ID, id))
}

// education
func (h *handlers) createEducation(w http.ResponseWriter, r *http.Request) {
	var in store.Education
	if err := readJSON(w, r, &in); err != nil || strings.TrimSpace(in.Degree) == "" {
		writeError(w, http.StatusBadRequest, "degree is required")
		return
	}
	if !h.ensureResume(w, r) {
		return
	}
	v, err := h.Store.CreateEducation(r.Context(), currentUser(r).ID, in)
	h.writeItemResult(w, r, v, err)
}
func (h *handlers) updateEducation(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	var in store.Education
	if err := readJSON(w, r, &in); !ok || err != nil || strings.TrimSpace(in.Degree) == "" {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	in.ID = id
	v, err := h.Store.UpdateEducation(r.Context(), currentUser(r).ID, in)
	h.writeItemResult(w, r, v, err)
}
func (h *handlers) deleteEducation(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	h.writeItemResult(w, r, nil, h.Store.DeleteEducation(r.Context(), currentUser(r).ID, id))
}

// experience
func (h *handlers) createExperience(w http.ResponseWriter, r *http.Request) {
	var in store.Experience
	if err := readJSON(w, r, &in); err != nil || strings.TrimSpace(in.Company) == "" {
		writeError(w, http.StatusBadRequest, "company is required")
		return
	}
	if !h.ensureResume(w, r) {
		return
	}
	v, err := h.Store.CreateExperience(r.Context(), currentUser(r).ID, in)
	h.writeItemResult(w, r, v, err)
}
func (h *handlers) updateExperience(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	var in store.Experience
	if err := readJSON(w, r, &in); !ok || err != nil || strings.TrimSpace(in.Company) == "" {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	in.ID = id
	v, err := h.Store.UpdateExperience(r.Context(), currentUser(r).ID, in)
	h.writeItemResult(w, r, v, err)
}
func (h *handlers) deleteExperience(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	h.writeItemResult(w, r, nil, h.Store.DeleteExperience(r.Context(), currentUser(r).ID, id))
}

// bullets: /api/me/resume/experience/{id}/bullets[/{bid}]
func (h *handlers) createBullet(w http.ResponseWriter, r *http.Request) {
	expID, ok := pathID(r, "id")
	var in store.Bullet
	if err := readJSON(w, r, &in); !ok || err != nil || strings.TrimSpace(in.Text) == "" {
		writeError(w, http.StatusBadRequest, "text is required")
		return
	}
	v, err := h.Store.CreateBullet(r.Context(), currentUser(r).ID, expID, in)
	h.writeItemResult(w, r, v, err)
}
func (h *handlers) updateBullet(w http.ResponseWriter, r *http.Request) {
	expID, ok1 := pathID(r, "id")
	bid, ok2 := pathID(r, "bid")
	var in store.Bullet
	if err := readJSON(w, r, &in); !ok1 || !ok2 || err != nil || strings.TrimSpace(in.Text) == "" {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	in.ID = bid
	v, err := h.Store.UpdateBullet(r.Context(), currentUser(r).ID, expID, in)
	h.writeItemResult(w, r, v, err)
}
func (h *handlers) deleteBullet(w http.ResponseWriter, r *http.Request) {
	expID, ok1 := pathID(r, "id")
	bid, ok2 := pathID(r, "bid")
	if !ok1 || !ok2 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	h.writeItemResult(w, r, nil, h.Store.DeleteBullet(r.Context(), currentUser(r).ID, expID, bid))
}
func (h *handlers) reorderBullets(w http.ResponseWriter, r *http.Request) {
	expID, ok := pathID(r, "id")
	var in struct {
		IDs []int64 `json:"ids"`
	}
	if err := readJSON(w, r, &in); !ok || err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	h.writeItemResult(w, r, nil, h.Store.ReorderBullets(r.Context(), currentUser(r).ID, expID, in.IDs))
}
