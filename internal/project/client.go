package project

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aleri-godays/frontend"
	"github.com/aleri-godays/frontend/internal/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io"
	"net/http"
	"time"
)

type Client interface {
	GetProjectByID(ctx context.Context, id int, jwt string) (*frontend.Project, error)
	GetAllProjects(ctx context.Context, jwt string) (*[]frontend.Project, error)
}

type client struct {
	conf       *config.Config
	httpClient *http.Client
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

func (c *client) GetProjectByID(ctx context.Context, id int, jwt string) (*frontend.Project, error) {
	span := newSpanFromContext(ctx, "GetProjectByID")
	defer span.Finish()

	url := fmt.Sprintf("/api/v1/project/%d", id)
	req, err := c.buildHTTPRequest(url, "GET", jwt, nil, span)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request for project id '%d' failed (http code %d): %w", id, res.StatusCode, err)
	}

	var project frontend.Project
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return &project, nil
}

func (c *client) GetAllProjects(ctx context.Context, jwt string) (*[]frontend.Project, error) {
	span := newSpanFromContext(ctx, "GetAllProjects")
	defer span.Finish()

	url := "/api/v1/project"
	req, err := c.buildHTTPRequest(url, "GET", jwt, nil, span)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request for all projects failed (http code %d): %w", res.StatusCode, err)
	}

	var projects []frontend.Project
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return &projects, nil
}

func (c *client) buildHTTPRequest(path, method, jwt string, body io.Reader, span opentracing.Span) (*http.Request, error) {
	uri := fmt.Sprintf("%s%s", c.conf.ProjectURI, path)
	req, err := http.NewRequest(method, uri, body)
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
	span, _ := opentracing.StartSpanFromContext(ctx, "project: "+opName)
	ext.Component.Set(span, "project-client")
	return span
}
