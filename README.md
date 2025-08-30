# Galaplate CLI

Go REST API boilerplate generator with built-in project scaffolding and code generation.

## Features

- 🚀 **Multiple Templates**: API, Full-stack, and Microservice project templates
- 🔧 **Database Support**: PostgreSQL and MySQL with migrations
- 🎯 **Built-in Generators**: Models, Controllers, Jobs, DTOs via console commands
- 🌐 **Cross-platform**: Supports Linux, macOS, and Windows (amd64/arm64)
- ⚡ **Zero Dependencies**: Single binary installation

## Installation

### Quick Install (requires sudo for /usr/local/bin)

```bash
curl -sSL https://raw.githubusercontent.com/galaplate/cli/main/install.sh | sudo bash
```

### Install to Custom Directory

```bash
curl -sSL https://raw.githubusercontent.com/galaplate/cli/main/install.sh | bash -s -- -d ~/.local/bin
```

### Install Specific Version

```bash
curl -sSL https://raw.githubusercontent.com/galaplate/cli/main/install.sh | bash -s -- -v v0.0.1
```

### Manual Download

Download the binary for your platform from the [releases page](https://github.com/galaplate/cli/releases).

## Usage

### Create a New Project

```bash
# Basic API project
galaplate new my-api

# Full-stack project with MySQL
galaplate new my-app --template=full --db=mysql

# Microservice without git initialization  
galaplate new microservice --template=micro --no-git

# Custom module name
galaplate new my-project --module=github.com/myuser/my-project
```

### Available Templates

- **api** (default): REST API with database integration
- **full**: Full-stack application with frontend components
- **micro**: Lightweight microservice template

### Available Commands

```bash
galaplate new <name>      # Create new project
galaplate templates       # List available templates
galaplate version        # Show version info
galaplate help           # Show help message
```

### Built-in Project Generators

After creating a project, use the built-in console system:

```bash
cd my-project

# Database
go run main.go console db:up          # Run migrations
go run main.go console db:down        # Rollback migrations

# Code Generation
go run main.go console make:model User
go run main.go console make:controller UserController  
go run main.go console make:job EmailNotification
go run main.go console make:dto UserCreateRequest

# List all commands
go run main.go console list
```

## Development

### Prerequisites

- Go 1.22+

### Building

```bash
go build -o galaplate ./cmd/galaplate
```

### Running Tests

```bash
go test ./...
```

## CI/CD

The project uses GitHub Actions for automated releases:

- **Triggers**: On tag push (v*)
- **Platforms**: Linux, macOS, Windows (amd64/arm64)
- **Artifacts**: Cross-compiled binaries uploaded to GitHub releases

## License

MIT License
