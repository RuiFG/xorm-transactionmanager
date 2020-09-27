# Xorm-TransactionManager
######
xorm transaction-manager is the xorm automatic transaction manager. With this library, 
you don’t need to worry about global transaction.
### Routine Transaction Manager
routine manager save session in the `sync.map`,in same go routine,DynamicSession function returns the same session
### Context Transaction Manager
context manager save session in the `context.Context`,DynamicSession function will get the session in the context

# Installation
```shell script
go get github.com/RuiFG/xorm-transactionmanager 
```
# Simple Example
## Routine Manager
```go
package main

import (
"context"
"github.com/RuiFG/xorm-transactionmanager"
"github.com/go-xorm/xorm"
)

func main() {
	DB, _ := xorm.NewEngine("mysql", "bp_ops:bp_ops@baishan.com@tcp(172.18.2.21:3306)/bp_ops?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
    // use routine transaction manager
	manager := transaction_manager.NewRoutineTransactionManager(DB)
	err := manager.Do(context.Background(), func(ctx context.Context, session *xorm.Session) error {
		return manager.Do(ctx, func(ctx context.Context, session *xorm.Session) error {
			m := make([]int, 0)
			session, _ = manager.DynamicSessionFunc()(ctx)
			_ = session.Table("test").Cols("id").Find(&m)
			fmt.Println(m)
			return nil
		})
	})
	fmt.Print(err)
}
```
## Context Manager
```go
package main

import (
"context"
"fmt"
"github.com/RuiFG/xorm-transactionmanager"
"github.com/go-xorm/xorm"
)

func main() {
	DB, _ := xorm.NewEngine("mysql", "bp_ops:bp_ops@baishan.com@tcp(172.18.2.21:3306)/bp_ops?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
    // use context transaction manager
	manager := transaction_manager.NewContextTransactionManager(DB)
	err := manager.Do(context.Background(), func(ctx context.Context, session *xorm.Session) error {
		return manager.Do(ctx, func(ctx context.Context, session *xorm.Session) error {
			m := make([]int, 0)
			session, _ = manager.DynamicSessionFunc()(ctx)
			_ = session.Table("test").Cols("id").Find(&m)
			fmt.Println(m)
			return nil
		})
	})
	fmt.Print(err)
}
```

# License
This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.