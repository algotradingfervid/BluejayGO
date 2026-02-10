package public

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

type SitemapHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
	baseURL string
}

func NewSitemapHandler(queries *sqlc.Queries, logger *slog.Logger, baseURL string) *SitemapHandler {
	return &SitemapHandler{queries: queries, logger: logger, baseURL: baseURL}
}

func (h *SitemapHandler) Sitemap(c echo.Context) error {
	now := time.Now().Format("2006-01-02")

	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	// Static pages
	staticPages := []struct {
		path       string
		changeFreq string
		priority   string
	}{
		{"/", "weekly", "1.0"},
		{"/products", "weekly", "0.9"},
		{"/solutions", "weekly", "0.9"},
		{"/blog", "daily", "0.8"},
		{"/case-studies", "weekly", "0.8"},
		{"/whitepapers", "weekly", "0.8"},
		{"/about", "monthly", "0.7"},
		{"/contact", "monthly", "0.6"},
		{"/partners", "monthly", "0.7"},
	}

	for _, page := range staticPages {
		urlset.URLs = append(urlset.URLs, URL{
			Loc:        h.baseURL + page.path,
			LastMod:    now,
			ChangeFreq: page.changeFreq,
			Priority:   page.priority,
		})
	}

	// Solutions
	solutions, err := h.queries.ListPublishedSolutions(c.Request().Context())
	if err != nil {
		h.logger.Error("sitemap: failed to list solutions", "error", err)
	} else {
		for _, s := range solutions {
			u := URL{
				Loc:        fmt.Sprintf("%s/solutions/%s", h.baseURL, s.Slug),
				ChangeFreq: "weekly",
				Priority:   "0.8",
			}
			if s.UpdatedAt.Valid {
				u.LastMod = s.UpdatedAt.Time.Format("2006-01-02")
			}
			urlset.URLs = append(urlset.URLs, u)
		}
	}

	// Blog posts
	posts, err := h.queries.ListPublishedPosts(c.Request().Context(), sqlc.ListPublishedPostsParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		h.logger.Error("sitemap: failed to list blog posts", "error", err)
	} else {
		for _, p := range posts {
			urlset.URLs = append(urlset.URLs, URL{
				Loc:        fmt.Sprintf("%s/blog/%s", h.baseURL, p.Slug),
				ChangeFreq: "monthly",
				Priority:   "0.7",
			})
		}
	}

	// Case studies
	caseStudies, err := h.queries.ListCaseStudies(c.Request().Context())
	if err != nil {
		h.logger.Error("sitemap: failed to list case studies", "error", err)
	} else {
		for _, cs := range caseStudies {
			urlset.URLs = append(urlset.URLs, URL{
				Loc:        fmt.Sprintf("%s/case-studies/%s", h.baseURL, cs.Slug),
				ChangeFreq: "monthly",
				Priority:   "0.7",
			})
		}
	}

	// Whitepapers
	whitepapers, err := h.queries.ListPublishedWhitepapers(c.Request().Context())
	if err != nil {
		h.logger.Error("sitemap: failed to list whitepapers", "error", err)
	} else {
		for _, w := range whitepapers {
			urlset.URLs = append(urlset.URLs, URL{
				Loc:        fmt.Sprintf("%s/whitepapers/%s", h.baseURL, w.Slug),
				ChangeFreq: "monthly",
				Priority:   "0.7",
			})
		}
	}

	xmlData, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		h.logger.Error("sitemap: failed to marshal XML", "error", err)
		return c.String(http.StatusInternalServerError, "failed to generate sitemap")
	}

	xmlData = append([]byte(xml.Header), xmlData...)
	return c.Blob(http.StatusOK, "application/xml", xmlData)
}

func (h *SitemapHandler) RobotsTxt(c echo.Context) error {
	robots := `User-agent: *
Allow: /
Disallow: /admin/

Sitemap: ` + h.baseURL + `/sitemap.xml`
	return c.String(http.StatusOK, robots)
}
