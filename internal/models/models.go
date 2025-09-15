package models

import "encoding/xml"

// SearchResult 表示一个搜索结果
type SearchResult struct {
	FilePath    string
	LineNumber  int
	LineContent string
	MatchType   string
}

// GenericXMLContent 表示通用XML内容结构
type GenericXMLContent struct {
	XMLName xml.Name `xml:""`          // 匹配任何根元素
	Content string   `xml:",innerxml"` // 获取所有内部XML内容
}

// NETCONFConfig 保持向后兼容性（已废弃，请使用 GenericXMLContent）
// Deprecated: Use GenericXMLContent instead
type NETCONFConfig = GenericXMLContent
