package plan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var ErrNotFound = errors.New("plan not found")

const schema = `
CREATE TABLE IF NOT EXISTS warehouse_plan (
    warehouse_code text PRIMARY KEY,
    width          integer NOT NULL,
    height         integer NOT NULL
);
CREATE TABLE IF NOT EXISTS zone (
    id             bigserial PRIMARY KEY,
    warehouse_code text NOT NULL REFERENCES warehouse_plan (warehouse_code),
    label          text NOT NULL,
    x              integer NOT NULL,
    y              integer NOT NULL,
    width          integer NOT NULL,
    height         integer NOT NULL,
    UNIQUE (warehouse_code, label)
);
`

type warehousePlanRow struct {
	bun.BaseModel `bun:"table:warehouse_plan"`

	WarehouseCode string `bun:"warehouse_code,pk"`
	Width         int    `bun:"width"`
	Height        int    `bun:"height"`
}

type zoneRow struct {
	bun.BaseModel `bun:"table:zone"`

	ID            int64  `bun:"id,pk,autoincrement"`
	WarehouseCode string `bun:"warehouse_code"`
	Label         string `bun:"label"`
	X             int    `bun:"x"`
	Y             int    `bun:"y"`
	Width         int    `bun:"width"`
	Height        int    `bun:"height"`
}

type Repository struct {
	db *bun.DB
}

func Open(databaseURL string) *Repository {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(databaseURL)))
	return &Repository{db: bun.NewDB(sqldb, pgdialect.New())}
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *Repository) Migrate(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	return nil
}

func (r *Repository) Seed(ctx context.Context, plans []Plan) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		for _, p := range plans {
			header := warehousePlanRow{WarehouseCode: p.WarehouseCode, Width: p.Width, Height: p.Height}
			if _, err := tx.NewInsert().Model(&header).On("CONFLICT DO NOTHING").Exec(ctx); err != nil {
				return fmt.Errorf("insert plan %s: %w", p.WarehouseCode, err)
			}
			for _, z := range p.Zones {
				row := zoneRow{
					WarehouseCode: p.WarehouseCode, Label: z.Label, X: z.X, Y: z.Y, Width: z.Width, Height: z.Height,
				}
				if _, err := tx.NewInsert().Model(&row).On("CONFLICT DO NOTHING").Exec(ctx); err != nil {
					return fmt.Errorf("insert zone %s/%s: %w", p.WarehouseCode, z.Label, err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("seed plans: %w", err)
	}
	return nil
}

func (r *Repository) FindByCode(ctx context.Context, code string) (Plan, error) {
	var header warehousePlanRow
	err := r.db.NewSelect().Model(&header).Where("warehouse_code = ?", code).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return Plan{}, ErrNotFound
	}
	if err != nil {
		return Plan{}, fmt.Errorf("select plan %s: %w", code, err)
	}

	var rows []zoneRow
	err = r.db.NewSelect().Model(&rows).Where("warehouse_code = ?", code).Order("id").Scan(ctx)
	if err != nil {
		return Plan{}, fmt.Errorf("select zones of %s: %w", code, err)
	}

	zones := make([]Zone, len(rows))
	for i, row := range rows {
		zones[i] = Zone{Label: row.Label, X: row.X, Y: row.Y, Width: row.Width, Height: row.Height}
	}
	return Plan{WarehouseCode: header.WarehouseCode, Width: header.Width, Height: header.Height, Zones: zones}, nil
}
