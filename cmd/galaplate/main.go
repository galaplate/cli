package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var (
	version = "0.1.0-dev"
	commit  = "none"
	date    = "unknown"
)

type TemplateConfig struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Version      string            `json:"version"`
	Features     []string          `json:"features"`
	Dependencies map[string]string `json:"dependencies"`
	Structure    struct {
		Include []string `json:"include"`
		Exclude []string `json:"exclude"`
	} `json:"structure"`
	Replacements map[string]string `json:"replacements"`
	SetupSteps   []string          `json:"setup_steps"`
}

type ProjectOptions struct {
	Name     string
	Template string
	DbType   string
	Module   string
	NoGit    bool
	Force    bool
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "new":
		if err := handleNewProject(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		showVersion()
	case "templates":
		if err := listTemplates(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "--help", "-h":
		showHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		showHelp()
		os.Exit(1)
	}
}

func handleNewProject(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("project name is required")
	}

	opts := parseProjectArgs(args)

	fmt.Printf("🚀 Creating new Galaplate project: %s\n", opts.Name)
	fmt.Printf("📦 Template: %s\n", opts.Template)
	fmt.Printf("🗄️  Database: %s\n", opts.DbType)
	fmt.Printf("📝 Module: %s\n", opts.Module)
	fmt.Println()

	// Check if directory exists
	if _, err := os.Stat(opts.Name); err == nil && !opts.Force {
		return fmt.Errorf("directory '%s' already exists. Use --force to overwrite", opts.Name)
	}

	// Download from GitHub and create project
	return createProjectFromGitHub(opts)
}

func parseProjectArgs(args []string) *ProjectOptions {
	opts := &ProjectOptions{
		Name:     args[0],
		Template: "api",
		DbType:   "postgres",
		Module:   args[0],
		NoGit:    false,
		Force:    false,
	}

	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "--template=") {
			opts.Template = strings.TrimPrefix(arg, "--template=")
		} else if strings.HasPrefix(arg, "--db=") {
			opts.DbType = strings.TrimPrefix(arg, "--db=")
		} else if strings.HasPrefix(arg, "--module=") {
			opts.Module = strings.TrimPrefix(arg, "--module=")
		} else if arg == "--no-git" {
			opts.NoGit = true
		} else if arg == "--force" {
			opts.Force = true
		}
	}

	return opts
}

func createProjectFromGitHub(opts *ProjectOptions) error {
	// Download the latest release from GitHub
	fmt.Println("📥 Downloading Galaplate template from GitHub...")

	tempDir, err := downloadGitHubRepo("galaplate", "galaplate", "main")
	if err != nil {
		return fmt.Errorf("failed to download template: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("📁 Template downloaded to: %s\n", tempDir)

	// Create project directory
	if err := os.MkdirAll(opts.Name, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	// Copy files from downloaded template
	if err := copyProjectFiles(tempDir, opts.Name, opts); err != nil {
		return fmt.Errorf("failed to copy project files: %v", err)
	}

	// Setup project
	if err := setupProject(opts); err != nil {
		return fmt.Errorf("failed to setup project: %v", err)
	}

	fmt.Printf("\n✅ Project '%s' created successfully!\n\n", opts.Name)
	fmt.Println("🎯 Next steps:")
	fmt.Printf("   cd %s\n", opts.Name)
	fmt.Println("   cp .env.example .env")
	fmt.Println("   # Edit .env with your database settings")
	fmt.Println("   go mod tidy")
	fmt.Println("   go run main.go console db:create create_users_table")
	fmt.Println("   go run main.go console db:up")
	fmt.Println("   go run main.go")
	fmt.Println()
	fmt.Println("🛠️  Available generators:")
	fmt.Println("   go run main.go console make:model User")
	fmt.Println("   go run main.go console make:controller UserController")
	fmt.Println("   go run main.go console list  # See all commands")

	return nil
}

func downloadGitHubRepo(owner, repo, branch string) (string, error) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "galaplate-")
	if err != nil {
		return "", err
	}

	// Download the repository archive
	url := fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/%s.tar.gz", owner, repo, branch)
	resp, err := http.Get(url)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to download: HTTP %d", resp.StatusCode)
	}

	// Extract the tar.gz file
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.RemoveAll(tempDir)
			return "", err
		}

		// Skip the top-level directory (e.g., "galaplate-main/")
		parts := strings.Split(header.Name, "/")
		if len(parts) <= 1 {
			continue
		}

		// Reconstruct path without the top-level directory
		path := filepath.Join(tempDir, filepath.Join(parts[1:]...))

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				os.RemoveAll(tempDir)
				return "", err
			}
		case tar.TypeReg:
			// Create directory if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				os.RemoveAll(tempDir)
				return "", err
			}

			// Create file
			file, err := os.Create(path)
			if err != nil {
				os.RemoveAll(tempDir)
				return "", err
			}

			if _, err := io.Copy(file, tr); err != nil {
				file.Close()
				os.RemoveAll(tempDir)
				return "", err
			}
			file.Close()

			// Set file permissions
			if err := os.Chmod(path, os.FileMode(header.Mode)); err != nil {
				os.RemoveAll(tempDir)
				return "", err
			}
		}
	}

	return tempDir, nil
}

func copyProjectFiles(srcRoot, destDir string, opts *ProjectOptions) error {
	excludeDirs := []string{
		"cli", ".git", "node_modules", "tmp", "storage",
		"cmd", "test-api", // Exclude any test directories
		".opencode", // Exclude opencode metadata
	}

	excludeFiles := []string{
		"main",      // Exclude compiled binary
		"server",    // Exclude compiled server binary
		"galaplate", // Exclude compiled CLI binary
		"AGENTS.md", "CORE_EXTRACTION_PLAN.md", "IMPLEMENTATION_ROADMAP.md",
		"PHASE1_PROGRESS_REPORT.md", "WORKFLOW_EXAMPLE.md", // Exclude development docs
	}

	return filepath.Walk(srcRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from source root
		relPath, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}

		// Skip if it's the root directory
		if relPath == "." {
			return nil
		}

		// Check if we should exclude this path
		pathParts := strings.Split(relPath, string(filepath.Separator))
		for _, excludeDir := range excludeDirs {
			if pathParts[0] == excludeDir {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Check if we should exclude this file
		if slices.Contains(excludeFiles, info.Name()) {
			return nil
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		return copyFile(path, destPath, opts)
	})
}

func copyFile(src, dest string, opts *ProjectOptions) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// For template files, apply replacements
	if strings.HasSuffix(src, ".template") {
		// Remove .template extension from destination
		dest = strings.TrimSuffix(dest, ".template")
		destFile.Close()

		// Read template content
		content, err := io.ReadAll(srcFile)
		if err != nil {
			return err
		}

		// Apply template replacements
		contentStr := string(content)
		contentStr = strings.ReplaceAll(contentStr, "{{PROJECT_NAME}}", opts.Name)
		contentStr = strings.ReplaceAll(contentStr, "{{MODULE_NAME}}", opts.Module)
		contentStr = strings.ReplaceAll(contentStr, "{{DB_TYPE}}", opts.DbType)
		contentStr = strings.ReplaceAll(contentStr, "{{GALAPLATE_CORE_VERSION}}", "v0.0.0")

		// Write processed content
		return os.WriteFile(dest, []byte(contentStr), 0644)
	}

	// Special handling for go.mod and Go files
	if filepath.Base(src) == "go.mod" || strings.HasSuffix(src, ".go") {
		content, err := io.ReadAll(srcFile)
		if err != nil {
			return err
		}

		contentStr := string(content)

		if filepath.Base(src) == "go.mod" {
			// Update module name in go.mod
			contentStr = strings.ReplaceAll(contentStr, "module github.com/galaplate/galaplate", "module "+opts.Module)
		} else {
			// Replace galaplate module references in Go files with the new module name
			contentStr = replaceModuleReferences(contentStr, opts.Module)
		}

		destFile.Close()
		return os.WriteFile(dest, []byte(contentStr), 0644)
	}

	// For regular files, just copy
	_, err = io.Copy(destFile, srcFile)
	return err
}

func replaceModuleReferences(content, newModule string) string {
	// Replace import paths that reference github.com/galaplate/galaplate
	re := regexp.MustCompile(`"github\.com/galaplate/galaplate([^"]*)"`)
	content = re.ReplaceAllString(content, `"`+newModule+`$1"`)

	// Replace any other references to the old module
	content = strings.ReplaceAll(content, "github.com/galaplate/galaplate", newModule)

	return content
}

func setupProject(opts *ProjectOptions) error {
	projectDir := opts.Name

	fmt.Printf("🔧 Setting up project in %s...\n", projectDir)

	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(projectDir); err != nil {
		return err
	}

	// Initialize git repository (unless disabled)
	if !opts.NoGit {
		fmt.Println("📦 Initializing git repository...")
		if err := exec.Command("git", "init").Run(); err != nil {
			fmt.Println("⚠️  Warning: Failed to initialize git repository")
		} else {
			// Create initial commit
			exec.Command("git", "add", ".").Run()
			exec.Command("git", "commit", "-m", "Initial commit from Galaplate").Run()
		}
	}

	// Run go mod tidy to download dependencies
	fmt.Println("📥 Installing dependencies...")
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Println("⚠️  Warning: Failed to run 'go mod tidy'. Please run it manually.")
	}

	fmt.Println("✅ Project setup completed!")
	return nil
}

func showVersion() {
	fmt.Printf("Galaplate CLI v%s\n", version)
	fmt.Printf("Built: %s\n", date)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Println("galaplate core: v0.0.0-local")
}

func listTemplates() error {
	fmt.Println("📋 Available Templates:")
	fmt.Println()
	fmt.Println("  🔹 api     - REST API only (default)")
	fmt.Println("           Features: HTTP server, Database, Auth, Jobs, Console")
	fmt.Println("           Best for: Backend APIs, microservices")
	fmt.Println()
	fmt.Println("  🔹 full    - Full-stack application")
	fmt.Println("           Features: API + Frontend templates + Static assets")
	fmt.Println("           Best for: Web applications with UI")
	fmt.Println()
	fmt.Println("  🔹 micro   - Microservice template")
	fmt.Println("           Features: Minimal API, Service discovery, Health checks")
	fmt.Println("           Best for: Distributed systems, minimal services")
	fmt.Println()
	fmt.Println("Usage: galaplate new my-project --template=api")

	return nil
}

func showHelp() {
	fmt.Println(`🚀 Galaplate CLI - Go REST API Boilerplate Generator

Usage:
  galaplate <command> [arguments]

Commands:
  new <project-name>     Create a new Galaplate project
  version               Show CLI version information
  templates             List available project templates
  help                  Show this help message

Flags for 'new' command:
  --template=<name>     Template to use (api, full, micro) [default: api]
  --db=<type>          Database type (postgres, mysql) [default: postgres]
  --module=<name>      Custom Go module name [default: project-name]
  --no-git            Skip git initialization
  --force             Overwrite existing directory

Examples:
  galaplate new my-api                                    # Simple API project
  galaplate new my-app --template=full --db=mysql        # Full-stack with MySQL
  galaplate new microservice --template=micro --no-git   # Microservice without git

🛠️  After project creation:
  cd my-api
  cp .env.example .env        # Configure your environment
  go mod tidy                 # Install dependencies
  go run main.go console db:up  # Run database migrations
  go run main.go              # Start the server

🎯 Built-in Generators (available in your project):
  go run main.go console make:model User
  go run main.go console make:controller UserController  
  go run main.go console make:job EmailNotification
  go run main.go console make:dto UserCreateRequest
  go run main.go console list                    # See all available commands

💡 The CLI downloads project templates. All powerful generators are built into
   the project itself via galaplate core's console system.

For more information: https://github.com/galaplate/galaplate`)
}
