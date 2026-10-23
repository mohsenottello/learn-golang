// Package httpapi is everything that knows about HTTP: routing, middleware,
// request decoding, response encoding. It calls into store and mail; it never
// touches SQL directly.
package httpapi

import (
	"io/fs"
	"net/http"

	"github.com/mohsenottello/golang/internal/mail"
	"github.com/mohsenottello/golang/internal/store"
)

type Deps struct {
	Store   *store.Store
	Mailer  *mail.Mailer
	BaseURL string
	Static  fs.FS
}

type handlers struct{ Deps }

// NewRouter is the single map of the API.
func NewRouter(d Deps) http.Handler {
	h := &handlers{d}
	mux := http.NewServeMux()
	admin := h.requireAdmin
	user := h.requireUser

	// Users: self-service auth
	mux.HandleFunc("POST /api/auth/register", h.userRegister)
	mux.HandleFunc("POST /api/auth/verify", h.userVerify)
	mux.HandleFunc("POST /api/auth/login", h.userLogin)
	mux.HandleFunc("POST /api/auth/logout", h.userLogout)
	mux.HandleFunc("POST /api/auth/forgot-password", h.forgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", h.resetPassword)
	mux.Handle("GET /api/auth/me", user(h.userMe))
	mux.Handle("PUT /api/auth/me", user(h.updateName))
	mux.Handle("PUT /api/auth/password", user(h.changePassword))

	// Users: own resume
	mux.Handle("GET /api/me/resume", user(h.getMyResume))
	mux.Handle("PUT /api/me/resume", user(h.updateMyResume))
	mux.Handle("PUT /api/me/resume/photo", user(h.uploadMyPhoto))
	mux.Handle("GET /api/me/resume/photo", user(h.myPhoto))

	mux.Handle("PUT /api/me/resume/order/{section}", user(h.reorderSection))

	mux.Handle("POST /api/me/resume/software", user(h.createSoftware))
	mux.Handle("PUT /api/me/resume/software/{id}", user(h.updateSoftware))
	mux.Handle("DELETE /api/me/resume/software/{id}", user(h.deleteSoftware))

	mux.Handle("POST /api/me/resume/languages", user(h.createLanguage))
	mux.Handle("PUT /api/me/resume/languages/{id}", user(h.updateLanguage))
	mux.Handle("DELETE /api/me/resume/languages/{id}", user(h.deleteLanguage))

	mux.Handle("POST /api/me/resume/education", user(h.createEducation))
	mux.Handle("PUT /api/me/resume/education/{id}", user(h.updateEducation))
	mux.Handle("DELETE /api/me/resume/education/{id}", user(h.deleteEducation))

	mux.Handle("POST /api/me/resume/experience", user(h.createExperience))
	mux.Handle("PUT /api/me/resume/experience/{id}", user(h.updateExperience))
	mux.Handle("DELETE /api/me/resume/experience/{id}", user(h.deleteExperience))
	mux.Handle("POST /api/me/resume/experience/{id}/bullets", user(h.createBullet))
	mux.Handle("PUT /api/me/resume/experience/{id}/bullets/{bid}", user(h.updateBullet))
	mux.Handle("DELETE /api/me/resume/experience/{id}/bullets/{bid}", user(h.deleteBullet))
	mux.Handle("PUT /api/me/resume/experience/{id}/bullets-order", user(h.reorderBullets))

	// Public resumes
	mux.HandleFunc("GET /api/resumes/{slug}", h.publicResume)
	mux.HandleFunc("GET /api/resumes/{slug}/photo", h.publicPhoto)

	// Admins
	mux.HandleFunc("POST /api/admin/login", h.adminLogin)
	mux.HandleFunc("POST /api/admin/logout", h.adminLogout)
	mux.Handle("GET /api/admin/me", admin(h.adminMe))
	mux.Handle("GET /api/admin/users", admin(h.adminListUsers)) // replaces listUsers
	mux.Handle("PUT /api/admin/users/{id}", admin(h.adminUpdateUser))
	mux.Handle("DELETE /api/admin/users/{id}", admin(h.deleteUser))
	mux.Handle("GET /api/admin/resumes", admin(h.adminListResumes))
	mux.Handle("PUT /api/admin/resumes/{id}/published", admin(h.adminSetPublished))
	mux.Handle("PUT /api/admin/password", admin(h.adminChangePassword))

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.Handle("/", spaHandler(d.Static))

	return recoverPanic(logRequests(mux))
}
