package base

type SourceCode struct {
	Code     string //raw code            -> ctx.GetText()
	LineNum  int    //line number         -> ctx.GetStart().GetLine()
	Column   int    //line start location -> ctx.GetStart().GetColumn()
	LineStop int    //line end location   -> ctx.GetStop().GetColumn()
}
