package plan

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func newRepository(t *testing.T) *Repository {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:16.4",
		postgres.WithDatabase("plan"),
		postgres.WithUsername("plan"),
		postgres.WithPassword("plan"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Errorf("terminate postgres container: %v", err)
		}
	})

	url, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}
	repo := Open(url)
	t.Cleanup(func() { _ = repo.Close() })
	return repo
}

func TestRepositoryReturnsSeededPlanWithZonesInInsertionOrder(t *testing.T) {
	ctx := context.Background()
	repo := newRepository(t)
	if err := repo.Migrate(ctx); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	if err := repo.Seed(ctx, DemoPlans); err != nil {
		t.Fatalf("Seed: %v", err)
	}

	got, err := repo.FindByCode(ctx, "LYS-01")
	if err != nil {
		t.Fatalf("FindByCode: %v", err)
	}

	want := DemoPlans[0]
	if got.WarehouseCode != want.WarehouseCode || got.Width != want.Width || got.Height != want.Height {
		t.Errorf("plan = %+v, want %+v", got, want)
	}
	if len(got.Zones) != len(want.Zones) {
		t.Fatalf("zones = %d, want %d", len(got.Zones), len(want.Zones))
	}
	for i := range want.Zones {
		if got.Zones[i] != want.Zones[i] {
			t.Errorf("zone %d = %+v, want %+v", i, got.Zones[i], want.Zones[i])
		}
	}
}

func TestRepositoryReturnsErrNotFoundForUnknownWarehouse(t *testing.T) {
	ctx := context.Background()
	repo := newRepository(t)
	if err := repo.Migrate(ctx); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	_, err := repo.FindByCode(ctx, "XXX-99")

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("FindByCode error = %v, want ErrNotFound", err)
	}
}

func TestMigrateAndSeedCanRunTwiceWithoutDuplicatingData(t *testing.T) {
	ctx := context.Background()
	repo := newRepository(t)
	for i := 0; i < 2; i++ {
		if err := repo.Migrate(ctx); err != nil {
			t.Fatalf("Migrate run %d: %v", i+1, err)
		}
		if err := repo.Seed(ctx, DemoPlans); err != nil {
			t.Fatalf("Seed run %d: %v", i+1, err)
		}
	}

	got, err := repo.FindByCode(ctx, "NTE-02")
	if err != nil {
		t.Fatalf("FindByCode: %v", err)
	}

	if len(got.Zones) != len(DemoPlans[1].Zones) {
		t.Errorf("zones after second seed = %d, want %d", len(got.Zones), len(DemoPlans[1].Zones))
	}
}
