package base

type StatementsHolder interface {
	AcceptStatements(statement *Statements) error
}
