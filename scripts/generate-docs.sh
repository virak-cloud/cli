#!/bin/bash

# Documentation Generation Script for Virak CLI
# This script helps generate and validate documentation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DOCS_DIR="$PROJECT_ROOT/docs"

echo -e "${BLUE}Virak CLI Documentation Generator${NC}"
echo "=================================="

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if required tools are installed
check_dependencies() {
    print_status "Checking dependencies..."
    
    # Check for required tools
    for tool in go jq; do
        if ! command -v $tool &> /dev/null; then
            print_error "$tool is not installed. Please install it first."
            exit 1
        fi
    done
    
    print_status "All dependencies are installed."
}

# Function to extract command documentation from Go files
extract_command_docs() {
    print_status "Extracting command documentation from Go files..."
    
    # Find all cmd directories
    CMD_DIRS=$(find "$PROJECT_ROOT/cmd" -type d -mindepth 1 | sort)
    
    for cmd_dir in $CMD_DIRS; do
        cmd_name=$(basename "$cmd_dir")
        doc_file="$cmd_dir/doc.go"
        
        if [ -f "$doc_file" ]; then
            print_status "Processing $cmd_name command documentation..."
            
            # Extract documentation content
            # This is a simplified extraction - you may need to enhance this
            # based on your actual doc.go format
            echo "# $cmd_name Command Reference" > "$DOCS/reference/$cmd_name.md.tmp"
            echo "" >> "$DOCS/reference/$cmd_name.md.tmp"
            
            # Extract package documentation
            grep -A 1000 "^// Package" "$doc_file" | grep -B 1000 "^package" | head -n -1 | sed 's/^\/\/ //' >> "$DOCS/reference/$cmd_name.md.tmp"
            
            print_status "Generated temporary docs for $cmd_name"
        fi
    done
}

# Function to validate markdown files
validate_markdown() {
    print_status "Validating markdown files..."
    
    # Find all markdown files
    MD_FILES=$(find "$DOCS_DIR" -name "*.md" | sort)
    
    for md_file in $MD_FILES; do
        # Check for common markdown issues
        if grep -q " \t" "$md_file"; then
            print_warning "Found tabs in $md_file - consider using spaces"
        fi
        
        # Check for trailing whitespace
        if grep -q " $" "$md_file"; then
            print_warning "Found trailing whitespace in $md_file"
        fi
        
        # Check for broken internal links
        # This is a simplified check - you may want to use a markdown linter
        internal_links=$(grep -o '\[.*\](docs/[^)]*)' "$md_file" || true)
        for link in $internal_links; do
            link_path=$(echo "$link" | sed 's/.*(\(.*\))/\1/')
            full_path="$PROJECT_ROOT/$link_path"
            if [ ! -f "$full_path" ]; then
                print_error "Broken link in $md_file: $link_path"
            fi
        done
    done
    
    print_status "Markdown validation completed."
}

# Function to generate table of contents
generate_toc() {
    print_status "Generating table of contents..."
    
    # Generate TOC for main documentation
    cat > "$DOCS_DIR/SUMMARY.md" << EOF
# Virak CLI Documentation Summary

## User Guide
- [User Guide Overview](user-guide/README.md)
- [Installation Guide](user-guide/installation.md)
- [Authentication Guide](user-guide/authentication.md)
- [Getting Started](user-guide/getting-started.md)
- [Configuration Guide](user-guide/configuration.md)
- [Tutorials](user-guide/tutorials/)

## Command Reference
- [Command Reference Overview](reference/README.md)
EOF

    # Add command references
    for cmd_dir in "$PROJECT_ROOT/cmd"/*; do
        if [ -d "$cmd_dir" ]; then
            cmd_name=$(basename "$cmd_dir")
            echo "- [$cmd_name Commands](reference/$cmd_name.md)" >> "$DOCS_DIR/SUMMARY.md"
        fi
    done
    
    cat >> "$DOCS_DIR/SUMMARY.md" << EOF

## Developer Documentation
- [Developer Guide Overview](developer/README.md)

## Troubleshooting
- [Troubleshooting Overview](troubleshooting/README.md)

EOF
    
    print_status "Table of contents generated."
}

# Function to check documentation coverage
check_coverage() {
    print_status "Checking documentation coverage..."
    
    # Check which commands have documentation
    documented_cmds=$(find "$DOCS_DIR/reference" -name "*.md" -exec basename {} .md \; | sort)
    all_cmds=$(find "$PROJECT_ROOT/cmd" -maxdepth 1 -type d -not -path "$PROJECT_ROOT/cmd" -exec basename {} \; | sort)
    
    echo "Documentation Coverage Report:"
    echo "============================"
    
    for cmd in $all_cmds; do
        if echo "$documented_cmds" | grep -q "^$cmd$"; then
            echo -e "${GREEN}✓${NC} $cmd - Documented"
        else
            echo -e "${RED}✗${NC} $cmd - Missing documentation"
        fi
    done
    
    # Calculate coverage percentage
    total_cmds=$(echo "$all_cmds" | wc -l)
    documented_count=$(echo "$documented_cmds" | wc -l)
    coverage=$((documented_count * 100 / total_cmds))
    
    echo ""
    echo "Coverage: $documented_count/$total_cmds ($coverage%)"
    
    if [ $coverage -lt 80 ]; then
        print_warning "Documentation coverage is below 80%"
    else
        print_status "Documentation coverage is good ($coverage%)"
    fi
}

# Function to create a documentation build script
create_build_script() {
    print_status "Creating documentation build script..."
    
    cat > "$PROJECT_ROOT/scripts/build-docs.sh" << 'EOF'
#!/bin/bash

# Build Documentation Script
# This script builds the documentation for deployment

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DOCS_DIR="$PROJECT_ROOT/docs"

echo "Building documentation..."

# Create build directory
BUILD_DIR="$PROJECT_ROOT/build/docs"
mkdir -p "$BUILD_DIR"

# Copy documentation files
cp -r "$DOCS_DIR"/* "$BUILD_DIR/"

# Process markdown files (add header/footer, etc.)
find "$BUILD_DIR" -name "*.md" -exec sed -i '1i\---\nlayout: default\n---\n' {} \;

# Generate index
echo "Documentation build complete!"
echo "Built documentation is available in: $BUILD_DIR"
EOF
    
    chmod +x "$PROJECT_ROOT/scripts/build-docs.sh"
    print_status "Build script created."
}

# Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -e, --extract       Extract documentation from Go files"
    echo "  -v, --validate      Validate markdown files"
    echo "  -t, --toc           Generate table of contents"
    echo "  -c, --coverage      Check documentation coverage"
    echo "  -b, --build         Create build script"
    echo "  -a, --all           Run all tasks"
    echo ""
    echo "Examples:"
    echo "  $0 --all            # Run all documentation tasks"
    echo "  $0 --extract        # Extract docs from Go files"
    echo "  $0 --validate       # Validate markdown files"
}

# Main script logic
main() {
    # Create scripts directory if it doesn't exist
    mkdir -p "$PROJECT_ROOT/scripts"
    
    # Parse command line arguments
    case "${1:-}" in
        -h|--help)
            show_help
            exit 0
            ;;
        -e|--extract)
            check_dependencies
            extract_command_docs
            ;;
        -v|--validate)
            validate_markdown
            ;;
        -t|--toc)
            generate_toc
            ;;
        -c|--coverage)
            check_coverage
            ;;
        -b|--build)
            create_build_script
            ;;
        -a|--all)
            check_dependencies
            extract_command_docs
            validate_markdown
            generate_toc
            check_coverage
            create_build_script
            print_status "All documentation tasks completed!"
            ;;
        "")
            print_error "No option specified. Use --help for usage information."
            exit 1
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"