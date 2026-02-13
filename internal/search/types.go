package search

// SchemaDocument represents an indexed GraphQL schema element.
type SchemaDocument struct {
	ID          string   `json:"id"`          // e.g., "type:Account", "mutation:createAccount"
	Kind        string   `json:"kind"`        // type, input, enum, mutation, query, scalar, directive
	Name        string   `json:"name"`        // Account, createAccount
	Description string   `json:"description"` // GraphQL description comment
	SourceFile  string   `json:"sourceFile"`  // account.graphql
	Definition  string   `json:"definition"`  // Full GraphQL definition text
	Fields      []string `json:"fields"`      // Field names (for types)
	Arguments   []string `json:"arguments"`   // Argument names (for mutations/queries)
}

// DocDocument represents an indexed documentation page.
type DocDocument struct {
	ID          string   `json:"id"`          // Path-based ID
	Title       string   `json:"title"`       // From YAML front matter
	Description string   `json:"description"` // From front matter
	Path        string   `json:"path"`        // File path
	Category    string   `json:"category"`    // introduction, accounting-core, guides, etc.
	Content     string   `json:"content"`     // Full markdown
	Headers     []string `json:"headers"`     // H2/H3 headers
}

// ExampleDocument represents an indexed example.
type ExampleDocument struct {
	ID           string   `json:"id"`           // Directory path
	Category     string   `json:"category"`     // examples, reference, fixtures
	Subcategory  string   `json:"subcategory"`  // basicTranCodeFlow, clientCRUD, etc.
	Name         string   `json:"name"`         // Directory name
	Query        string   `json:"query"`        // GraphQL query/mutation text
	Variables    string   `json:"variables"`    // Variables JSON if present
	Instructions string   `json:"instructions"` // Instructions.txt if present
	Operations   []string `json:"operations"`   // Extracted operation names
	Types        []string `json:"types"`        // Referenced type names
}

// SearchResult represents a generic search result.
type SearchResult struct {
	ID       string  `json:"id"`
	Score    float64 `json:"score"`
	Document any     `json:"document"`
}
