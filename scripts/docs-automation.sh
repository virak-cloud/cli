#!/bin/bash

# Documentation Automation Script for Virak CLI
# This script provides automated documentation tasks for local development and CI

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

# Default values
SKIP_VALIDATION=false
SKIP_LINK_CHECK=false
SKIP_COVERAGE_CHECK=false
GENERATE_COVERAGE_REPORT=false
DEPLOY_DOCS=false
OUTPUT_FORMAT="table"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --skip-validation)
      SKIP_VALIDATION=true
      shift
      ;;
    --skip-link-check)
      SKIP_LINK_CHECK=true
      shift
      ;;
    --skip-coverage-check)
      SKIP_COVERAGE_CHECK=true
      shift
      ;;
    --generate-coverage-report)
      GENERATE_COVERAGE_REPORT=true
      shift
      ;;
    --deploy)
      DEPLOY_DOCS=true
      shift
      ;;
    --output-format)
      OUTPUT_FORMAT="$2"
      shift 2
      ;;
    --help|-h)
      echo "Usage: $0 [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  --skip-validation          Skip markdown validation"
      echo "  --skip-link-check          Skip broken link checking"
      echo "  --skip-coverage-check      Skip documentation coverage check"
      echo "  --generate-coverage-report Generate a detailed coverage report"
      echo "  --deploy                   Deploy documentation to GitHub Pages"
      echo "  --output-format FORMAT     Output format for coverage report (table, json)"
      echo "  --help, -h                 Show this help message"
      echo ""
      echo "Examples:"
      echo "  $0                         # Run all documentation tasks"
      echo "  $0 --skip-validation       # Skip validation"
      echo "  $0 --generate-coverage-report --output-format json"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      echo "Use --help for usage information"
      exit 1
      ;;
  esac
done

echo -e "${BLUE}Virak CLI Documentation Automation${NC}"
echo "======================================"

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
    for tool in go jq markdownlint markdown-link-check; do
        if ! command -v $tool &> /dev/null; then
            print_error "$tool is not installed. Please install it first."
            exit 1
        fi
    done
    
    print_status "All dependencies are installed."
}

# Function to validate markdown files
validate_markdown() {
    if [ "$SKIP_VALIDATION" = true ]; then
        print_warning "Skipping markdown validation"
        return
    fi
    
    print_status "Validating markdown files..."
    
    # Find all markdown files
    MD_FILES=$(find "$DOCS_DIR" -name "*.md" | sort)
    
    # Create temporary file for validation results
    VALIDATION_RESULTS=$(mktemp)
    
    for md_file in $MD_FILES; do
        # Run markdownlint
        if ! markdownlint "$md_file" 2>> "$VALIDATION_RESULTS"; then
            echo "Markdownlint issues found in $md_file" >> "$VALIDATION_RESULTS"
        fi
        
        # Check for common markdown issues
        if grep -q " \t" "$md_file"; then
            echo "Found tabs in $md_file - consider using spaces" >> "$VALIDATION_RESULTS"
        fi
        
        # Check for trailing whitespace
        if grep -q " $" "$md_file"; then
            echo "Found trailing whitespace in $md_file" >> "$VALIDATION_RESULTS"
        fi
    done
    
    # Report validation results
    if [ -s "$VALIDATION_RESULTS" ]; then
        print_error "Markdown validation failed:"
        cat "$VALIDATION_RESULTS"
        rm -f "$VALIDATION_RESULTS"
        exit 1
    else
        print_status "Markdown validation passed."
        rm -f "$VALIDATION_RESULTS"
    fi
}

# Function to check for broken links
check_links() {
    if [ "$SKIP_LINK_CHECK" = true ]; then
        print_warning "Skipping link check"
        return
    fi
    
    print_status "Checking for broken links..."
    
    # Create temporary file for link check results
    LINK_CHECK_RESULTS=$(mktemp)
    
    # Run markdown-link-check on all markdown files
    find "$DOCS_DIR" -name "*.md" -exec markdown-link-check {} \; 2>> "$LINK_CHECK_RESULTS" || true
    
    # Report link check results
    if grep -q "ERROR" "$LINK_CHECK_RESULTS"; then
        print_error "Link check failed:"
        cat "$LINK_CHECK_RESULTS"
        rm -f "$LINK_CHECK_RESULTS"
        exit 1
    else
        print_status "Link check passed."
        rm -f "$LINK_CHECK_RESULTS"
    fi
}

# Function to check documentation coverage
check_coverage() {
    if [ "$SKIP_COVERAGE_CHECK" = true ]; then
        print_warning "Skipping coverage check"
        return
    fi
    
    print_status "Checking documentation coverage..."
    
    # Check which commands have documentation
    documented_cmds=$(find "$DOCS_DIR/reference" -name "*.md" -exec basename {} .md \; 2>/dev/null | sort)
    all_cmds=$(find "$PROJECT_ROOT/cmd" -maxdepth 1 -type d -not -path "$PROJECT_ROOT/cmd" -exec basename {} \; 2>/dev/null | sort)
    
    # Create coverage report
    COVERAGE_REPORT=$(mktemp)
    
    echo "Documentation Coverage Report:" > "$COVERAGE_REPORT"
    echo "============================" >> "$COVERAGE_REPORT"
    echo "" >> "$COVERAGE_REPORT"
    
    documented_count=0
    total_count=0
    
    for cmd in $all_cmds; do
        total_count=$((total_count + 1))
        if echo "$documented_cmds" | grep -q "^$cmd$"; then
            echo "✓ $cmd - Documented" >> "$COVERAGE_REPORT"
            documented_count=$((documented_count + 1))
        else
            echo "✗ $cmd - Missing documentation" >> "$COVERAGE_REPORT"
        fi
    done
    
    # Calculate coverage percentage
    if [ $total_count -gt 0 ]; then
        coverage=$((documented_count * 100 / total_count))
    else
        coverage=0
    fi
    
    echo "" >> "$COVERAGE_REPORT"
    echo "Coverage: $documented_count/$total_count ($coverage%)" >> "$COVERAGE_REPORT"
    
    # Generate detailed report if requested
    if [ "$GENERATE_COVERAGE_REPORT" = true ]; then
        echo "" >> "$COVERAGE_REPORT"
        echo "Detailed Coverage Analysis:" >> "$COVERAGE_REPORT"
        echo "============================" >> "$COVERAGE_REPORT"
        echo "" >> "$COVERAGE_REPORT"
        
        # Check for missing sections in documented commands
        for cmd in $documented_cmds; do
            cmd_doc="$DOCS_DIR/reference/$cmd.md"
            if [ -f "$cmd_doc" ]; then
                echo "Command: $cmd" >> "$COVERAGE_REPORT"
                
                # Check for required sections
                if grep -q "^## Overview" "$cmd_doc"; then
                    echo "  ✓ Overview section" >> "$COVERAGE_REPORT"
                else
                    echo "  ✗ Missing Overview section" >> "$COVERAGE_REPORT"
                fi
                
                if grep -q "^## Commands" "$cmd_doc"; then
                    echo "  ✓ Commands section" >> "$COVERAGE_REPORT"
                else
                    echo "  ✗ Missing Commands section" >> "$COVERAGE_REPORT"
                fi
                
                if grep -q "^## Examples" "$cmd_doc"; then
                    echo "  ✓ Examples section" >> "$COVERAGE_REPORT"
                else
                    echo "  ✗ Missing Examples section" >> "$COVERAGE_REPORT"
                fi
                
                echo "" >> "$COVERAGE_REPORT"
            fi
        done
    fi
    
    # Output coverage report
    if [ "$OUTPUT_FORMAT" = "json" ]; then
        # Convert to JSON format
        echo "{"
        echo "  \"total_commands\": $total_count,"
        echo "  \"documented_commands\": $documented_count,"
        echo "  \"coverage_percentage\": $coverage,"
        echo "  \"commands\": ["
        
        first=true
        for cmd in $all_cmds; do
            if [ "$first" = false ]; then
                echo ","
            fi
            first=false
            
            documented=false
            if echo "$documented_cmds" | grep -q "^$cmd$"; then
                documented=true
            fi
            
            echo "    {\"name\": \"$cmd\", \"documented\": $documented}"
        done
        
        echo ""
        echo "  ]"
        echo "}"
    else
        # Output table format
        cat "$COVERAGE_REPORT"
    fi
    
    rm -f "$COVERAGE_REPORT"
    
    # Check coverage threshold
    if [ $coverage -lt 80 ]; then
        print_warning "Documentation coverage is below 80% ($coverage%)"
        if [ "$SKIP_COVERAGE_CHECK" = false ]; then
            exit 1
        fi
    else
        print_status "Documentation coverage is good ($coverage%)"
    fi
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
            if [ -f "$DOCS_DIR/reference/$cmd_name.md" ]; then
                echo "- [$cmd_name Commands](reference/$cmd_name.md)" >> "$DOCS_DIR/SUMMARY.md"
            fi
        fi
    done
    
    cat >> "$DOCS_DIR/SUMMARY.md" << EOF

## Developer Documentation
- [Developer Guide Overview](developer/README.md)
- [Architecture Overview](developer/architecture.md)
- [API Documentation](developer/api/README.md)
- [Contributing Guide](developer/contributing.md)
- [Development Setup](developer/development-setup.md)

## Troubleshooting
- [Troubleshooting Overview](troubleshooting/README.md)
- [Common Issues](troubleshooting/common-issues.md)

EOF
    
    print_status "Table of contents generated."
}

# Function to build documentation for deployment
build_docs() {
    print_status "Building documentation for deployment..."
    
    # Create build directory
    BUILD_DIR="$PROJECT_ROOT/build/docs"
    mkdir -p "$BUILD_DIR"
    
    # Copy documentation files
    cp -r "$DOCS_DIR"/* "$BUILD_DIR/"
    
    # Add Jekyll front matter to markdown files
    find "$BUILD_DIR" -name "*.md" -exec sed -i '1i\---\nlayout: default\n---\n' {} \;
    
    # Create .nojekyll file
    touch "$BUILD_DIR/.nojekyll"
    
    print_status "Documentation built successfully."
    echo "Built documentation is available in: $BUILD_DIR"
}

# Function to deploy documentation to GitHub Pages
deploy_docs() {
    if [ "$DEPLOY_DOCS" = false ]; then
        return
    fi
    
    print_status "Deploying documentation to GitHub Pages..."
    
    # Check if we're on the main branch
    BRANCH=$(git rev-parse --abbrev-ref HEAD)
    if [ "$BRANCH" != "main" ] && [ "$BRANCH" != "master" ]; then
        print_warning "Not on main or master branch. Skipping deployment."
        return
    fi
    
    # Check if gh CLI is installed
    if ! command -v gh &> /dev/null; then
        print_warning "GitHub CLI (gh) is not installed. Skipping deployment."
        return
    fi
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_warning "Not in a git repository. Skipping deployment."
        return
    fi
    
    # Build documentation
    build_docs
    
    # Deploy to GitHub Pages
    cd "$PROJECT_ROOT"
    
    # Create a temporary branch for deployment
    DEPLOY_BRANCH="gh-pages-$(date +%s)"
    git checkout --orphan "$DEPLOY_BRANCH"
    
    # Copy built documentation
    rm -rf *
    cp -r build/docs/* .
    
    # Add and commit changes
    git add .
    git commit -m "Deploy documentation $(date)"
    
    # Push to GitHub Pages
    git push origin "$DEPLOY_BRANCH":gh-pages --force
    
    # Return to original branch
    git checkout "$BRANCH"
    
    # Delete temporary branch
    git branch -D "$DEPLOY_BRANCH"
    
    print_status "Documentation deployed to GitHub Pages."
}

# Main script logic
main() {
    # Check dependencies
    check_dependencies
    
    # Run documentation tasks
    validate_markdown
    check_links
    check_coverage
    generate_toc
    build_docs
    
    # Deploy if requested
    deploy_docs
    
    print_status "All documentation automation tasks completed successfully!"
}

# Run main function
main