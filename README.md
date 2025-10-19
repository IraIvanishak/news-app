# News Application

## Prerequisites

- Docker
- Docker Compose
- Go 1.22+

## Quick Start

### Production Setup

1. Clone the repository
```bash
git clone https://github.com/IraIvanishak/news-app.git
cd news-app
```

2. Create a `.env` file with necessary configurations

3. Start the application
```bash
make up
```


### Common Commands

- Start application: `make up`
- Stop application: `make down`
- Rebuild: `make rebuild`
- Run tests: `make test`
- View logs: `make logs`
- Access MongoDB CLI: `make mongo-cli`

### Testing

Run repository tests:
```bash
make test
```

### Database

Access MongoDB shell:
```bash
make mongo-cli
```

## Environment Configuration

Ensure `.env` file includes:
- `PUBLIC_PORT`
- `LISTEN_PORT`
- Other necessary environment variables

## Troubleshooting

- Ensure Docker is running
- Check Docker Compose version compatibility
- Verify `.env` file configuration

## License

[Your License Here]