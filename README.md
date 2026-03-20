# FormCMS Go

A headless CMS framework in Go, rewritten from the [FormCMS](https://github.com/formcms/formcms) original C# implementation.

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

## Deployment

### Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `DOMAIN` | Your external domain name (e.g., `example.com`). If set, enables automatic HTTPS via `autocert`. | `""` |
| `PORT` | The port to listen on for HTTP. Ignored if `DOMAIN` is set. | `5000` |

### Automatic HTTPS (autocert)

FormCMS Go supports automatic TLS certificate provisioning via Let's Encrypt using `autocert`. To enable this:

1. Point your domain's DNS A/AAAA records to your server's IP.
2. Ensure ports `80` and `443` are open and not in use by other processes.
3. Run the application with the `DOMAIN` environment variable:

```bash
DOMAIN=yourdomain.com ./formcms-go
```

The server will automatically handle HTTP-to-HTTPS redirection and store certificates in a local `certs` directory.

## Project Structure

- `core/api`: HTTP handlers and routing.
- `core/descriptors`: Data models and schema definitions.
- `core/services`: Business logic and orchestration.
- `infrastructure/filestore`: File storage implementations (Local, S3).
- `infrastructure/relationdbdao`: Database abstraction layer (PostgreSQL, SQLite).
- `utils`: Shared utilities and data models.
