package search

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"gopkg.in/yaml.v3"
)

// BuildIndexesFromFS builds all three in-memory Bleve indexes from an fs.FS
// rooted at the content directory (containing schema/, docs/, examples/ subdirs).
func BuildIndexesFromFS(contentFS fs.FS) (*IndexManager, error) {
	im := &IndexManager{}

	var err error

	// Create in-memory indexes
	im.SchemaIndex, err = bleve.NewMemOnly(buildSchemaMapping())
	if err != nil {
		return nil, fmt.Errorf("failed to create schema index: %w", err)
	}

	im.DocsIndex, err = bleve.NewMemOnly(buildDocsMapping())
	if err != nil {
		im.SchemaIndex.Close()
		return nil, fmt.Errorf("failed to create docs index: %w", err)
	}

	im.ExamplesIndex, err = bleve.NewMemOnly(buildExamplesMapping())
	if err != nil {
		im.SchemaIndex.Close()
		im.DocsIndex.Close()
		return nil, fmt.Errorf("failed to create examples index: %w", err)
	}

	// Index schema
	schemaFS, err := fs.Sub(contentFS, "schema")
	if err != nil {
		im.Close()
		return nil, fmt.Errorf("failed to access schema dir: %w", err)
	}
	if err := indexSchemaFromFS(im, schemaFS); err != nil {
		im.Close()
		return nil, fmt.Errorf("failed to index schema: %w", err)
	}

	// Index docs
	docsFS, err := fs.Sub(contentFS, "docs")
	if err != nil {
		im.Close()
		return nil, fmt.Errorf("failed to access docs dir: %w", err)
	}
	if err := indexDocsFromFS(im, docsFS); err != nil {
		im.Close()
		return nil, fmt.Errorf("failed to index docs: %w", err)
	}

	// Index examples
	examplesFS, err := fs.Sub(contentFS, "examples")
	if err != nil {
		im.Close()
		return nil, fmt.Errorf("failed to access examples dir: %w", err)
	}
	if err := indexExamplesFromFS(im, examplesFS); err != nil {
		im.Close()
		return nil, fmt.Errorf("failed to index examples: %w", err)
	}

	return im, nil
}

func indexSchemaFromFS(im *IndexManager, schemaFS fs.FS) error {
	files, err := fs.Glob(schemaFS, "*.graphql")
	if err != nil {
		return err
	}

	var sources []*ast.Source
	for _, file := range files {
		content, err := fs.ReadFile(schemaFS, file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}
		sources = append(sources, &ast.Source{
			Name:  file,
			Input: string(content),
		})
	}

	schema, gqlErr := gqlparser.LoadSchema(sources...)
	if gqlErr != nil {
		log.Printf("Warning: schema parsing errors: %v", gqlErr)
	}

	if schema == nil {
		return fmt.Errorf("failed to parse schema")
	}

	// Index types
	for _, typeDef := range schema.Types {
		if strings.HasPrefix(typeDef.Name, "__") {
			continue
		}

		doc := &SchemaDocument{
			ID:          fmt.Sprintf("type:%s", typeDef.Name),
			Kind:        string(typeDef.Kind),
			Name:        typeDef.Name,
			Description: typeDef.Description,
			SourceFile:  typeDef.Position.Src.Name,
			Definition:  formatTypeDefinition(typeDef),
		}

		for _, field := range typeDef.Fields {
			doc.Fields = append(doc.Fields, field.Name)
		}

		if err := im.IndexSchemaDocument(doc); err != nil {
			log.Printf("Warning: failed to index type %s: %v", typeDef.Name, err)
		}
	}

	// Index directives
	for _, directive := range schema.Directives {
		if strings.HasPrefix(directive.Name, "__") {
			continue
		}

		doc := &SchemaDocument{
			ID:          fmt.Sprintf("directive:%s", directive.Name),
			Kind:        "directive",
			Name:        directive.Name,
			Description: directive.Description,
			Definition:  formatDirectiveDefinition(directive),
		}

		for _, arg := range directive.Arguments {
			doc.Arguments = append(doc.Arguments, arg.Name)
		}

		if err := im.IndexSchemaDocument(doc); err != nil {
			log.Printf("Warning: failed to index directive %s: %v", directive.Name, err)
		}
	}

	return nil
}

func indexDocsFromFS(im *IndexManager, docsFS fs.FS) error {
	return fs.WalkDir(docsFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(p, ".md") {
			return nil
		}

		if strings.Contains(p, "_template") {
			return nil
		}

		content, err := fs.ReadFile(docsFS, p)
		if err != nil {
			return err
		}

		doc := parseMarkdownDoc(string(content), p)
		if doc == nil {
			return nil
		}

		if err := im.IndexDocDocument(doc); err != nil {
			log.Printf("Warning: failed to index doc %s: %v", p, err)
		}

		return nil
	})
}

// parseMarkdownDoc parses a markdown file into a DocDocument.
// p is the fs.FS-relative path (forward slashes).
func parseMarkdownDoc(content, p string) *DocDocument {
	var frontMatter map[string]interface{}
	var body string

	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content[3:], "---", 2)
		if len(parts) == 2 {
			if err := yaml.Unmarshal([]byte(parts[0]), &frontMatter); err == nil {
				body = strings.TrimSpace(parts[1])
			} else {
				body = content
			}
		} else {
			body = content
		}
	} else {
		body = content
	}

	// Determine category from path (fs.FS uses forward slashes)
	relPath := strings.TrimSuffix(p, ".md")
	parts := strings.Split(relPath, "/")

	var category string
	if len(parts) > 0 {
		category = parts[0]
	}

	title := ""
	if t, ok := frontMatter["title"].(string); ok {
		title = t
	} else {
		lines := strings.Split(body, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "# ") {
				title = strings.TrimPrefix(line, "# ")
				break
			}
		}
	}

	description := ""
	if d, ok := frontMatter["description"].(string); ok {
		description = d
	}

	var headers []string
	headerRegex := regexp.MustCompile(`^#{2,3}\s+(.+)$`)
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		if matches := headerRegex.FindStringSubmatch(scanner.Text()); len(matches) > 1 {
			headers = append(headers, matches[1])
		}
	}

	return &DocDocument{
		ID:          relPath,
		Title:       title,
		Description: description,
		Path:        relPath,
		Category:    category,
		Content:     body,
		Headers:     headers,
	}
}

func indexExamplesFromFS(im *IndexManager, examplesFS fs.FS) error {
	entries, err := fs.ReadDir(examplesFS, ".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		category := entry.Name()

		// Skip non-example directories
		if category == "regression" || category == "seed" {
			continue
		}

		if err := indexExampleCategoryFromFS(im, examplesFS, category); err != nil {
			log.Printf("Warning: error indexing category %s: %v", category, err)
		}
	}

	return nil
}

func indexExampleCategoryFromFS(im *IndexManager, examplesFS fs.FS, category string) error {
	return fs.WalkDir(examplesFS, category, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() != "request.gql" {
			return nil
		}

		exampleDir := path.Dir(p)

		// Read request.gql
		query, err := fs.ReadFile(examplesFS, p)
		if err != nil {
			return nil
		}

		// Read optional files (instructions.txt is lowercase on disk)
		var variables, instructions string
		if data, err := fs.ReadFile(examplesFS, path.Join(exampleDir, "variables.json")); err == nil {
			variables = string(data)
		}
		if data, err := fs.ReadFile(examplesFS, path.Join(exampleDir, "instructions.txt")); err == nil {
			instructions = string(data)
		}

		// Determine subcategory from path parts
		parts := strings.Split(exampleDir, "/")
		subcategory := ""
		if len(parts) > 1 {
			subcategory = parts[1]
		}

		operations, types := extractQueryInfo(string(query))

		doc := &ExampleDocument{
			ID:           exampleDir,
			Category:     category,
			Subcategory:  subcategory,
			Name:         path.Base(exampleDir),
			Query:        string(query),
			Variables:    variables,
			Instructions: instructions,
			Operations:   operations,
			Types:        types,
		}

		if err := im.IndexExampleDocument(doc); err != nil {
			log.Printf("Warning: failed to index example %s: %v", exampleDir, err)
		}

		return nil
	})
}

func formatTypeDefinition(typeDef *ast.Definition) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %s", strings.ToLower(string(typeDef.Kind)), typeDef.Name))

	if len(typeDef.Interfaces) > 0 {
		sb.WriteString(" implements ")
		sb.WriteString(strings.Join(typeDef.Interfaces, " & "))
	}

	if len(typeDef.Fields) > 0 {
		sb.WriteString(" {\n")
		for _, field := range typeDef.Fields {
			if field.Description != "" {
				sb.WriteString(fmt.Sprintf("  # %s\n", field.Description))
			}
			sb.WriteString(fmt.Sprintf("  %s", field.Name))
			if len(field.Arguments) > 0 {
				sb.WriteString("(")
				args := make([]string, len(field.Arguments))
				for i, arg := range field.Arguments {
					args[i] = fmt.Sprintf("%s: %s", arg.Name, arg.Type.String())
				}
				sb.WriteString(strings.Join(args, ", "))
				sb.WriteString(")")
			}
			sb.WriteString(fmt.Sprintf(": %s\n", field.Type.String()))
		}
		sb.WriteString("}")
	}

	if len(typeDef.EnumValues) > 0 {
		sb.WriteString(" {\n")
		for _, val := range typeDef.EnumValues {
			sb.WriteString(fmt.Sprintf("  %s\n", val.Name))
		}
		sb.WriteString("}")
	}

	return sb.String()
}

func formatDirectiveDefinition(directive *ast.DirectiveDefinition) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("directive @%s", directive.Name))

	if len(directive.Arguments) > 0 {
		sb.WriteString("(")
		args := make([]string, len(directive.Arguments))
		for i, arg := range directive.Arguments {
			args[i] = fmt.Sprintf("%s: %s", arg.Name, arg.Type.String())
		}
		sb.WriteString(strings.Join(args, ", "))
		sb.WriteString(")")
	}

	if len(directive.Locations) > 0 {
		locs := make([]string, len(directive.Locations))
		for i, loc := range directive.Locations {
			locs[i] = string(loc)
		}
		sb.WriteString(" on ")
		sb.WriteString(strings.Join(locs, " | "))
	}

	return sb.String()
}

func extractQueryInfo(query string) (operations, types []string) {
	opRegex := regexp.MustCompile(`(query|mutation|subscription)\s+(\w+)`)
	typeRegex := regexp.MustCompile(`\b([A-Z][a-zA-Z]+(?:Input|Type|Connection|Edge)?)\b`)

	for _, match := range opRegex.FindAllStringSubmatch(query, -1) {
		if len(match) > 2 {
			operations = append(operations, match[2])
		}
	}

	seen := make(map[string]bool)
	for _, match := range typeRegex.FindAllStringSubmatch(query, -1) {
		if len(match) > 1 && !seen[match[1]] {
			seen[match[1]] = true
			types = append(types, match[1])
		}
	}

	return
}
