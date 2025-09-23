# Poetry Archive Website

A Go web server for displaying and searching through a collection of poems.

## Features

- Individual poem pages with clean layout
- Search functionality across poem content
- Poetry listing page with condensed spacing
- Minimal design with transparent sidebar
- Background image with 3% opacity overlay
- Responsive navigation and search bar

## Local Development

```bash
# Run the server locally
go run main.go data.go

# Or use air for hot reloading
air
```

## Deployment

This project is configured for deployment on Vercel with the following files:

- `vercel.json` - Vercel configuration
- `go.mod` - Go module configuration
- `build.sh` - Build script

## File Structure

- `main.go` - Main Go web server
- `data.go` - Legacy data (can be removed)
- `static/poems/` - Poem JSON files
- `static/archbgs-01.webp` - Background image
