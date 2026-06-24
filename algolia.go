// Package algolia is an Algolia driver for togo full-text search.
// Blank-import it and set SEARCH_DRIVER=algolia, ALGOLIA_APP_ID, ALGOLIA_API_KEY.
package algolia

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/togo-framework/search"
	"github.com/togo-framework/togo"
)

func init() {
	search.RegisterDriver("algolia", func(k *togo.Kernel) (search.Searcher, error) {
		app := os.Getenv("ALGOLIA_APP_ID")
		key := os.Getenv("ALGOLIA_API_KEY")
		if app == "" || key == "" {
			return nil, errors.New("search-algolia: ALGOLIA_APP_ID and ALGOLIA_API_KEY required")
		}
		return &searcher{app: app, key: key, client: &http.Client{Timeout: 15 * time.Second}}, nil
	})
}

type searcher struct {
	app, key string
	client   *http.Client
}

func (s *searcher) do(ctx context.Context, method, host, path string, body any) (*http.Response, error) {
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, "https://"+host+path, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Algolia-Application-Id", s.app)
	req.Header.Set("X-Algolia-API-Key", s.key)
	req.Header.Set("Content-Type", "application/json")
	return s.client.Do(req)
}

func (s *searcher) write() string { return s.app + ".algolia.net" }
func (s *searcher) read() string  { return s.app + "-dsn.algolia.net" }

func (s *searcher) Index(ctx context.Context, index, id string, doc map[string]any) error {
	d := map[string]any{"objectID": id}
	for k, v := range doc {
		d[k] = v
	}
	resp, err := s.do(ctx, http.MethodPut, s.write(), "/1/indexes/"+url.PathEscape(index)+"/"+url.PathEscape(id), d)
	if err != nil {
		return err
	}
	return drain(resp)
}

func (s *searcher) Search(ctx context.Context, index, query string, limit int) ([]search.Hit, error) {
	if limit <= 0 {
		limit = 20
	}
	resp, err := s.do(ctx, http.MethodPost, s.read(), "/1/indexes/"+url.PathEscape(index)+"/query", map[string]any{"query": query, "hitsPerPage": limit})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search-algolia: %s: %s", resp.Status, b)
	}
	var out struct {
		Hits []map[string]any `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	hits := make([]search.Hit, 0, len(out.Hits))
	for _, h := range out.Hits {
		id, _ := h["objectID"].(string)
		hits = append(hits, search.Hit{ID: id, Doc: h})
	}
	return hits, nil
}

func (s *searcher) Delete(ctx context.Context, index, id string) error {
	resp, err := s.do(ctx, http.MethodDelete, s.write(), "/1/indexes/"+url.PathEscape(index)+"/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	return drain(resp)
}

func drain(resp *http.Response) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("search-algolia: %s: %s", resp.Status, b)
	}
	io.Copy(io.Discard, resp.Body)
	return nil
}
