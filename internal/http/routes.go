package http

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//AddRoutes creates all routes
func (s *Server) AddRoutes(h *handler) {
	s.e.GET("/health", s.Health)

	api := s.e.Group("/api/v1")
	frontend := api.Group("/frontend")

	auth := newGithubAuth(s.conf)
	s.AddAuthRoutes(auth)

	frontend.Use(auth.AuthMiddleware())
	frontend.Use(AddLoginToContext())

	frontend.GET("/me", h.Me)

	frontend.POST("/entry", h.AddEntry)
	frontend.GET("/entry", h.GetEntries)
	frontend.GET("/entry/:id", h.GetEntry)
	frontend.PUT("/entry/:id", h.UpdateEntry)
	frontend.DELETE("/entry/:id", h.DeleteEntry)

	frontend.GET("/project", h.AllProjects)
	frontend.GET("/project/:id", h.GetProject)

	s.e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}

//AddRoutes creates all routes
func (s *Server) AddAuthRoutes(auth *githubAuth) {
	//auth related routes
	s.e.POST("/login", auth.Login)
	s.e.GET("/login", auth.Login)
	s.e.POST("/logout", auth.Logout)
	s.e.GET("/logout", auth.Logout)
	s.e.GET("/authcallback", auth.Callback)
	s.e.POST("/authcallback", auth.Callback)
}
