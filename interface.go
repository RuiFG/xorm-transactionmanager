package transaction_manager

import (
	"context"
	. "xorm.io/xorm"
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
	Do(context.Context, TransactionFunc) error
}
