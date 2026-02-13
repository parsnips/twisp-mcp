package search

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

const (
	SchemaIndexName   = "schema"
	DocsIndexName     = "docs"
	ExamplesIndexName = "examples"
)

// IndexManager manages the three Bleve indexes.
type IndexManager struct {
	SchemaIndex   bleve.Index
	DocsIndex     bleve.Index
	ExamplesIndex bleve.Index
}

// Close closes all indexes.
func (im *IndexManager) Close() error {
	var errs []error
	if im.SchemaIndex != nil {
		if err := im.SchemaIndex.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if im.DocsIndex != nil {
		if err := im.DocsIndex.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if im.ExamplesIndex != nil {
		if err := im.ExamplesIndex.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors closing indexes: %v", errs)
	}
	return nil
}

// SearchSchema searches the schema index.
func (im *IndexManager) SearchSchema(queryStr string, kind string, limit int) ([]SchemaDocument, error) {
	query := bleve.NewQueryStringQuery(queryStr)

	search := bleve.NewSearchRequest(query)
	search.Size = limit
	search.Fields = []string{"*"}

	if kind != "" {
		kindQuery := bleve.NewTermQuery(kind)
		kindQuery.SetField("kind")
		boolQuery := bleve.NewBooleanQuery()
		boolQuery.AddMust(query, kindQuery)
		search = bleve.NewSearchRequest(boolQuery)
		search.Size = limit
		search.Fields = []string{"*"}
	}

	results, err := im.SchemaIndex.Search(search)
	if err != nil {
		return nil, err
	}

	docs := make([]SchemaDocument, 0, len(results.Hits))
	for _, hit := range results.Hits {
		doc := SchemaDocument{
			ID: hit.ID,
		}
		if v, ok := hit.Fields["kind"].(string); ok {
			doc.Kind = v
		}
		if v, ok := hit.Fields["name"].(string); ok {
			doc.Name = v
		}
		if v, ok := hit.Fields["description"].(string); ok {
			doc.Description = v
		}
		if v, ok := hit.Fields["sourceFile"].(string); ok {
			doc.SourceFile = v
		}
		if v, ok := hit.Fields["definition"].(string); ok {
			doc.Definition = v
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// SearchDocs searches the docs index.
func (im *IndexManager) SearchDocs(queryStr string, category string, limit int) ([]DocDocument, error) {
	query := bleve.NewQueryStringQuery(queryStr)

	search := bleve.NewSearchRequest(query)
	search.Size = limit
	search.Fields = []string{"*"}

	if category != "" {
		catQuery := bleve.NewTermQuery(category)
		catQuery.SetField("category")
		boolQuery := bleve.NewBooleanQuery()
		boolQuery.AddMust(query, catQuery)
		search = bleve.NewSearchRequest(boolQuery)
		search.Size = limit
		search.Fields = []string{"*"}
	}

	results, err := im.DocsIndex.Search(search)
	if err != nil {
		return nil, err
	}

	docs := make([]DocDocument, 0, len(results.Hits))
	for _, hit := range results.Hits {
		doc := DocDocument{
			ID: hit.ID,
		}
		if v, ok := hit.Fields["title"].(string); ok {
			doc.Title = v
		}
		if v, ok := hit.Fields["description"].(string); ok {
			doc.Description = v
		}
		if v, ok := hit.Fields["path"].(string); ok {
			doc.Path = v
		}
		if v, ok := hit.Fields["category"].(string); ok {
			doc.Category = v
		}
		if v, ok := hit.Fields["content"].(string); ok {
			doc.Content = v
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// SearchExamples searches the examples index.
func (im *IndexManager) SearchExamples(queryStr string, category string, limit int) ([]ExampleDocument, error) {
	query := bleve.NewQueryStringQuery(queryStr)

	search := bleve.NewSearchRequest(query)
	search.Size = limit
	search.Fields = []string{"*"}

	if category != "" {
		catQuery := bleve.NewTermQuery(category)
		catQuery.SetField("category")
		boolQuery := bleve.NewBooleanQuery()
		boolQuery.AddMust(query, catQuery)
		search = bleve.NewSearchRequest(boolQuery)
		search.Size = limit
		search.Fields = []string{"*"}
	}

	results, err := im.ExamplesIndex.Search(search)
	if err != nil {
		return nil, err
	}

	docs := make([]ExampleDocument, 0, len(results.Hits))
	for _, hit := range results.Hits {
		doc := ExampleDocument{
			ID: hit.ID,
		}
		if v, ok := hit.Fields["category"].(string); ok {
			doc.Category = v
		}
		if v, ok := hit.Fields["subcategory"].(string); ok {
			doc.Subcategory = v
		}
		if v, ok := hit.Fields["name"].(string); ok {
			doc.Name = v
		}
		if v, ok := hit.Fields["query"].(string); ok {
			doc.Query = v
		}
		if v, ok := hit.Fields["variables"].(string); ok {
			doc.Variables = v
		}
		if v, ok := hit.Fields["instructions"].(string); ok {
			doc.Instructions = v
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// IndexSchemaDocument indexes a schema document.
func (im *IndexManager) IndexSchemaDocument(doc *SchemaDocument) error {
	return im.SchemaIndex.Index(doc.ID, doc)
}

// IndexDocDocument indexes a documentation document.
func (im *IndexManager) IndexDocDocument(doc *DocDocument) error {
	return im.DocsIndex.Index(doc.ID, doc)
}

// IndexExampleDocument indexes an example document.
func (im *IndexManager) IndexExampleDocument(doc *ExampleDocument) error {
	return im.ExamplesIndex.Index(doc.ID, doc)
}

// GetSchemaDocument retrieves a schema document by ID.
func (im *IndexManager) GetSchemaDocument(id string) (*SchemaDocument, error) {
	query := bleve.NewDocIDQuery([]string{id})
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}

	results, err := im.SchemaIndex.Search(search)
	if err != nil {
		return nil, err
	}
	if len(results.Hits) == 0 {
		return nil, fmt.Errorf("document not found: %s", id)
	}

	hit := results.Hits[0]
	result := &SchemaDocument{ID: id}
	if v, ok := hit.Fields["kind"].(string); ok {
		result.Kind = v
	}
	if v, ok := hit.Fields["name"].(string); ok {
		result.Name = v
	}
	if v, ok := hit.Fields["description"].(string); ok {
		result.Description = v
	}
	if v, ok := hit.Fields["sourceFile"].(string); ok {
		result.SourceFile = v
	}
	if v, ok := hit.Fields["definition"].(string); ok {
		result.Definition = v
	}
	return result, nil
}

// GetDocDocument retrieves a doc document by ID.
func (im *IndexManager) GetDocDocument(id string) (*DocDocument, error) {
	query := bleve.NewDocIDQuery([]string{id})
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}

	results, err := im.DocsIndex.Search(search)
	if err != nil {
		return nil, err
	}
	if len(results.Hits) == 0 {
		return nil, fmt.Errorf("document not found: %s", id)
	}

	hit := results.Hits[0]
	result := &DocDocument{ID: id}
	if v, ok := hit.Fields["title"].(string); ok {
		result.Title = v
	}
	if v, ok := hit.Fields["description"].(string); ok {
		result.Description = v
	}
	if v, ok := hit.Fields["path"].(string); ok {
		result.Path = v
	}
	if v, ok := hit.Fields["category"].(string); ok {
		result.Category = v
	}
	if v, ok := hit.Fields["content"].(string); ok {
		result.Content = v
	}
	return result, nil
}

// GetExampleDocument retrieves an example document by ID.
func (im *IndexManager) GetExampleDocument(id string) (*ExampleDocument, error) {
	query := bleve.NewDocIDQuery([]string{id})
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}

	results, err := im.ExamplesIndex.Search(search)
	if err != nil {
		return nil, err
	}
	if len(results.Hits) == 0 {
		return nil, fmt.Errorf("document not found: %s", id)
	}

	hit := results.Hits[0]
	result := &ExampleDocument{ID: id}
	if v, ok := hit.Fields["category"].(string); ok {
		result.Category = v
	}
	if v, ok := hit.Fields["subcategory"].(string); ok {
		result.Subcategory = v
	}
	if v, ok := hit.Fields["name"].(string); ok {
		result.Name = v
	}
	if v, ok := hit.Fields["query"].(string); ok {
		result.Query = v
	}
	if v, ok := hit.Fields["variables"].(string); ok {
		result.Variables = v
	}
	if v, ok := hit.Fields["instructions"].(string); ok {
		result.Instructions = v
	}
	return result, nil
}

// SchemaStats returns the document count in the schema index.
func (im *IndexManager) SchemaStats() (uint64, error) {
	return im.SchemaIndex.DocCount()
}

// DocsStats returns the document count in the docs index.
func (im *IndexManager) DocsStats() (uint64, error) {
	return im.DocsIndex.DocCount()
}

// ExamplesStats returns the document count in the examples index.
func (im *IndexManager) ExamplesStats() (uint64, error) {
	return im.ExamplesIndex.DocCount()
}

func buildSchemaMapping() mapping.IndexMapping {
	schemaMapping := bleve.NewDocumentMapping()

	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = "en"

	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = "keyword"

	schemaMapping.AddFieldMappingsAt("id", keywordFieldMapping)
	schemaMapping.AddFieldMappingsAt("kind", keywordFieldMapping)
	schemaMapping.AddFieldMappingsAt("name", textFieldMapping)
	schemaMapping.AddFieldMappingsAt("description", textFieldMapping)
	schemaMapping.AddFieldMappingsAt("sourceFile", keywordFieldMapping)
	schemaMapping.AddFieldMappingsAt("definition", textFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = schemaMapping
	return indexMapping
}

func buildDocsMapping() mapping.IndexMapping {
	docMapping := bleve.NewDocumentMapping()

	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = "en"

	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = "keyword"

	docMapping.AddFieldMappingsAt("id", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("title", textFieldMapping)
	docMapping.AddFieldMappingsAt("description", textFieldMapping)
	docMapping.AddFieldMappingsAt("path", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("category", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("content", textFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = docMapping
	return indexMapping
}

func buildExamplesMapping() mapping.IndexMapping {
	exampleMapping := bleve.NewDocumentMapping()

	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = "en"

	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = "keyword"

	exampleMapping.AddFieldMappingsAt("id", keywordFieldMapping)
	exampleMapping.AddFieldMappingsAt("category", keywordFieldMapping)
	exampleMapping.AddFieldMappingsAt("subcategory", keywordFieldMapping)
	exampleMapping.AddFieldMappingsAt("name", textFieldMapping)
	exampleMapping.AddFieldMappingsAt("query", textFieldMapping)
	exampleMapping.AddFieldMappingsAt("variables", textFieldMapping)
	exampleMapping.AddFieldMappingsAt("instructions", textFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = exampleMapping
	return indexMapping
}
