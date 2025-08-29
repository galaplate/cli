# 🚀 Galaplate CLI

> A powerful Go REST API boilerplate generator that gets you from zero to production-ready API in seconds.

Galaplate CLI is the official command-line tool for bootstrapping Go REST API projects using the [Galaplate framework](https://github.com/galaplate/galaplate). It downloads the latest template from GitHub and automatically configures everything with your custom module name.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/galaplate/galaplate.svg)](https://github.com/galaplate/galaplate/releases)

## ✨ Features

- 🏗️ **Instant Project Setup** - Create production-ready Go APIs in seconds
- 🌐 **GitHub Integration** - Downloads latest templates directly from GitHub
- 🔧 **Smart Module Replacement** - Automatically replaces module names throughout the project
- 🎯 **Multiple Templates** - Choose from API-only, full-stack, or microservice templates
- 📦 **Database Support** - Built-in support for PostgreSQL and MySQL
- 🛠️ **Rich Generators** - Built-in code generators for models, controllers, jobs, and DTOs
- 🔐 **Authentication** - JWT-based auth system included
- 🚀 **Production Ready** - Includes Docker, migrations, validation, and more

## 📦 Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/galaplate/galaplate/main/install.sh | bash
```

### Custom Installation Directory

```bash
curl -sSL https://raw.githubusercontent.com/galaplate/galaplate/main/install.sh | bash -s -- -d ~/.local/bin
```

### Install Specific Version

```bash
curl -sSL https://raw.githubusercontent.com/galaplate/galaplate/main/install.sh | bash -s -- -v 1.0.0
```

### Manual Installation

1. Download the binary for your platform from the [releases page](https://github.com/galaplate/galaplate/releases)
2. Make it executable: `chmod +x galaplate`
3. Move to your PATH: `mv galaplate /usr/local/bin/`

### Verify Installation

```bash
galaplate version
```

## 🎯 Quick Start

Create your first Galaplate project:

```bash
# Create a new API project
galaplate new my-api

# Create with custom module name
galaplate new my-api --module=github.com/myuser/my-api

# Create full-stack project with MySQL
galaplate new my-app --template=full --db=mysql

# Create microservice
galaplate new my-service --template=micro
```

## 📚 Usage

### Basic Commands

```bash
galaplate new <project-name>     # Create a new project
galaplate templates             # List available templates  
galaplate version              # Show version information
galaplate help                 # Show help
```

### Project Creation Options

```bash
galaplate new my-project [options]

Options:
  --template=<name>    Template type (api, full, micro) [default: api]
  --db=<type>         Database type (postgres, mysql) [default: postgres]  
  --module=<name>     Go module name [default: project-name]
  --no-git          Skip git initialization
  --force           Overwrite existing directory
```

### Examples

```bash
# Simple API project
galaplate new blog-api

# Full-stack application with MySQL
galaplate new my-webapp --template=full --db=mysql

# Microservice with custom module
galaplate new user-service \
  --template=micro \
  --module=github.com/mycompany/user-service

# Skip git initialization  
galaplate new my-api --no-git

# Force overwrite existing directory
galaplate new my-api --force
```

## 🏗️ Available Templates

### 🔹 API Template (Default)
**Best for:** REST APIs, backend services, microservices

**Features:**
- HTTP server with routing
- Database integration (PostgreSQL/MySQL)
- JWT authentication
- Background job processing
- Database migrations
- Request validation
- Comprehensive testing setup
- Docker configuration

```bash
galaplate new my-api --template=api
```

### 🔹 Full Template  
**Best for:** Full-stack web applications

**Features:**
- Everything from API template
- Frontend template integration
- Static asset handling
- Server-side rendering support

```bash
galaplate new my-webapp --template=full
```

### 🔹 Micro Template
**Best for:** Minimal microservices, distributed systems

**Features:**
- Lightweight HTTP server
- Health check endpoints
- Service discovery ready
- Minimal dependencies
- Fast startup time

```bash
galaplate new my-service --template=micro
```

## 🛠️ After Project Creation

Once your project is created, follow these steps:

```bash
# Navigate to your project
cd my-api

# Copy environment file
cp .env.example .env

# Edit database configuration
nano .env

# Install dependencies
go mod tidy

# Create database and run migrations
go run main.go console db:create create_users_table
go run main.go console db:up

# Start the development server
go run main.go
```

## 🎨 Built-in Generators

Your Galaplate project comes with powerful code generators:

```bash
# Generate a new model
go run main.go console make:model User

# Generate a controller  
go run main.go console make:controller UserController

# Generate a background job
go run main.go console make:job EmailNotification

# Generate a DTO (Data Transfer Object)
go run main.go console make:dto UserCreateRequest

# List all available commands
go run main.go console list
```

## 🏗️ Project Structure

```
my-api/
├── main.go                 # Application entry point
├── core/                   # Galaplate core framework
├── app/
│   ├── controllers/        # HTTP controllers
│   ├── models/            # Database models  
│   ├── jobs/              # Background jobs
│   ├── dto/               # Data transfer objects
│   └── middleware/        # HTTP middleware
├── database/
│   └── migrations/        # Database migrations
├── routes/                # Route definitions
├── storage/               # File storage
├── docker-compose.yml     # Docker configuration
├── Dockerfile            # Container definition
├── .env.example         # Environment template
└── README.md            # Project documentation
```

## 🌟 What You Get

Every Galaplate project includes:

### 🔧 **Development Tools**
- Hot reload during development
- Comprehensive logging system
- Environment-based configuration
- Database migration system
- Robust error handling

### 🗄️ **Database Features**  
- Database abstraction layer
- Automatic model binding
- Query builder
- Connection pooling
- Migration system

### 🔐 **Security**
- JWT authentication
- Password hashing (bcrypt)
- Rate limiting
- CORS support
- Request validation

### 🧪 **Testing**
- Unit test examples
- Integration test setup
- Test database configuration
- Coverage reporting

### 🚀 **Production Ready**
- Docker configuration
- Health check endpoints
- Graceful shutdown
- Performance monitoring hooks
- Structured logging

## 🔧 Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/galaplate/galaplate.git
cd galaplate/cli

# Build the CLI
go build -o galaplate cmd/galaplate/main.go

# Install locally
mv galaplate /usr/local/bin/
```

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🤝 Support & Community

- 📖 **Documentation:** [https://galaplate.dev](https://galaplate.dev)
- 🐛 **Issues:** [GitHub Issues](https://github.com/galaplate/galaplate/issues)
- 💬 **Discussions:** [GitHub Discussions](https://github.com/galaplate/galaplate/discussions)
- 🗨️ **Discord:** [Join our community](https://discord.gg/galaplate)

## 📋 Requirements

- **Go 1.22+** - Required for building and running projects
- **Git** - For project initialization (optional with `--no-git`)
- **Docker** - Optional, for containerized development

## 🗺️ Roadmap

- [ ] **Templates** - Additional project templates (GraphQL, gRPC)
- [ ] **Databases** - MongoDB and Redis support  
- [ ] **Cloud** - AWS, GCP deployment templates
- [ ] **Frontend** - React, Vue.js integration templates
- [ ] **Monitoring** - Prometheus, Grafana configurations
- [ ] **CI/CD** - GitHub Actions, GitLab CI templates

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with ❤️ by the Galaplate team
- Inspired by Laravel Artisan and Ruby on Rails generators
- Special thanks to all [contributors](https://github.com/galaplate/galaplate/contributors)

---

<div align="center">

**[Website](https://galaplate.dev)** • **[Documentation](https://docs.galaplate.dev)** • **[Examples](https://github.com/galaplate/examples)** • **[Discord](https://discord.gg/galaplate)**

Made with ❤️ by the Galaplate team

</div>