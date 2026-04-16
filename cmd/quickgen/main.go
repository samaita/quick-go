package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	schemaPath := flag.String("schema", "", "path to DDL SQL file (required)")
	outDir := flag.String("out", ".", "project root directory (default: current dir)")
	force := flag.Bool("force", false, "overwrite existing generated files")
	flag.Parse()

	if *schemaPath == "" {
		fmt.Fprintln(os.Stderr, "quickgen: --schema is required")
		flag.Usage()
		os.Exit(1)
	}

	tables, err := ParseDDL(*schemaPath)
	if err != nil {
		log.Fatalf("quickgen: parse DDL: %v", err)
	}
	log.Printf("quickgen: parsed %d table(s):", len(tables))
	for _, t := range tables {
		log.Printf("  - %s (%d columns)", t.Name, len(t.Columns))
	}

	tmplDir := filepath.Join(*outDir, "templates")

	// Write the shared handler helpers file (once, not per-table)
	if err := renderFile(tmplDir, *outDir, "handler_helpers.go.tmpl",
		filepath.Join("generated", "handler", "helpers.go"), nil, *force); err != nil {
		log.Fatalf("quickgen: handler helpers: %v", err)
	}

	// Generate per-table files
	for _, table := range tables {
		if err := renderFile(tmplDir, *outDir, "model.tmpl",
			filepath.Join("generated", "model", table.Name+".go"), table, *force); err != nil {
			log.Fatalf("quickgen: model %s: %v", table.Name, err)
		}
		if err := renderFile(tmplDir, *outDir, "repo.tmpl",
			filepath.Join("generated", "repo", table.Name+".go"), table, *force); err != nil {
			log.Fatalf("quickgen: repo %s: %v", table.Name, err)
		}
		if err := renderFile(tmplDir, *outDir, "handler.tmpl",
			filepath.Join("generated", "handler", table.Name+".go"), table, *force); err != nil {
			log.Fatalf("quickgen: handler %s: %v", table.Name, err)
		}

		// Hugo content pages
		if err := renderFile(tmplDir, *outDir, "hugo/list.tmpl",
			filepath.Join("frontend", "hugo-site", "content", table.URLPath(), "_index.md"), table, *force); err != nil {
			log.Fatalf("quickgen: hugo list %s: %v", table.Name, err)
		}
		if err := renderFile(tmplDir, *outDir, "hugo/form.tmpl",
			filepath.Join("frontend", "hugo-site", "content", table.URLPath(), "form.md"), table, *force); err != nil {
			log.Fatalf("quickgen: hugo form %s: %v", table.Name, err)
		}
	}

	// Generate routes.go (takes all tables)
	if err := renderFile(tmplDir, *outDir, "routes.tmpl",
		filepath.Join("generated", "routes.go"), tables, *force); err != nil {
		log.Fatalf("quickgen: routes: %v", err)
	}

	log.Println("quickgen: done. run 'go build ./...' to verify.")
}

func renderFile(tmplDir, outDir, tmplName, relOut string, data any, force bool) error {
	outPath := filepath.Join(outDir, relOut)

	if !force {
		if _, err := os.Stat(outPath); err == nil {
			log.Printf("quickgen: skip (exists, use --force to overwrite): %s", relOut)
			return nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(outPath), err)
	}

	tmplPath := filepath.Join(tmplDir, tmplName)
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("read template %s: %w", tmplPath, err)
	}

	tmpl, err := template.New(filepath.Base(tmplName)).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", tmplName, err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", outPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		os.Remove(outPath) // clean up partial write
		return fmt.Errorf("execute template %s → %s: %w", tmplName, relOut, err)
	}

	log.Printf("quickgen: wrote %s", relOut)
	return nil
}
