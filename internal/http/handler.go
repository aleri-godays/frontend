package http

import (
	"fmt"
	"github.com/aleri-godays/frontend"
	"github.com/aleri-godays/frontend/internal/project"
	"github.com/aleri-godays/frontend/internal/timetracking"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type handler struct {
	entries  timetracking.Client
	projects project.Client
}

func (h *handler) Me(c echo.Context) error {
	login := c.Get("login").(string)
	token := c.Get("user_jwt").(string)
	return c.JSONPretty(http.StatusOK, map[string]string{
		"login": login,
		"jwt":   token,
	}, "  ")
}

type Entry struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	ProjectID int       `json:"project_id"`
	Comment   string    `json:"comment"`
	Duration  int64     `json:"duration"`
}

func (h *handler) AddEntry(c echo.Context) error {
	logger := c.Get("logger").(*log.Entry)
	jwt := c.Get("user_jwt").(string)

	entry, err := h.getEntryFromBody(c)
	if err != nil {
		return err
	}

	addedEntry, err := h.entries.Add(c.Request().Context(), entry, jwt)
	if err != nil {
		logger.WithFields(log.Fields{
			"user":       entry.User,
			"project_id": entry.ProjectID,
			"error":      err,
		}).Warn("could not add entry")
		return jsonError(c, "could not add entry", http.StatusInternalServerError)
	}

	return c.JSONPretty(http.StatusOK, map[string]int{"entry_id": addedEntry.ID}, " ")
}

func (h *handler) GetEntry(c echo.Context) error {
	jwt := c.Get("user_jwt").(string)

	id, err := h.getIDFromPath(c)
	if err != nil {
		return err
	}

	ee, err := h.entries.Get(c.Request().Context(), *id, jwt)
	if err != nil {
		return jsonError(c, "could not get entry", http.StatusInternalServerError)
	}

	if ee == nil {
		return c.NoContent(http.StatusNotFound)
	}

	entry := &Entry{
		ID:        ee.ID,
		Date:      ee.Date,
		ProjectID: ee.ProjectID,
		Comment:   ee.Comment,
		Duration:  ee.Duration,
	}

	return c.JSONPretty(http.StatusOK, entry, "  ")
}

func (h *handler) UpdateEntry(c echo.Context) error {
	jwt := c.Get("user_jwt").(string)

	entry, err := h.getEntryFromBody(c)
	if err != nil {
		return err
	}
	if err := h.entries.Update(c.Request().Context(), entry, jwt); err != nil {
		return jsonError(c, fmt.Sprintf("could not update entry %d", entry.ID), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *handler) DeleteEntry(c echo.Context) error {
	jwt := c.Get("user_jwt").(string)

	id, err := h.getIDFromPath(c)
	if err != nil {
		return err
	}
	if err := h.entries.Delete(c.Request().Context(), *id, jwt); err != nil {
		return jsonError(c, fmt.Sprintf("could not delete entry %d", *id), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *handler) GetEntries(c echo.Context) error {
	jwt := c.Get("user_jwt").(string)

	login := c.Get("login").(string)
	ees, err := h.entries.GetByUser(c.Request().Context(), login, jwt)
	if err != nil {
		return jsonError(c, fmt.Sprintf("could not get entriesfor user %s", login), http.StatusInternalServerError)
	}

	if len(ees) == 0 {
		return c.JSONPretty(http.StatusOK, []string{}, "  ")
	}

	entries := make([]*Entry, 0, len(ees))
	for _, ee := range ees {
		entry := &Entry{
			ID:        ee.ID,
			Date:      ee.Date,
			ProjectID: ee.ProjectID,
			Comment:   ee.Comment,
			Duration:  ee.Duration,
		}
		entries = append(entries, entry)
	}

	return c.JSONPretty(http.StatusOK, entries, "  ")
}

func (h *handler) AllProjects(c echo.Context) error {
	logger := c.Get("logger").(*log.Entry)
	jwt := c.Get("user_jwt").(string)

	projects, err := h.projects.GetAllProjects(c.Request().Context(), jwt)
	if err != nil {
		logger.WithFields(log.Fields{
			"error": err,
		}).Warn("could no request all projects")
		return jsonError(c, "failed ro request projects", http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, projects)
}

func (h *handler) GetProject(c echo.Context) error {
	logger := c.Get("logger").(*log.Entry)
	jwt := c.Get("user_jwt").(string)

	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return jsonError(c, "id must be an integer", http.StatusBadRequest)
	}

	prj, err := h.projects.GetProjectByID(c.Request().Context(), id, jwt)
	if err != nil {
		logger.WithFields(log.Fields{
			"error": err,
		}).Warn("could no request all projects")
		return jsonError(c, "failed ro request projects", http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, prj)
}

func (h *handler) getEntryFromBody(c echo.Context) (*frontend.Entry, error) {
	logger := c.Get("logger").(*log.Entry)

	var entry Entry
	if err := c.Bind(&entry); err != nil {
		logger.WithFields(log.Fields{
			"error": err,
		}).Warn("could parse request body")
		return nil, jsonError(c, "could parse request body", http.StatusBadRequest)
	}

	login := c.Get("login").(string)
	ee := &frontend.Entry{
		ID:        entry.ID,
		Date:      entry.Date,
		ProjectID: entry.ProjectID,
		User:      login,
		Comment:   entry.Comment,
		Duration:  entry.Duration,
	}

	return ee, nil
}

func (h *handler) getIDFromPath(c echo.Context) (*int, error) {
	entryID := c.Param("id")
	if entryID == "" {
		return nil, jsonError(c, "empty id", http.StatusBadRequest)
	}

	id, err := strconv.Atoi(entryID)
	if err != nil {
		return nil, jsonError(c, "id must be an integer", http.StatusBadRequest)
	}
	return &id, nil
}
