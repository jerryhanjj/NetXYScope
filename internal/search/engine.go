package search

import (
	"encoding/xml"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jerryhanjj/NetXYScope/internal/models"
)

// Engine 搜索引擎
type Engine struct{}

// NewEngine 创建新的搜索引擎
func NewEngine() *Engine {
	return &Engine{}
}

// SearchFiles 在指定目录中搜索NETCONF文件
func (e *Engine) SearchFiles(directory, searchTerm string) ([]models.SearchResult, error) {
	var results []models.SearchResult

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext == ".xml" || ext == ".yang" || strings.HasSuffix(path, ".yin") {
			fileResults, err := e.searchFile(path, searchTerm)
			if err != nil {
				return err
			}
			results = append(results, fileResults...)
		}

		return nil
	})

	return results, err
}

// searchFile 在单个文件中搜索
func (e *Engine) searchFile(filePath, searchTerm string) ([]models.SearchResult, error) {
	var results []models.SearchResult

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// For XML files, try to parse and search with better context
	if strings.HasSuffix(filePath, ".xml") {
		xmlResults, err := e.searchXMLFile(content, filePath, searchTerm)
		if err == nil {
			results = append(results, xmlResults...)
		}
	}

	// Fallback to line-based search for all files
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(searchTerm)) {
			matchType := "text"
			if strings.HasSuffix(filePath, ".xml") {
				matchType = "xml"
			} else if strings.HasSuffix(filePath, ".yang") || strings.HasSuffix(filePath, ".yin") {
				matchType = "yang"
			}

			// Avoid duplicate results from XML parsing
			if !containsResult(results, filePath, i+1) {
				results = append(results, models.SearchResult{
					FilePath:    filePath,
					LineNumber:  i + 1,
					LineContent: strings.TrimSpace(line),
					MatchType:   matchType,
				})
			}
		}
	}

	return results, nil
}

// searchXMLFile 在XML文件中搜索
func (e *Engine) searchXMLFile(content []byte, filePath, searchTerm string) ([]models.SearchResult, error) {
	var results []models.SearchResult
	var config models.NETCONFConfig

	err := xml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	// Search in XML content with better context
	if strings.Contains(strings.ToLower(config.Content), strings.ToLower(searchTerm)) {
		// For simplicity, we'll use line-based search within XML content
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), strings.ToLower(searchTerm)) {
				results = append(results, models.SearchResult{
					FilePath:    filePath,
					LineNumber:  i + 1,
					LineContent: strings.TrimSpace(line),
					MatchType:   "xml",
				})
			}
		}
	}

	return results, nil
}

// containsResult 检查是否已包含相同的结果
func containsResult(results []models.SearchResult, filePath string, lineNumber int) bool {
	for _, result := range results {
		if result.FilePath == filePath && result.LineNumber == lineNumber {
			return true
		}
	}
	return false
}

// HighlightSearchTerm 高亮显示搜索词
func HighlightSearchTerm(content, searchTerm string) string {
	lowerContent := strings.ToLower(content)
	lowerTerm := strings.ToLower(searchTerm)

	if idx := strings.Index(lowerContent, lowerTerm); idx >= 0 {
		before := content[:idx]
		match := content[idx : idx+len(searchTerm)]
		after := content[idx+len(searchTerm):]
		return before + "\033[1;31m" + match + "\033[0m" + after
	}
	return content
}
