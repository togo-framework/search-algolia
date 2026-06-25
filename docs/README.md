# search-algolia — documentation

Algolia driver for togo full-text search

## Overview

Package algolia is an Algolia driver for togo full-text search.
Blank-import it and set SEARCH_DRIVER=algolia, ALGOLIA_APP_ID, ALGOLIA_API_KEY.

## Install

```bash
togo install togo-framework/search-algolia
```

Set `SEARCH_DRIVER=algolia`.

## Configuration

Environment variables read by this plugin (extracted from the source — see the gateway/provider docs for each value):

| Env var |
|---|
| `ALGOLIA_API_KEY"` |
| `ALGOLIA_APP_ID"` |

## Usage

```go
s := k.Search
s.Index(ctx, "posts", doc)
hits, _ := s.Search(ctx, "posts", "query")
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/search-algolia
- Full README: ../README.md
