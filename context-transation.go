package transaction_manager

import (
	"context"
	"errors"
	. "xorm.io/xorm"
)

var (
	ErrContextParameter               = errors.New("context transaction manager need context")
	ErrNeverNotSupported              = errors.New("never transaction definition not supported transaction")
	ErrUnSupportTransactionDefinition = errors.New("unsupported transaction")
	ErrMandatoryInTransaction         = errors.New("mandatory needs to be in the transaction")
)

type ContextSessionKey struct{}

type contextTransactionManager struct {
	engine EngineInterface
}

func (ctm contextTransactionManager) DynamicSessionFunc() DynamicSession {
	return func(ctx ...context.Context) (*Session, error) {
		if len(ctx) == 0 {
			panic(ErrContextParameter)
		}
		sessionCtx := ctx[0]
		value := sessionCtx.Value(ContextSessionKey{})
		session, ok := value.(*Session)
		if ok {
			return session, nil
		}
		return ctm.engine.NewSession(), nil

	}
}

func (ctm contextTransactionManager) Do(ctx context.Context, transactionFunc TransactionFunc, propagations ...TransactionDefinition) error {
	propagation := PROPAGATION_REQUIRED
	if len(propagations) != 0 {
		propagation = propagations[0]
	}
	switch propagation {
	case PROPAGATION_REQUIRED:
		return ctm.Required(ctx, transactionFunc)
	case PROPAGATION_SUPPORTS:
		return ctm.Supports(ctx, transactionFunc)
	case PROPAGATION_NEVER:
		return ctm.Never(ctx, transactionFunc)
	case PROPAGATION_NOT_SUPPORTED:
		return ctm.NotSupported(ctx, transactionFunc)
	case PROPAGATION_NESTED:
		return ctm.Nested(ctx, transactionFunc)
	case PROPAGATION_MANDATORY:
		return ctm.Mandatory(ctx, transactionFunc)
	case PROPAGATION_REQUIRES_NEW:
		return ctm.RequiresNew(ctx, transactionFunc)
	default:
		return ErrUnSupportTransactionDefinition
	}
}

func (ctm contextTransactionManager) Required(ctx context.Context, transactionFunc TransactionFunc) error {
	if ctm.IsInTransaction(ctx) {
		value := ctx.Value(ContextSessionKey{})
		session := value.(*Session)
		// when in transaction,transaction is managed by the first called
		return transactionFunc(ctx, session)
	} else {
		session := ctm.engine.NewSession()
		sessionCtx := context.WithValue(ctx, ContextSessionKey{}, session)
		if err := session.Begin(); err != nil {
			return err
		}
		defer session.Close()
		err := transactionFunc(sessionCtx, session)
		if err != nil {
			return err
		}
		if err := session.Commit(); err != nil {
			return err
		}
		return err
	}

}

func (ctm contextTransactionManager) Supports(ctx context.Context, transactionFunc TransactionFunc) error {
	if ctm.IsInTransaction(ctx) {
		value := ctx.Value(ContextSessionKey{})
		session := value.(*Session)
		// when in transaction,transaction is managed by the first called
		return transactionFunc(ctx, session)
	} else {
		session := ctm.engine.NewSession()
		defer session.Close()
		return transactionFunc(ctx, session)
	}
}

func (ctm contextTransactionManager) Mandatory(ctx context.Context, transactionFunc TransactionFunc) error {
	if ctm.IsInTransaction(ctx) {
		value := ctx.Value(ContextSessionKey{})
		session := value.(*Session)
		// when in transaction,transaction is managed by the first called
		return transactionFunc(ctx, session)
	} else {
		return ErrMandatoryInTransaction
	}
}

func (ctm contextTransactionManager) RequiresNew(ctx context.Context, transactionFunc TransactionFunc) error {
	session := ctm.engine.NewSession()
	sessionCtx := context.WithValue(ctx, ContextSessionKey{}, session)
	if err := session.Begin(); err != nil {
		return err
	}
	defer session.Close()
	err := transactionFunc(sessionCtx, session)
	if err != nil {
		return err
	}
	if err := session.Commit(); err != nil {
		return err
	}
	return err
}

func (ctm contextTransactionManager) NotSupported(ctx context.Context, transactionFunc TransactionFunc) error {
	sessionCtx := context.WithValue(ctx, ContextSessionKey{}, nil)
	session := ctm.engine.NewSession()
	defer session.Close()
	return transactionFunc(sessionCtx, session)
}

func (ctm contextTransactionManager) Never(ctx context.Context, transactionFunc TransactionFunc) error {
	if ctm.IsInTransaction(ctx) {
		return ErrNeverNotSupported
	}
	session := ctm.engine.NewSession()
	defer session.Close()
	return transactionFunc(ctx, session)
}

func (ctm contextTransactionManager) Nested(ctx context.Context, transactionFunc TransactionFunc) error {
	panic("implement me")
}

func (ctm contextTransactionManager) IsInTransaction(ctx context.Context) bool {
	value := ctx.Value(ContextSessionKey{})
	_, ok := value.(*Session)
	return ok
}

func NewContextTransactionManager(engine EngineInterface) TransactionManager {
	return contextTransactionManager{engine: engine}
}
