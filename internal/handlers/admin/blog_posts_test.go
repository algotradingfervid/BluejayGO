package admin

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func createTestCategory(t *testing.T, q *sqlc.Queries, name, slug, color string) sqlc.BlogCategory {
	t.Helper()
	cat, err := q.CreateBlogCategory(context.Background(), sqlc.CreateBlogCategoryParams{
		Name:     name,
		Slug:     slug,
		ColorHex: color,
	})
	if err != nil {
		t.Fatalf("failed to create test category: %v", err)
	}
	return cat
}

func createTestAuthor(t *testing.T, q *sqlc.Queries, name, slug string) sqlc.BlogAuthor {
	t.Helper()
	author, err := q.CreateBlogAuthor(context.Background(), sqlc.CreateBlogAuthorParams{
		Name:  name,
		Slug:  slug,
		Title: "Test Author",
	})
	if err != nil {
		t.Fatalf("failed to create test author: %v", err)
	}
	return author
}

func createTestPost(t *testing.T, q *sqlc.Queries, cat sqlc.BlogCategory, author sqlc.BlogAuthor, title, slug, status string, publishedAt sql.NullTime) sqlc.BlogPost {
	t.Helper()
	post, err := q.CreateBlogPost(context.Background(), sqlc.CreateBlogPostParams{
		Title:      title,
		Slug:       slug,
		Excerpt:    "Test excerpt for " + title,
		Body:       "<p>Test body for " + title + "</p>",
		CategoryID: cat.ID,
		AuthorID:   author.ID,
		Status:     status,
		PublishedAt: publishedAt,
	})
	if err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}
	return post
}

func TestListPublishedPosts_EmptyDatabase(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	posts, err := queries.ListPublishedPosts(context.Background(), sqlc.ListPublishedPostsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 0 {
		t.Errorf("expected 0 posts, got %d", len(posts))
	}
}

func TestListPublishedPosts_ReturnsOnlyPublished(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "News", "news", "#FF0000")
	author := createTestAuthor(t, queries, "John", "john")

	// Draft post
	createTestPost(t, queries, cat, author, "Draft Post", "draft-post", "draft", sql.NullTime{})
	// Published post
	createTestPost(t, queries, cat, author, "Published Post", "published-post", "published", sql.NullTime{Time: time.Now(), Valid: true})

	posts, err := queries.ListPublishedPosts(context.Background(), sqlc.ListPublishedPostsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("expected 1 published post, got %d", len(posts))
	}
}

func TestGetPublishedPostBySlug_NotFound(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	_, err := queries.GetPublishedPostBySlug(context.Background(), "non-existent")
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows, got %v", err)
	}
}

func TestGetPublishedPostBySlug_DraftNotReturned(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "News", "news", "#FF0000")
	author := createTestAuthor(t, queries, "John", "john")
	createTestPost(t, queries, cat, author, "Draft", "draft-slug", "draft", sql.NullTime{})

	_, err := queries.GetPublishedPostBySlug(context.Background(), "draft-slug")
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows for draft post, got %v", err)
	}
}

func TestListPublishedPostsByCategory(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	cat1 := createTestCategory(t, queries, "News", "news", "#FF0000")
	cat2 := createTestCategory(t, queries, "Updates", "updates", "#00FF00")
	author := createTestAuthor(t, queries, "John", "john")

	createTestPost(t, queries, cat1, author, "News Post", "news-post", "published", sql.NullTime{Time: time.Now(), Valid: true})
	createTestPost(t, queries, cat2, author, "Update Post", "update-post", "published", sql.NullTime{Time: time.Now(), Valid: true})

	posts, err := queries.ListPublishedPostsByCategory(context.Background(), sqlc.ListPublishedPostsByCategoryParams{
		Slug:   "news",
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("expected 1 post in news category, got %d", len(posts))
	}
}

func TestGetRelatedPosts_ExcludesCurrentPost(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "News", "news", "#FF0000")
	author := createTestAuthor(t, queries, "John", "john")

	post1 := createTestPost(t, queries, cat, author, "Post 1", "post-1", "published", sql.NullTime{Time: time.Now(), Valid: true})
	createTestPost(t, queries, cat, author, "Post 2", "post-2", "published", sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true})

	related, err := queries.GetRelatedPosts(context.Background(), sqlc.GetRelatedPostsParams{
		CategoryID: cat.ID,
		ID:         post1.ID,
		Limit:      3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(related) != 1 {
		t.Errorf("expected 1 related post, got %d", len(related))
	}
	if len(related) > 0 && related[0].Title == "Post 1" {
		t.Error("related posts should not include current post")
	}
}

func TestTagOperations(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "News", "news", "#FF0000")
	author := createTestAuthor(t, queries, "John", "john")
	post := createTestPost(t, queries, cat, author, "Tagged Post", "tagged-post", "published", sql.NullTime{Time: time.Now(), Valid: true})

	tag, err := queries.CreateBlogTag(context.Background(), sqlc.CreateBlogTagParams{
		Name: "Test Tag",
		Slug: "test-tag",
	})
	if err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	// Add tag
	err = queries.AddTagToPost(context.Background(), sqlc.AddTagToPostParams{
		BlogPostID: post.ID,
		BlogTagID:  tag.ID,
	})
	if err != nil {
		t.Fatalf("failed to add tag: %v", err)
	}

	tags, err := queries.GetPostTagsByPostID(context.Background(), post.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}

	// Clear tags
	err = queries.ClearPostTags(context.Background(), post.ID)
	if err != nil {
		t.Fatalf("failed to clear tags: %v", err)
	}

	tags, err = queries.GetPostTagsByPostID(context.Background(), post.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 tags after clear, got %d", len(tags))
	}
}

func TestCountPublishedPosts(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "News", "news", "#FF0000")
	author := createTestAuthor(t, queries, "John", "john")

	createTestPost(t, queries, cat, author, "Post 1", "post-1", "published", sql.NullTime{Time: time.Now(), Valid: true})
	createTestPost(t, queries, cat, author, "Post 2", "post-2", "published", sql.NullTime{Time: time.Now(), Valid: true})
	createTestPost(t, queries, cat, author, "Draft", "draft", "draft", sql.NullTime{})

	count, err := queries.CountPublishedPosts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}

func TestGetFeaturedPost_ReturnsMostRecent(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "News", "news", "#FF0000")
	author := createTestAuthor(t, queries, "John", "john")

	createTestPost(t, queries, cat, author, "Old Post", "old-post", "published", sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true})
	createTestPost(t, queries, cat, author, "New Post", "new-post", "published", sql.NullTime{Time: time.Now(), Valid: true})

	featured, err := queries.GetFeaturedPost(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if featured.Title != "New Post" {
		t.Errorf("expected featured post to be 'New Post', got '%s'", featured.Title)
	}
}
