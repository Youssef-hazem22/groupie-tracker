# Groupie Trackers

A Go web app that displays music artists and their concerts, using data from the Groupie Trackers API. Browse all artists, search them, and open an artist page to see members, first album, and every concert location with its dates.

## Features

- **Home page:** responsive grid of artist cards
- **Search:** by artist name, member, first album, creation year, or concert city
- **Artist page:** photo, name, active since, first album, members, and a concerts table with a show count per location
- **Bright, responsive design** that works on desktop and mobile
- **Safe by default:** uses `html/template`, so output is escaped automatically

## Tech Stack

- Go (standard library only, no external dependencies)
- `html/template` for server-side rendering
- Plain CSS (grid, flexbox, CSS variables)

## Project Structure

```
.
├── main.go          # Server setup and routes
├── handlers.go      # Data fetching, search, and HTTP handlers
├── objects.go       # Data structures
├── go.mod
└── templates/
    ├── index.html   # Home page (artist grid + search)
    ├── artist.html  # Artist detail page
    └── styling.css  # Styles for both pages
```

## Requirements

- Go 1.22 or newer (the routes use method patterns like `GET /{$}`)
- Internet access at startup (the app fetches data from the API)

## Run Locally

```bash
git clone <your-repo-url>
cd <your-project-folder>
go run .
```

Open <http://localhost:8080>.

Templates are loaded once at startup, so restart the server after editing any `.html` file. CSS changes only need a browser refresh.

## Routes

| Route | Description |
|---|---|
| `GET /` | Home page. Add `?q=text` to search |
| `GET /artist?id=N` | Artist detail page |
| `GET /styling.css` | Stylesheet |

Errors: `400` for an invalid artist id, `404` for an unknown artist, `500` for a template error.

## Data Source

The app fetches artists and relations (locations with dates) from the Groupie Trackers API once at startup and keeps them in memory. If the API is unreachable at startup, the server exits with an error.

## Deployment

**Docker:**

```bash
docker build -t groupie .
docker run --rm -p 8080:8080 groupie
```

**Vercel / Render / Railway:** import the GitHub repo and deploy. The server reads the port from the `PORT` environment variable and defaults to `8080`.

## Author

Youssef Hazem
