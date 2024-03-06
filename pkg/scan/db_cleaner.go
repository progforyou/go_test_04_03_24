package scan

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"testing/dating/api/pkg/domain"
	"time"
)

var interval = time.Minute * 2

type ScanHandler struct {
	cancelFn context.CancelFunc
	db       *gorm.DB
	log      zerolog.Logger
	ctx      context.Context
}

func NewScanHandler(ctx context.Context, db *gorm.DB, logger zerolog.Logger) (*ScanHandler, error) {
	res := &ScanHandler{
		db:  db,
		log: logger.With().Str("module", "db scanner").Logger(),
		ctx: ctx,
	}

	ctx, res.cancelFn = context.WithCancel(ctx)
	res.log.Debug().Msg("Scan handler started")
	return res, nil
}

func (c *ScanHandler) Start() {
	go c.ScanDeletedTasks()
}

func (c *ScanHandler) Stop() {
	c.log.Info().Msg("Scan handler stopped")
	c.cancelFn()
}

func (c *ScanHandler) ScanDeletedTasks() {
	for {
		select {
		case <-time.After(interval):
			err := c.DeleteUserActivate()
			if err != nil {
				c.log.Error().Err(err).Msg("task delete handler error")
			}
			break
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *ScanHandler) DeleteUserActivate() error {
	c.log.Debug().Msg("scan for useless user activate codes...")
	tx := c.db.Where("created_at < ?", time.Now().Add(-interval)).Delete(&domain.UserActivate{})
	if tx.RowsAffected != 0 {
		c.log.Debug().Msg(fmt.Sprintf("found %d user activate codes, deleting..", tx.RowsAffected))
	}
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
