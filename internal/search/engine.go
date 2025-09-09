package search

import (
	"encoding/xml"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jerryhanjj/NetXYScope/internal/models"
)

// Engine 搜索引擎
type Engine struct{}

// NewEngine 创建新的搜索引擎
func NewEngine() *Engine {
	return &Engine{}
}

// SearchFiles 在指定目录中搜索NETCONF文件（并行版本）
func (e *Engine) SearchFiles(directory, searchTerm string) ([]models.SearchResult, error) {
	return e.SearchFilesParallel(directory, searchTerm, 8) // 默认8个worker
}

// SearchFilesParallel 并行搜索文件
func (e *Engine) SearchFilesParallel(directory, searchTerm string, numWorkers int) ([]models.SearchResult, error) {
	// 先收集所有NETCONF文件
	files, err := e.findNETCONFFiles(directory)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	// 创建worker pool
	fileCh := make(chan string, len(files))
	resultsCh := make(chan []models.SearchResult, len(files))
	var wg sync.WaitGroup

	// 添加文件到通道
	for _, file := range files {
		fileCh <- file
	}
	close(fileCh)

	// 启动worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileCh {
				fileResults, err := e.searchFile(file, searchTerm)
				if err == nil {
					resultsCh <- fileResults
				} else {
					resultsCh <- nil
				}
			}
		}()
	}

	// 等待所有worker完成
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// 收集结果
	var allResults []models.SearchResult
	for results := range resultsCh {
		if results != nil {
			allResults = append(allResults, results...)
		}
	}

	return allResults, nil
}

// findNETCONFFiles 查找所有NETCONF相关文件
func (e *Engine) findNETCONFFiles(directory string) ([]string, error) {
	var files []string
	
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext == ".xml" || ext == ".yang" || strings.HasSuffix(path, ".yin") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// searchFile 在单个文件中搜索（使用Boyer-Moore优化）
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

	// 使用Boyer-Moore算法进行高效搜索
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if e.boyerMooreSearch(strings.ToLower(line), strings.ToLower(searchTerm)) {
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

// boyerMooreSearch Boyer-Moore字符串搜索算法
func (e *Engine) boyerMooreSearch(text, pattern string) bool {
	if len(pattern) == 0 {
		return true
	}
	if len(text) == 0 || len(text) < len(pattern) {
		return false
	}

	// 构建坏字符表
	badCharTable := make(map[byte]int)
	for i := 0; i < len(pattern); i++ {
		badCharTable[pattern[i]] = len(pattern) - i - 1
	}

	i := len(pattern) - 1
	for i < len(text) {
		j := len(pattern) - 1
		
		// 从右向左匹配
		for j >= 0 && text[i] == pattern[j] {
			i--
			j--
		}
		
		// 完全匹配
		if j < 0 {
			return true
		}
		
		// 根据坏字符规则移动
		if shift, exists := badCharTable[text[i]]; exists {
			i += max(shift, len(pattern)-j)
		} else {
			i += len(pattern)
		}
	}
	
	return false
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
