# Test Plan: Footer Configuration Management

## Summary
Tests the complex footer management interface with dynamic column configuration, link management, social settings, and legal links using JavaScript for dynamic form manipulation.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Footer settings initialized with default values
- JavaScript functions available: updateColumnCount(), updateColumnType(), addLink(), addLegalLink()
- Delete-and-recreate strategy for columns and links

## User Journey Steps
1. Navigate to /admin/footer
2. View layout section with column count radio (2/3/4) and background style radio (dark/light/primary)
3. Select footer_columns (e.g., 3 columns)
4. Verify updateColumnCount() JavaScript shows appropriate column config sections
5. For each column (col_0, col_1, col_2):
   - Set col_N_heading
   - Select col_N_type radio (links/text/contact)
   - If type=text: fill col_N_content textarea
   - If type=links: add multiple links using addLink() (col_N_link_label[], col_N_link_url[])
   - If type=contact: appropriate contact info display
6. Configure social section: check footer_show_social, select footer_social_style radio
7. Fill footer_copyright text
8. Add legal links using addLegalLink() (legal_link_label[], legal_link_url[])
9. Submit form
10. Verify redirect to /admin/footer?saved=1
11. Confirm success message and verify footer renders correctly
12. Test removing links and columns

## Test Cases

### Happy Path
- **Load footer settings**: Verifies GET /admin/footer loads with current configuration
- **Select 2 columns**: Chooses footer_columns=2, verifies 2 column sections appear
- **Select 3 columns**: Chooses footer_columns=3, verifies 3 column sections appear
- **Select 4 columns**: Chooses footer_columns=4, verifies 4 column sections appear
- **Change background dark**: Selects footer_bg_style=dark, verifies save
- **Change background light**: Selects footer_bg_style=light, verifies save
- **Column type links**: Sets col_0_type=links, adds 3 links, verifies save
- **Column type text**: Sets col_1_type=text, fills col_1_content textarea, verifies save
- **Column type contact**: Sets col_2_type=contact, verifies contact info display
- **Add link to column**: Clicks addLink() button, fills label and URL, verifies added
- **Remove link from column**: Removes link using delete/remove button, verifies removal
- **Add legal link**: Clicks addLegalLink() button, fills label and URL, verifies added
- **Remove legal link**: Removes legal link, verifies removal
- **Enable footer social**: Checks footer_show_social, selects style, verifies save
- **Disable footer social**: Unchecks footer_show_social, verifies social hidden
- **Update copyright**: Changes footer_copyright text, verifies save
- **Column heading**: Updates col_0_heading, verifies display

### Edge Cases / Error States
- **No columns selected**: Tests if footer_columns must be selected
- **Column count change**: Changes from 4 columns to 2, verifies col_2 and col_3 removed
- **Empty column heading**: Leaves col_N_heading empty, checks validation
- **Column type not selected**: Doesn't select col_N_type radio, checks default or validation
- **Links type with no links**: Selects col_N_type=links but adds no links, checks validation/empty
- **Text type with empty content**: Selects col_N_type=text but leaves col_N_content empty, checks validation
- **Link with empty label**: Adds link with URL but no label, checks validation
- **Link with empty URL**: Adds link with label but no URL, checks validation
- **Invalid link URL**: Enters malformed URL in col_N_link_url[], checks validation
- **Duplicate link labels**: Adds multiple links with same label, verifies handling
- **Legal link validation**: Tests empty label/URL in legal links
- **Invalid legal URL**: Enters malformed URL in legal_link_url[], checks validation
- **Empty copyright**: Leaves footer_copyright empty, checks if required
- **Very long copyright**: Enters 500+ chars in footer_copyright, checks limit
- **Social style without enable**: Selects footer_social_style without checking footer_show_social, checks handling
- **Delete-and-recreate strategy**: Verifies backend deletes old columns/links and recreates from form data
- **JavaScript errors**: Tests updateColumnCount() with invalid input, updateColumnType() with invalid column

## Selectors & Elements
- Form: `form[action="/admin/footer"][method="POST"]`
- Column count radios: `input[name="footer_columns"][value="2"]`, `[value="3"]`, `[value="4"]`
- Background style radios: `input[name="footer_bg_style"][value="dark"]`, `[value="light"]`, `[value="primary"]`
- Column sections: `#column-0`, `#column-1`, `#column-2`, `#column-3`
- Column heading: `input[name="col_0_heading"]`, `input[name="col_1_heading"]`, etc.
- Column type radios: `input[name="col_0_type"][value="links"]`, `[value="text"]`, `[value="contact"]`
- Column content textarea: `textarea[name="col_0_content"]`, `textarea[name="col_1_content"]`, etc.
- Link label inputs: `input[name="col_0_link_label[]"]`, `input[name="col_1_link_label[]"]`, etc.
- Link URL inputs: `input[name="col_0_link_url[]"]`, `input[name="col_1_link_url[]"]`, etc.
- Add link button: `button.add-link[data-column="0"]`, `[data-column="1"]`, etc.
- Remove link button: `button.remove-link` or `.delete-link`
- Footer show social checkbox: `input[name="footer_show_social"][type="checkbox"]`
- Footer social style radios: `input[name="footer_social_style"][value="icons"]`, `[value="icons_labels"]`
- Copyright input: `input[name="footer_copyright"]` or `textarea[name="footer_copyright"]`
- Legal link labels: `input[name="legal_link_label[]"]`
- Legal link URLs: `input[name="legal_link_url[]"]`
- Add legal link button: `button#add-legal-link` or `.add-legal-link`
- Remove legal link button: `button.remove-legal-link`
- Submit button: `button[type="submit"]`
- Success banner: `.alert-success` (when ?saved=1)

## HTMX Interactions
- None - standard form POST with full page redirect
- JavaScript handles dynamic column/link addition and removal client-side

## Dependencies
- Database footer settings table with complex structure
- JavaScript file with updateColumnCount(), updateColumnType(), addLink(), addLegalLink() functions
- Delete-and-recreate logic in backend handler
- Template: templates/admin/pages/footer-management.html
- Handler: internal/handlers/footer.go (GetFooterSettings, PostFooterSettings)
- Column types: links, text, contact
- Background styles: dark, light, primary
