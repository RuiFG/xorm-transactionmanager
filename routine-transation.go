package transaction_manager

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	. "xorm.io/xorm"
)

var (
	MapError = errors.New("error loaded session in the map")
)
var managerCont uint64 = 0

type routineTransactionManager struct {
	id                    string
	synchronizeSessionMap sync.Map
	engine                EngineInterface
}

func (tm *routineTransactionManager) DynamicSessionFunc() DynamicSession {
	return func(...context.Context) (*Session, error) {
		store, _ := tm.synchronizeSessionMap.LoadOrStore(tm.id+strconv.FormatUint(curGoroutineID(), 10), tm.engine.NewSession())
		session, ok := store.(*Session)
		if !ok {
			return nil, MapError
		}
		return session, nil

	}
}

func (tm *routineTransactionManager) Do(ctx context.Context, do TransactionFunc) error {
	var result error
	store, loaded := tm.synchronizeSessionMap.LoadOrStore(tm.id+strconv.FormatUint(curGoroutineID(), 10), tm.engine.NewSession())
	session, ok := store.(*Session)
	if !ok {
		return MapError
	}
	//transaction is managed by the first called
	if !loaded {
		if err := session.Begin(); err != nil {
			return err
		}
		defer func() {
			// commit session  if result is nil
			if result == nil {
				result = session.Commit()
			}
			if err := session.Close(); err != nil {
				result = err
			}
			tm.synchronizeSessionMap.Delete(curGoroutineID())
		}()
	}
	result = do(ctx, session)
	return result
}

// NewRoutineTransactionManager
func NewRoutineTransactionManager(engine EngineInterface) TransactionManager {

	return &routineTransactionManager{id: strconv.FormatUint(atomic.AddUint64(&managerCont, 1), 10), engine: engine}
}
