package admin_test

import (
	"context"
	"database/sql"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	admin "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func blogPublicationFixture(t *testing.T) (*sqlc.Queries, *echo.Echo, *admin.BlogPostsHandler, url.Values) {
	t.Helper()
	_, q, cleanup := testutil.SetupTestDB(t)
	t.Cleanup(cleanup)
	ctx := context.Background()
	cat, err := q.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{Name: "News", Slug: "publication-news", ColorHex: "#000000"})
	if err != nil {
		t.Fatal(err)
	}
	author, err := q.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{Name: "Editor", Slug: "publication-editor", Title: "Editor"})
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	f := url.Values{"title": {"Publication test"}, "slug": {"publication-test"}, "category_id": {strconv.FormatInt(cat.ID, 10)}, "author_id": {strconv.FormatInt(author.ID, 10)}, "submit_action": {"publish"}}
	return q, e, admin.NewBlogPostsHandler(q, logger, services.NewCache()), f
}

func TestBlogPublicationRejectsEmptyBodyAndAllowsIncompleteDraft(t *testing.T) {
	q, e, h, f := blogPublicationFixture(t)
	for i, tc := range []struct {
		name, body, action string
		allowed            bool
	}{
		{"empty", "", "publish", false},
		{"whitespace", " \n\t", "publish", false},
		{"editor empty markup", "<div><br></div>", "publish", false},
		{"entities and format characters", "<p>&nbsp;&#160;&#x200b;&#xfeff;</p>", "publish", false},
		{"comment", "<!-- article not written -->", "publish", false},
		{"script and style", "<script>text</script><style>p {color:red}</style>", "publish", false},
		{"template", "<template><p>hidden story</p></template>", "publish", false},
		{"hidden", "<div hidden><p>hidden story</p></div>", "publish", false},
		{"empty image", "<img alt='description only'>", "publish", false},
		{"orphan source", "<source src='/uploads/article.mp4'>", "publish", false},
		{"plain text", "An article", "publish", true},
		{"formatted text", "<p>Saved <strong>story</strong></p>", "publish", true},
		{"image", "<img src='/uploads/article.png' alt='Chart'>", "publish", true},
		{"video", "<video controls><source src='/uploads/article.mp4'></video>", "publish", true},
		{"unfinished draft", "", "draft", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before, err := q.CountBlogPostsAdminFiltered(context.Background(), sqlc.CountBlogPostsAdminFilteredParams{FilterStatus: "", FilterCategory: 0, FilterAuthor: 0, FilterSearch: ""})
			if err != nil {
				t.Fatal(err)
			}
			f.Set("slug", "publication-test-"+strconv.Itoa(i))
			f.Set("body", tc.body)
			f.Set("submit_action", tc.action)
			rec, c := postForm(e, "/admin/blog/posts", f)
			if err := h.Create(c); err != nil {
				t.Fatal(err)
			}
			after, err := q.CountBlogPostsAdminFiltered(context.Background(), sqlc.CountBlogPostsAdminFilteredParams{FilterStatus: "", FilterCategory: 0, FilterAuthor: 0, FilterSearch: ""})
			if err != nil {
				t.Fatal(err)
			}
			if tc.allowed {
				if rec.Code != http.StatusSeeOther || after != before+1 {
					t.Fatalf("allowed submission: HTTP %d, count %d→%d", rec.Code, before, after)
				}
			} else {
				if rec.Code != http.StatusUnprocessableEntity || after != before {
					t.Fatalf("invalid publication wrote a record: HTTP %d, count %d→%d", rec.Code, before, after)
				}
				if !strings.Contains(rec.Body.String(), "Add article text or media before publishing") || !strings.Contains(rec.Body.String(), "data-dirty-on-load") {
					t.Fatal("missing actionable error or unsaved state")
				}
				if strings.Contains(rec.Body.String(), "Preview saved version") {
					t.Fatal("unsaved new post must not offer a saved preview")
				}
			}
		})
	}
}

func TestBlogPublicationFailureRetainsInputsAndExistingPublishedArticle(t *testing.T) {
	q, e, h, f := blogPublicationFixture(t)
	ctx := context.Background()
	f.Set("body", "<p>Original published article</p>")
	rec, c := postForm(e, "/admin/blog/posts", f)
	if err := h.Create(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Fatal(rec.Code)
	}
	post, err := q.GetPublishedPostBySlug(ctx, f.Get("slug"))
	if err != nil {
		t.Fatal(err)
	}
	id := strconv.FormatInt(post.ID, 10)
	tag, err := q.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "Retained tag", Slug: "retained-tag"})
	if err != nil {
		t.Fatal(err)
	}
	cat, err := q.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{Name: "Devices", Slug: "devices", Description: "Devices", Icon: "devices"})
	if err != nil {
		t.Fatal(err)
	}
	product, err := q.CreateProduct(ctx, sqlc.CreateProductParams{Name: "Retained product", Slug: "retained-product", Sku: "RETAINED", Description: "Test product", CategoryID: cat.ID, Status: "published"})
	if err != nil {
		t.Fatal(err)
	}
	f.Set("title", "Revised title & details")
	f.Set("slug", "revised-title")
	f.Set("body", "<div><br></div>")
	f.Set("excerpt", "Retained excerpt")
	f.Set("featured_image_url", "https://example.com/image.png")
	f.Set("featured_image_alt", "Retained image description")
	f.Set("meta_title", "Retained SEO title")
	f.Set("meta_description", "Retained SEO description")
	f.Set("published_at", "2030-01-02T12:30")
	f.Set("reading_time_minutes", "7")
	f.Set("tag_ids", strconv.FormatInt(tag.ID, 10))
	f.Set("product_ids", strconv.FormatInt(product.ID, 10))
	rec, c = postForm(e, "/admin/blog/posts/"+id, f)
	c.SetParamNames("id")
	c.SetParamValues(id)
	if err := h.Update(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status %d", rec.Code)
	}
	for _, key := range []string{"title", "slug", "body", "excerpt", "featured_image_url", "featured_image_alt", "meta_title", "meta_description", "published_at"} {
		if !strings.Contains(rec.Body.String(), html.EscapeString(f.Get(key))) {
			t.Errorf("lost submitted %s", key)
		}
	}
	for _, text := range []string{"Retained tag", "Retained product", `value="7"`, `aria-invalid="true"`, `data-dirty-on-load`, `name="category_id"`, `name="author_id"`} {
		if !strings.Contains(rec.Body.String(), text) {
			t.Errorf("missing retained field/state %q", text)
		}
	}
	if !strings.Contains(rec.Body.String(), `href="/blog/publication-test?preview=true"`) {
		t.Fatal("failed slug edit changed the saved-version preview link")
	}
	saved, err := q.GetBlogPost(ctx, post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if saved.Title != "Publication test" || saved.Body != "<p>Original published article</p>" || saved.Status != "published" {
		t.Fatal("invalid update changed published content")
	}
	tags, _ := q.GetPostTagsByPostID(ctx, post.ID)
	products, _ := q.GetPostProductsByPostID(ctx, post.ID)
	if len(tags) != 0 || len(products) != 0 {
		t.Fatal("invalid update mutated associations")
	}
	// Save the same unfinished work explicitly as a draft, then correct and publish it.
	for _, action := range []string{"draft", "publish"} {
		f.Set("submit_action", action)
		if action == "publish" {
			f.Set("body", "<p>Complete revised article</p>")
		}
		rec, c = postForm(e, "/admin/blog/posts/"+id, f)
		c.SetParamNames("id")
		c.SetParamValues(id)
		if err := h.Update(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusSeeOther {
			t.Fatal(rec.Code)
		}
		_, err = q.GetPublishedPostBySlug(ctx, "revised-title")
		if action == "draft" && err != sql.ErrNoRows {
			t.Fatal("unfinished draft remained public")
		}
		if action == "publish" && err != nil {
			t.Fatal("corrected article did not publish", err)
		}
	}
}
