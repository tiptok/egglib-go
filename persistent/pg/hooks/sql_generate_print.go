package hooks

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/linmadan/egglib-go/log"
)

type SqlGeneratePrintHook struct {
	Logger log.Logger
}

func (hook SqlGeneratePrintHook) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (hook SqlGeneratePrintHook) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	sqlStr, err := q.FormattedQuery()
	if err != nil {
		return err
	}
	if hook.Logger == nil {
		fmt.Println(string(sqlStr))
		return nil
	}
	hook.Logger.Debug(string(sqlStr))
	return nil
}
