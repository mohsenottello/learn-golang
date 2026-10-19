package httpapi

// Public, unauthenticated read of a published resume, localized:
//   GET /api/resumes/{slug}?lang=de|en   (default: the resume's default_lang)

import (
	"errors"
	"net/http"

	"github.com/mohsenottello/golang/internal/store"
)

// pick returns the English text when asked for and available, else the default.
func pick(lang, def, en string) string {
	if lang == "en" && en != "" {
		return en
	}
	return def
}

// localize collapses each translatable pair into the base field and blanks the *_en twins,
// so the public page only ever sees one language.
func localize(f *FullResume, lang string) {
	f.Lang = lang
	r := f.Resume
	r.Title, r.TitleEn = pick(lang, r.Title, r.TitleEn), ""
	r.Subtitle, r.SubtitleEn = pick(lang, r.Subtitle, r.SubtitleEn), ""
	r.About, r.AboutEn = pick(lang, r.About, r.AboutEn), ""
	r.UserID = 0

	for i := range f.Languages {
		l := &f.Languages[i]
		l.Name, l.NameEn = pick(lang, l.Name, l.NameEn), ""
		l.Level, l.LevelEn = pick(lang, l.Level, l.LevelEn), ""
	}
	for i := range f.Education {
		e := &f.Education[i]
		e.Degree, e.DegreeEn = pick(lang, e.Degree, e.DegreeEn), ""
		e.School, e.SchoolEn = pick(lang, e.School, e.SchoolEn), ""
	}
	for i := range f.Experience {
		x := &f.Experience[i]
		x.Role, x.RoleEn = pick(lang, x.Role, x.RoleEn), ""
		x.Subtitle, x.SubtitleEn = pick(lang, x.Subtitle, x.SubtitleEn), ""
		x.Location, x.LocationEn = pick(lang, x.Location, x.LocationEn), ""
		for j := range x.Bullets {
			b := &x.Bullets[j]
			b.Text, b.TextEn = pick(lang, b.Text, b.TextEn), ""
		}
	}
}

func (h *handlers) publicResume(w http.ResponseWriter, r *http.Request) {
	res, err := h.Store.ResumeBySlug(r.Context(), r.PathValue("slug"))
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "resume not found")
		return
	case err != nil:
		serverError(w, r, err)
		return
	}
	owner, err := h.Store.GetUser(r.Context(), res.UserID)
	if err != nil {
		serverError(w, r, err)
		return
	}
	full, err := h.loadFull(r, res, owner.Name)
	if err != nil {
		serverError(w, r, err)
		return
	}
	lang := r.URL.Query().Get("lang")
	if !validLangs[lang] {
		lang = res.DefaultLang
	}
	localize(full, lang)
	writeJSON(w, http.StatusOK, full)
}

func (h *handlers) publicPhoto(w http.ResponseWriter, r *http.Request) {
	data, ct, err := h.Store.ResumePhotoBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Write(data)
}
