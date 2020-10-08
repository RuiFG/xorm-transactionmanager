package transaction_manager

import (
	"context"
	. "xorm.io/xorm"
)

type TransactionDefinition uint

const (
	_ TransactionDefinition = iota
	PROPAGATION_REQUIRED
	PROPAGATION_SUPPORTS
	PROPAGATION_MANDATORY
	PROPAGATION_REQUIRES_NEW
	PROPAGATION_NOT_SUPPORTED
	PROPAGATION_NEVER
	PROPAGATION_NESTED
)

//return session of the transaction manager
type DynamicSession func(...context.Context) (*Session, error)

//transaction function
type TransactionFunc func(context.Context, *Session) error

//transaction manager interface
type TransactionManager interface {
	//DynamicSessionFunc return a function
	DynamicSessionFunc() DynamicSession
	//IsInTransaction returns a value to determine whether it is currently in a transaction
	IsInTransaction(ctx context.Context) bool
	Do(context.Context, TransactionFunc, ...TransactionDefinition) error
	Required(context.Context, TransactionFunc) error
	Supports(context.Context, TransactionFunc) error
	Mandatory(context.Context, TransactionFunc) error
	RequiresNew(context.Context, TransactionFunc) error
	NotSupported(context.Context, TransactionFunc) error
	Never(context.Context, TransactionFunc) error
	Nested(context.Context, TransactionFunc) error
}
