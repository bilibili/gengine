package iparser

import parser "github.com/bilibili/gengine/internal/iantlr/alr"

type GengineParserVisitor struct {
	parser.BasegengineVisitor
}

func NewGengineParserVisitor() *GengineParserVisitor {
	return &GengineParserVisitor{}
}
