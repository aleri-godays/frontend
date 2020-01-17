package timetracking

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aleri-godays/frontend"
	"github.com/aleri-godays/frontend/internal/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	Add(ctx context.Context, entry *frontend.Entry, jwt string) (*frontend.Entry, error)
	Get(ctx context.Context, id int, jwt string) (*frontend.Entry, error)
	Update(ctx context.Context, entry *frontend.Entry, jwt string) error
	Delete(ctx context.Context, id int, jwt string) error
	GetByUser(ctx context.Context, user string, jwt string) ([]*frontend.Entry, error)
}

type client struct {
	conf       *config.Config
	httpClient *http.Client
}

type Entry struct {
	ID       int    `json:"id"`
	DateTS   int64  `json:"date_ts"`
	Project  int    `json:"project"`
	User     string `json:"user"`
	Comment  string `json:"comment"`
	Duration int64  `json:"duration"`
}

func NewClient(conf *config.Config) Client {
	c := &client{
		conf: conf,
		httpClient: &http.Client{
			Timeout: time.Second * 60,
		},
	}

	return c
}

func (c *client) Add(ctx context.Context, entry *frontend.Entry, jwt string) (*frontend.Entry, error) {
	span := newSpanFromContext(ctx, "Add")
	defer span.Finish()

	e := &Entry{
		ID:       entry.ID,
		DateTS:   entry.Date.Unix(),
		Project:  entry.ProjectID,
		User:     entry.User,
		Comment:  entry.Comment,
		Duration: entry.Duration,
	}

	req, err := c.buildHTTPRequest("/api/v1/timetracking", "POST", jwt, e, nil, span)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("request for entry failed (http code %d): %w", res.StatusCode, err)
	}

	var addedEntry Entry
	if err := json.NewDecoder(res.Body).Decode(&addedEntry); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	entry.ID = addedEntry.ID

	return entry, nil
}

func (c *client) Get(ctx context.Context, id int, jwt string) (*frontend.Entry, error) {
	span := newSpanFromContext(ctx, "Get")
	defer span.Finish()

	path := fmt.Sprintf("/api/v1/timetracking/%d", id)
	req, err := c.buildHTTPRequest(path, "GET", jwt, nil, nil, span)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request for entry id '%d' failed (http code %d): %w", id, res.StatusCode, err)
	}

	var entry Entry
	if err := json.NewDecoder(res.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	ee := &frontend.Entry{
		ID:        entry.ID,
		Date:      time.Unix(entry.DateTS, 0),
		ProjectID: entry.Project,
		User:      entry.User,
		Comment:   entry.Comment,
		Duration:  entry.Duration,
	}

	return ee, nil
}

func (c *client) Update(ctx context.Context, entry *frontend.Entry, jwt string) error {
	span := newSpanFromContext(ctx, "Update")
	defer span.Finish()

	e := &Entry{
		ID:       entry.ID,
		DateTS:   entry.Date.Unix(),
		Project:  entry.ProjectID,
		User:     entry.User,
		Comment:  entry.Comment,
		Duration: entry.Duration,
	}

	path := fmt.Sprintf("/api/v1/timetracking/%d", entry.ID)
	req, err := c.buildHTTPRequest(path, "PUT", jwt, e, nil, span)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("request for entry failed (http code %d): %w", res.StatusCode, err)
	}

	return nil
}

func (c *client) Delete(ctx context.Context, id int, jwt string) error {
	span := newSpanFromContext(ctx, "Delete")
	defer span.Finish()

	path := fmt.Sprintf("/api/v1/timetracking/%d", id)
	req, err := c.buildHTTPRequest(path, "DELETE", jwt, nil, nil, span)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return errors.New("entry does not exist")
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("request for entry failed (http code %d): %w", res.StatusCode, err)
	}

	return nil
}

func (c *client) GetByUser(ctx context.Context, user string, jwt string) ([]*frontend.Entry, error) {
	span := newSpanFromContext(ctx, "GetByUser")
	defer span.Finish()

	path := fmt.Sprintf("/api/v1/timetracking")
	req, err := c.buildHTTPRequest(path, "GET", jwt, nil, map[string]string{"user": user}, span)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request for entries failed (http code %d): %w", res.StatusCode, err)
	}

	var entries []*Entry
	if err := json.NewDecoder(res.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(entries) == 0 {
		return nil, nil
	}

	entryList := make([]*frontend.Entry, 0, len(entries))
	for _, entry := range entries {
		ee := &frontend.Entry{
			ID:        entry.ID,
			Date:      time.Unix(entry.DateTS, 0),
			ProjectID: entry.Project,
			User:      entry.User,
			Comment:   entry.Comment,
			Duration:  entry.Duration,
		}
		entryList = append(entryList, ee)
	}

	return entryList, nil
}

func (c *client) buildHTTPRequest(path, method, jwt string, data interface{}, params map[string]string, span opentracing.Span) (*http.Request, error) {
	var body io.Reader
	if data != nil {
		buf, err := json.Marshal(&data)
		if err != nil {
			return nil, fmt.Errorf("could not marshal to json: %w", err)
		}
		body = bytes.NewReader(buf)
	}

	uri := fmt.Sprintf("%s%s", c.conf.TimeTrackingURI, path)
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid url '%s': %w", uri, err)
	}
	for key, value := range params {
		u.Query().Set(key, value)
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, path)
	ext.HTTPMethod.Set(span, method)
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	return req, nil
}

func newSpanFromContext(ctx context.Context, opName string) opentracing.Span {
	span, _ := opentracing.StartSpanFromContext(ctx, "timetracking: "+opName)
	ext.Component.Set(span, "timetracking-client")
	return span
}
