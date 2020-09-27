package transaction_manager

import (
	"context"
	. "github.com/go-xorm/xorm"
)

//return session of the transaction manager
type DynamicSession func(...context.Context) (*Session, error)

//transaction function
type TransactionFunc func(context.Context, *Session) error

//transaction manager interface
type TransactionManager interface {
	DynamicSessionFunc() DynamicSession
	Do(context.Context, TransactionFunc) error
}
