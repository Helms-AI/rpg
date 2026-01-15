# bookmark-manager

A personal bookmark manager with tagging, full-text search, and automatic metadata extraction.

## Overview

This is a self-hosted web application for saving and organizing bookmarks. When a user saves a URL, the system automatically fetches the page title, description, and favicon. Users can organize bookmarks with tags and find them later through search.

The application should feel fast and responsive. Use optimistic UI updates where appropriate.

## Target Languages

- typescript
- go
- python

## Architecture

This is a traditional fullstack web application with a REST API backend and a single-page frontend.

### Backend

The backend exposes a REST API for managing bookmarks. It handles:
- CRUD operations for bookmarks
- Background fetching of URL metadata
- Full-text search across titles, descriptions, and tags
- Tag management and suggestions

Use SQLite for storage - it's simple to deploy and sufficient for personal use. The full-text search should use SQLite's FTS5 extension.

### Frontend

A simple, fast SPA built with vanilla JavaScript or a lightweight framework. No build step required for the MVP - just serve static files.

The UI should have:
- A prominent input for adding new URLs
- A searchable list of bookmarks
- Tag filtering sidebar
- Quick actions (edit, delete, copy URL)

## Data Model

### Bookmark

The core entity. Each bookmark represents a saved URL.

- `id` - unique identifier
- `url` - the saved URL (required, must be valid URL)
- `title` - page title (auto-fetched, editable)
- `description` - page description or user notes
- `favicon` - URL to the site's favicon
- `tags` - list of associated tags
- `created_at` - when the bookmark was saved
- `updated_at` - last modification time

### Tag

Tags are created implicitly when added to bookmarks. Track usage count for suggestions.

- `name` - the tag text (lowercase, no spaces)
- `count` - number of bookmarks using this tag

## API Design

All endpoints return JSON. Use standard HTTP status codes.

### Bookmarks

```
GET    /api/bookmarks          - list all bookmarks (supports ?q=search&tag=filter)
POST   /api/bookmarks          - create bookmark (body: {url, tags?})
GET    /api/bookmarks/:id      - get single bookmark
PUT    /api/bookmarks/:id      - update bookmark
DELETE /api/bookmarks/:id      - delete bookmark
```

### Tags

```
GET    /api/tags               - list all tags with counts
GET    /api/tags/:name         - get bookmarks by tag
```

### Metadata

```
POST   /api/fetch-metadata     - fetch title/description/favicon for a URL
```

## Key Behaviors

### Adding a Bookmark

When a user submits a URL:
1. Immediately create the bookmark with just the URL
2. Return success to the user (optimistic)
3. In the background, fetch the page metadata
4. Update the bookmark with title, description, favicon
5. If the frontend is connected, push the update

### Search

Search should query across:
- Bookmark titles
- Bookmark descriptions
- Bookmark URLs
- Tag names

Results should be ranked by relevance. Consider boosting exact matches and recent bookmarks.

### Tag Suggestions

When editing tags, suggest:
- Existing tags that match the typed prefix
- Most frequently used tags
- Tags commonly used together with already-selected tags

## Configuration

The app should be configurable via environment variables:

- `PORT` - HTTP port (default: 3000)
- `DATABASE_PATH` - SQLite file location (default: ./data/bookmarks.db)
- `FETCH_TIMEOUT` - Timeout for metadata fetching (default: 10s)

## Error Handling

- Invalid URLs should be rejected with a clear message
- Failed metadata fetches should not prevent bookmark creation
- Database errors should return 500 with a generic message (log details server-side)

## Security Considerations

This is designed for personal/single-user use, but consider:
- Sanitize any HTML in fetched metadata before displaying
- Validate URLs before fetching (no internal network access)
- Rate limit the metadata fetch endpoint

## Future Ideas

Not for initial implementation, but keep in mind:
- Browser extension for one-click saving
- Import/export (Netscape bookmark format, JSON)
- Bookmark health checking (detect dead links)
- Reading list / archive mode
- Share individual bookmarks or collections
