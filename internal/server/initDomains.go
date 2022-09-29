package server

import (
	"github.com/gmhafiz/go8/internal/domain/report"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gmhafiz/go8/internal/middleware"
	"github.com/gmhafiz/go8/internal/utility/respond"
)

type Domain struct {
}

func (s *Server) InitDomains() {
	s.rootRouter.Route("/api/v0", func(apiRouter chi.Router) {
		apiRouter.Use(middleware.Json)
		apiRouter.Route("/version", s.versionRouter)
		apiRouter.Route("/report", report.ReportRouter)
	})
	s.initSwagger()
}

func (s *Server) versionRouter(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		respond.Json(w, http.StatusOK, map[string]string{"version": s.Version})
	})

}

func (s *Server) initSwagger() {
	if s.Config().Api.RunSwagger {
		fileServer := http.FileServer(http.Dir(swaggerDocsAssetPath))

		s.rootRouter.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
		})
		s.rootRouter.Handle("/swagger/", http.StripPrefix("/swagger", middleware.ContentType(fileServer)))
		s.rootRouter.Handle("/swagger/*", http.StripPrefix("/swagger", middleware.ContentType(fileServer)))
	}
}
