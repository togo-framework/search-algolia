# search-algolia — documentation

  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />

## Overview

Package algolia is an Algolia driver for togo full-text search.
Blank-import it and set SEARCH_DRIVER=algolia, ALGOLIA_APP_ID, ALGOLIA_API_KEY.

## Install

```bash
togo install togo-framework/search-algolia
```

Set `SEARCH_DRIVER=algolia`.

## Configuration

Environment variables read by this plugin (extracted from the source):

| Env var | Notes |
|---|---|
| `ALGOLIA_API_KEY` | _see provider docs_ |
| `ALGOLIA_APP_ID` | _see provider docs_ |
| `G` | _see provider docs_ |

## Usage

```go
s := k.Search
s.Index(ctx, "posts", doc)
hits, _ := s.Search(ctx, "posts", "query")
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/search-algolia
- README: ../README.md
