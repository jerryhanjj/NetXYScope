package models

import "encoding/xml"

// SearchResult 表示一个搜索结果
type SearchResult struct {
	FilePath    string
	LineNumber  int
	LineContent string
	MatchType   string
}

// NETCONFConfig 表示NETCONF配置结构
type NETCONFConfig struct {
	XMLName xml.Name `xml:"config"`
	Content string   `xml:",innerxml"`
}
