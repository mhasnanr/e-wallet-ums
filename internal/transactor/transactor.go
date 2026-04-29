package transactor

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type transactor struct {
	db *gorm.DB
}

type TxKey struct{}

func NewTransactor(db *gorm.DB) *transactor {
	return &transactor{db: db}
}

func (t *transactor) WithinTransaction(ctx context.Context, txFunc func(context.Context) error) error {
	tx := t.db.WithContext(ctx).Begin(nil)
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	txCtx := context.WithValue(ctx, TxKey{}, tx)
	defer func() {
		_ = tx.Rollback().Error
	}()
	if err := txFunc(txCtx); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
