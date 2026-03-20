# FormCMS Go

A headless CMS framework in Go, rewritten from the original C# implementation.

## Features

- **Dynamic Data Modeling**: Define entities and attributes via UI or API.
- **REST & GraphQL APIs**: Auto-generated CRUD and GraphQL endpoints.
- **File Storage**: Local and S3 support with image processing.
- **Social Engagement**: Built-in likes, bookmarks, and comments.
- **Embedded UI**: React Admin panel and GrapesJS page builder included.

## Getting Started

### Prerequisites

- Go 1.25+
- (Optional) Docker

### Running Locally

```bash
cd formcms-go
go run main.go
```

The server will start on `http://localhost:5000`.

### Building

```bash
cd formcms-go
go build -o formcms-go main.go
```

### Running with Docker

```bash
docker build -t formcms-go .
docker run -p 5000:5000 formcms-go
```

## Project Structure

- `core/api`: HTTP handlers and routing.
- `core/descriptors`: Data models and schema definitions.
- `core/services`: Business logic and orchestration.
- `infrastructure/filestore`: File storage implementations (Local, S3).
- `infrastructure/relationdbdao`: Database abstraction layer (PostgreSQL, SQLite).
- `utils`: Shared utilities and data models.
