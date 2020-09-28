package transaction_manager

import (
	"context"
	"errors"
	. "xorm.io/xorm"
)

var (
	ContextParmaError = errors.New("context transaction manager need context")
)

type ContextSessionKey struct{}

type contextTransactionManager struct {
	engine EngineInterface
}

func (ctm contextTransactionManager) DynamicSessionFunc() DynamicSession {
	return func(ctx ...context.Context) (*Session, error) {
		if len(ctx) == 0 {
			panic(ContextParmaError)
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

func (ctm contextTransactionManager) Do(ctx context.Context, do TransactionFunc) error {
	value := ctx.Value(ContextSessionKey{})
	session, ok := value.(*Session)
	// when context has session,transaction is managed by the first called
	if ok {
		return do(ctx, session)
	}
	session = ctm.engine.NewSession()
	sessionCtx := context.WithValue(ctx, ContextSessionKey{}, session)
	if err := session.Begin(); err != nil {
		return err
	}
	err := do(sessionCtx, session)
	if err != nil {
		return err
	}
	if err := session.Commit(); err != nil {
		return err
	}
	return session.Close()

}

func NewContextTransactionManager(engine EngineInterface) TransactionManager {
	return contextTransactionManager{engine: engine}
}
