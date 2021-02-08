package base

type AssignmentHolder interface {
	AcceptAssignment(a *Assignment) error
}
