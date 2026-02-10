# Test Plan: Blog Post Tag Assignment

## Summary
Testing the HTMX-driven tag search, selection, quick-create, and removal flow within blog post create/edit forms.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with 13 tags and blog posts
- On blog post create/edit form at /admin/blog/posts/new or /admin/blog/posts/:id/edit

## User Journey Steps
1. Navigate to http://localhost:28090/admin/blog/posts/new
2. Locate tag-search input (id="tag-search")
3. Type search query to trigger hx-get="/admin/blog/tags/search"
4. Verify suggestions appear in #tag-suggestions after 200ms delay
5. Click a suggestion to call addTag(id, name) JavaScript function
6. Verify tag chip with hidden tag_ids[] input added to #selected-tags
7. Type new tag name in search
8. Click "Quick Create" button in suggestions dropdown
9. Verify hx-post="/admin/blog/tags/quick-create" creates tag and appends chip
10. Click X on tag chip to remove tag
11. Verify hidden tag_ids[] input removed from #selected-tags

## Test Cases

### Happy Path
- **Search existing tags**: Type "java", see matching tags in dropdown within 200ms
- **Select tag from suggestions**: Click "JavaScript" suggestion, chip added to #selected-tags with hidden input tag_ids[]=:id
- **Multiple tag selection**: Add 3 different tags, verify 3 chips and 3 hidden inputs exist
- **Quick-create new tag**: Type "DevOps", click "Quick Create", new tag created and chip appended
- **Remove tag**: Click X on "JavaScript" chip, chip and hidden input removed from DOM
- **Focus trigger**: Focusing tag-search input shows recent/all suggestions
- **Form submission**: Submit post form with 3 tags, verify tag_ids[] array sent to server

### Edge Cases / Error States
- **Input delay**: Typing rapidly waits 200ms after last keystroke before HTMX request
- **Empty search**: Focusing empty search input shows all tags or placeholder
- **No results**: Search for "xyz123nonexistent" shows "No tags found" or quick-create option
- **Duplicate tag selection**: Selecting same tag twice prevented by JS or shows warning
- **Quick-create duplicate name**: Creating tag with existing name handled gracefully
- **Remove last tag**: Removing all tags leaves #selected-tags empty
- **JavaScript disabled**: Fallback behavior if addTag() function unavailable

## Selectors & Elements
- Tag search input: id="tag-search", hx-get="/admin/blog/tags/search", hx-trigger="input changed delay:200ms, focus", hx-target="#tag-suggestions"
- Suggestions container: id="tag-suggestions" (receives tag_suggestions.html partial)
- Selected tags container: id="selected-tags" (contains tag chips)
- Tag chip structure: <div class="tag-chip">TagName <button onclick="removeTag(id)">X</button><input type="hidden" name="tag_ids[]" value=":id"></div>
- Quick-create button: hx-post="/admin/blog/tags/quick-create" hx-vals='{"name":"..."}' hx-target="#selected-tags" hx-swap="beforeend"
- JavaScript functions: addTag(id, name), removeTag(id)

## HTMX Interactions
- **Tag search**: hx-get="/admin/blog/tags/search?q=keyword" hx-trigger="input changed delay:200ms, focus" hx-target="#tag-suggestions" (returns tag_suggestions.html partial with clickable suggestions)
- **Quick-create**: hx-post="/admin/blog/tags/quick-create" hx-vals='{"name":"NewTag"}' hx-target="#selected-tags" hx-swap="beforeend" (returns tag_chip.html partial appended to container)
- **Chip addition**: Clicking suggestion triggers addTag() JS which manually creates chip HTML and appends to #selected-tags
- **Chip removal**: Clicking X button calls removeTag() JS which removes chip element from DOM

## Dependencies
- 16-blog-tags.md (tag search endpoint and quick-create endpoint)
- 15-blog-posts-crud.md (parent blog post form)
