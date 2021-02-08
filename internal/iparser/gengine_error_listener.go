package iparser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"strconv"
)

type GengineErrorListener struct {
	antlr.ErrorListener
	GrammarErrors []string
}

func NewGengineErrorListener() *GengineErrorListener {
	return &GengineErrorListener{
		GrammarErrors: make([]string, 0),
	}
}

/**
syntax err check
*/
func (el *GengineErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	el.GrammarErrors = append(el.GrammarErrors, "line"+" "+strconv.Itoa(line)+":"+strconv.Itoa(column)+" "+msg)
}

func (el *GengineErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {

}
func (el *GengineErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {

}
func (el *GengineErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {

}
