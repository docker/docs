package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
)

// Config represents the configuration file structure
type Config struct {
	Repos []RepoConfig `json:"repos"`
}

type RepoConfig struct {
	Owner               string `json:"owner"`
	Repo                string `json:"repo"`
	ContentPath         string `json:"content_path"`
	TitlePrefix         string `json:"title_prefix"`
	DescriptionTemplate string `json:"description_template"`
	KeywordsBase        string `json:"keywords_base"`
	FetchLimit          int    `json:"fetch_limit"`
	IncludePrereleases  bool   `json:"include_prereleases"`
}

// GitHubRelease represents the release data from GitHub API
type GitHubRelease struct {
	TagName     string        `json:"tagName"`
	Name        string        `json:"name"`
	Body        string        `json:"body"`
	PublishedAt time.Time     `json:"publishedAt"`
	IsPrerelease bool         `json:"isPrerelease"`
	URL         string        `json:"url"`
	Author      GitHubAuthor  `json:"author"`
	Assets      []GitHubAsset `json:"assets"`
}

type GitHubAuthor struct {
	Login string `json:"login"`
}

type GitHubAsset struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	Size          int64  `json:"size"`
	DownloadCount int    `json:"downloadCount"`
}

// TemplateData represents the data passed to the markdown template
type TemplateData struct {
	TitlePrefix string
	Tag         string
	Description string
	Keywords    string
	Date        string
	URL         string
	IsPrerelease bool
	Author      string
	Body        string
	Assets      []Asset
	Binaries    []Binary
	Checksums   []Checksum
}

type Asset struct {
	Name          string
	URL           string
	Size          int64
	DownloadCount int
}

type Binary struct {
	Name      string
	URL       string
	Artifacts string
	SizeMB    string
}

type Checksum struct {
	Name   string
	URL    string
	SizeKB string
}

var (
	cleanMode   bool
	projectRoot string
	execDir     string
)

func init() {
	flag.BoolVar(&cleanMode, "clean", false, "Remove existing release notes before fetching")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <repo> [version]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Fetch release notes from GitHub repositories and generate Hugo markdown files.\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  repo       Repository name (e.g., 'buildx', 'compose')\n")
		fmt.Fprintf(os.Stderr, "  version    Optional specific version to fetch (e.g., 'v0.18.0')\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s buildx                 # Fetch all buildx releases\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s buildx v0.18.0         # Fetch specific version\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --clean buildx         # Clean and refetch all\n", os.Args[0])
	}
}

func main() {
	flag.Parse()

	// Determine executable directory and project root
	execPath, err := os.Executable()
	if err != nil {
		// Fall back to current directory logic
		cwd, err := os.Getwd()
		if err != nil {
			fatal("Failed to get current directory: %v", err)
		}
		execDir = cwd
		projectRoot = filepath.Join(cwd, "../..")
	} else {
		execDir = filepath.Dir(execPath)
		projectRoot = filepath.Join(execDir, "../..")
	}
	projectRoot, _ = filepath.Abs(projectRoot)
	execDir, _ = filepath.Abs(execDir)

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	repoFilter := args[0]
	versionFilter := ""
	if len(args) > 1 {
		versionFilter = args[1]
	}

	// Check if gh CLI is available
	if _, err := exec.LookPath("gh"); err != nil {
		fatal("gh CLI is required but not installed. Install from: https://cli.github.com/")
	}

	// Read config
	config, err := readConfig()
	if err != nil {
		fatal("Failed to read config: %v", err)
	}

	fmt.Println("\033[34mStarting release notes fetch...\033[0m")
	fmt.Println()

	// Process each repository
	for _, repo := range config.Repos {
		if repoFilter != "" && repo.Repo != repoFilter {
			continue
		}

		if err := processRepo(repo, versionFilter); err != nil {
			fmt.Fprintf(os.Stderr, "\033[33m  Warning: %v\033[0m\n", err)
		}
		fmt.Println()
	}

	fmt.Println("\033[32m✅ Release notes fetch complete\033[0m")
}

func readConfig() (*Config, error) {
	configPath := filepath.Join(execDir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &config, nil
}

func processRepo(repo RepoConfig, versionFilter string) error {
	fmt.Printf("\033[34mFetching releases for %s/%s...\033[0m\n", repo.Owner, repo.Repo)

	outputDir := filepath.Join(projectRoot, repo.ContentPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Clean existing files if requested
	if cleanMode {
		fmt.Println("\033[33m  Cleaning existing release notes...\033[0m")
		entries, err := os.ReadDir(outputDir)
		if err != nil {
			return fmt.Errorf("reading output directory: %w", err)
		}
		for _, entry := range entries {
			if !entry.IsDir() && entry.Name() != "_index.md" && strings.HasSuffix(entry.Name(), ".md") {
				filePath := filepath.Join(outputDir, entry.Name())
				if err := os.Remove(filePath); err != nil {
					fmt.Fprintf(os.Stderr, "\033[33m  Warning: failed to remove %s: %v\033[0m\n", entry.Name(), err)
				}
			}
		}
	}

	count := 0

	// Fetch specific version or all releases
	if versionFilter != "" {
		if err := fetchAndGenerateRelease(repo, versionFilter, outputDir); err != nil {
			return fmt.Errorf("fetching version %s: %w", versionFilter, err)
		}
		count++
	} else {
		releases, err := listReleases(repo)
		if err != nil {
			return fmt.Errorf("listing releases: %w", err)
		}

		for _, tag := range releases {
			sanitized := sanitizeVersion(tag)
			outputFile := filepath.Join(outputDir, sanitized+".md")

			// Skip if file exists and not in clean mode
			if !cleanMode {
				if _, err := os.Stat(outputFile); err == nil {
					continue
				}
			}

			if err := fetchAndGenerateRelease(repo, tag, outputDir); err != nil {
				fmt.Fprintf(os.Stderr, "\033[33m  Skipping %s: %v\033[0m\n", tag, err)
				continue
			}
			count++
		}
	}

	// Create or update _index.md
	indexFile := filepath.Join(outputDir, "_index.md")
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		if err := createIndexFile(repo, indexFile); err != nil {
			return fmt.Errorf("creating index file: %w", err)
		}
		fmt.Println("\033[32m  ✓ Created _index.md\033[0m")
	}

	fmt.Printf("\033[32mFetched %d releases for %s/%s\033[0m\n", count, repo.Owner, repo.Repo)
	return nil
}

func listReleases(repo RepoConfig) ([]string, error) {
	cmd := exec.Command("gh", "release", "list",
		"--repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Repo),
		"--limit", fmt.Sprintf("%d", repo.FetchLimit),
		"--json", "tagName,isPrerelease")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh release list failed: %w", err)
	}

	// Parse JSON output
	type ReleaseListItem struct {
		TagName      string `json:"tagName"`
		IsPrerelease bool   `json:"isPrerelease"`
	}

	var releases []ReleaseListItem
	if err := json.Unmarshal(output, &releases); err != nil {
		return nil, fmt.Errorf("parsing release list: %w", err)
	}

	var tags []string
	for _, release := range releases {
		// Filter out prereleases if not wanted
		if !repo.IncludePrereleases && release.IsPrerelease {
			continue
		}
		tags = append(tags, release.TagName)
	}

	return tags, nil
}

func fetchAndGenerateRelease(repo RepoConfig, tag string, outputDir string) error {
	release, err := fetchRelease(repo, tag)
	if err != nil {
		return err
	}

	data := buildTemplateData(repo, release)

	sanitized := sanitizeVersion(tag)
	outputFile := filepath.Join(outputDir, sanitized+".md")

	if err := generateMarkdown(data, outputFile); err != nil {
		return fmt.Errorf("generating markdown: %w", err)
	}

	fmt.Printf("\033[32m  ✓ Generated %s.md\033[0m\n", sanitized)
	return nil
}

func fetchRelease(repo RepoConfig, tag string) (*GitHubRelease, error) {
	cmd := exec.Command("gh", "release", "view", tag,
		"--repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Repo),
		"--json", "tagName,name,body,publishedAt,isPrerelease,url,author,assets")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh release view failed: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(output, &release); err != nil {
		return nil, fmt.Errorf("parsing release JSON: %w", err)
	}

	return &release, nil
}

func buildTemplateData(repo RepoConfig, release *GitHubRelease) *TemplateData {
	description := strings.ReplaceAll(repo.DescriptionTemplate, "{{version}}", release.TagName)
	keywords := fmt.Sprintf("%s, %s", repo.KeywordsBase, release.TagName)
	author := release.Author.Login
	if author == "" {
		author = "unknown"
	}

	data := &TemplateData{
		TitlePrefix:  repo.TitlePrefix,
		Tag:          release.TagName,
		Description:  description,
		Keywords:     keywords,
		Date:         release.PublishedAt.Format("2006-01-02"),
		URL:          release.URL,
		IsPrerelease: release.IsPrerelease,
		Author:       author,
		Body:         strings.TrimSpace(release.Body),
		Assets:       make([]Asset, 0),
		Binaries:     make([]Binary, 0),
		Checksums:    make([]Checksum, 0),
	}

	// Process assets in two passes:
	// Pass 1: Categorize all assets into binaries, checksums, and metadata files.
	//         Metadata files (SBOM, provenance, sigstore) are stored in a map keyed by
	//         the binary filename they relate to (e.g., "buildx-v0.31.0.linux-amd64").
	// Pass 2: Associate metadata files with their corresponding binaries by looking up
	//         each binary name in the metadata map and building artifact links.
	binaries := make([]Binary, 0)
	metadata := make(map[string][]GitHubAsset) // keyed by binary name

	for _, asset := range release.Assets {
		data.Assets = append(data.Assets, Asset{
			Name:          asset.Name,
			URL:           asset.URL,
			Size:          asset.Size,
			DownloadCount: asset.DownloadCount,
		})

		// Categorize assets
		if isChecksum(asset.Name) {
			data.Checksums = append(data.Checksums, Checksum{
				Name:   asset.Name,
				URL:    asset.URL,
				SizeKB: fmt.Sprintf("%.1f", float64(asset.Size)/1024.0),
			})
		} else if isMetadata(asset.Name) || isSigstore(asset.Name) {
			// Extract base binary name (remove .sbom.json, .provenance.json, .sigstore.json)
			baseName := strings.TrimSuffix(asset.Name, ".sbom.json")
			baseName = strings.TrimSuffix(baseName, ".provenance.json")
			baseName = strings.TrimSuffix(baseName, ".sigstore.json")
			metadata[baseName] = append(metadata[baseName], asset)
		} else {
			binaries = append(binaries, Binary{
				Name:   asset.Name,
				URL:    asset.URL,
				SizeMB: fmt.Sprintf("%.1f", float64(asset.Size)/1048576.0),
			})
		}
	}

	// Second pass: add artifacts links to binaries
	for i := range binaries {
		artifacts := make([]string, 0)
		if relatedMeta, ok := metadata[binaries[i].Name]; ok {
			for _, meta := range relatedMeta {
				var link string
				if strings.Contains(meta.Name, ".sbom.") {
					link = fmt.Sprintf("[`sbom.json`](%s)", meta.URL)
				} else if strings.Contains(meta.Name, ".provenance.") {
					link = fmt.Sprintf("[`provenance.json`](%s)", meta.URL)
				} else if strings.Contains(meta.Name, ".sigstore.") {
					link = fmt.Sprintf("[`sigstore.json`](%s)", meta.URL)
				}
				if link != "" {
					artifacts = append(artifacts, link)
				}
			}
		}
		if len(artifacts) > 0 {
			binaries[i].Artifacts = strings.Join(artifacts, " ")
		} else {
			binaries[i].Artifacts = "-"
		}
	}

	data.Binaries = binaries
	return data
}

func generateMarkdown(data *TemplateData, outputFile string) error {
	tmplPath := filepath.Join(execDir, "templates", "release-note.md.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}

func createIndexFile(repo RepoConfig, indexFile string) error {
	content := fmt.Sprintf(`---
title: %s release notes
linkTitle: Release notes
description: Release notes for %s
keywords: %s
weight: 120
type: release-note
---

Release notes for %s. Each release includes new features, bug fixes,
and improvements.

For the latest releases and downloads, visit the
[GitHub releases page](https://github.com/%s/%s/releases).
`, repo.TitlePrefix, repo.TitlePrefix, repo.KeywordsBase,
		repo.TitlePrefix, repo.Owner, repo.Repo)

	return os.WriteFile(indexFile, []byte(content), 0644)
}

func sanitizeVersion(version string) string {
	// Remove leading 'v' and replace problematic characters
	s := strings.TrimPrefix(version, "v")
	s = strings.ReplaceAll(s, "/", "-")
	return s
}

func isChecksum(name string) bool {
	checksumPattern := regexp.MustCompile(`(?i)(checksum|sha256|sha512|md5).*\.(txt|sum)$`)
	checksumNames := regexp.MustCompile(`(?i)^(SHA256SUMS|CHECKSUMS|checksums\.txt)$`)
	return checksumPattern.MatchString(name) || checksumNames.MatchString(name)
}

func isMetadata(name string) bool {
	metadataPattern := regexp.MustCompile(`(?i)\.(sbom|provenance|attestation)\.(json|jsonl)$`)
	return metadataPattern.MatchString(name)
}

func isSigstore(name string) bool {
	sigstorePattern := regexp.MustCompile(`(?i)\.sigstore\.json$`)
	return sigstorePattern.MatchString(name)
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[31mError: "+format+"\033[0m\n", args...)
	os.Exit(1)
}
