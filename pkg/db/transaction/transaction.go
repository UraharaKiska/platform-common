package transaction

import (
	"context"

	"github.com/UraharaKiska/platform-common/pkg/db"
	"github.com/UraharaKiska/platform-common/pkg/db/pg"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Transactor
}

func NewTransactorManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}
	// Стартуем новую транзакцию
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	// Кладем транзакцию в контекст
	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		// востанавливаемся после паники
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: #{r}")
		}

		// откатываем транзакцию, если произошла ошибка
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrap(err, "errRollback: #{errRollBack}")
			}
			return
		}

		// если ошибки не было, коммитим транзакцию
		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()
	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}
	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}
