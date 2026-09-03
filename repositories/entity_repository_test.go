package repositories

import (
	"cli-music-reviewer/interfaces"
	"cli-music-reviewer/models/dtos"
	"cli-music-reviewer/models/entities"
	"database/sql"
	"fmt"
	"reflect"
	"runtime/debug"
	"sync/atomic"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// testTimeout bounds every DB call made through withTimeout/withTimeoutErr.
// The known ScanRow nil-pointer bug (see runEntityRepositoryCRUD) doesn't just
// panic: it leaks the sqlite connection/cursor a panicking call had checked
// out, which (in shared-cache mode, or once a size-limited pool is exhausted)
// wedges every future query on that DB forever. Without a timeout, one such
// bug turns a test failure into the whole suite hanging until `go test
// -timeout` kills the binary.
const testTimeout = 2 * time.Second

func withTimeout[R any](t *testing.T, label string, fn func() (R, error)) (R, error) {
	t.Helper()
	type result struct {
		val R
		err error
	}
	ch := make(chan result, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				ch <- result{err: fmt.Errorf("panic: %v", r)}
			}
		}()
		v, err := fn()
		ch <- result{val: v, err: err}
	}()
	select {
	case res := <-ch:
		return res.val, res.err
	case <-time.After(testTimeout):
		t.Fatalf("%s: timed out after %s (likely a wedged sqlite connection left by an earlier panic)", label, testTimeout)
		var zero R
		return zero, nil
	}
}

func withTimeoutErr(t *testing.T, label string, fn func() error) error {
	t.Helper()
	v, err := withTimeout(t, label, func() (struct{}, error) { return struct{}{}, fn() })
	_ = v
	return err
}

var testDBCounter atomic.Int64

// newTestDB opens an in-memory DB with a shared cache so every pooled
// connection sees the same schema/data. A plain ":memory:" DSN gives each
// pooled connection its own private, empty database, which surfaces as
// flaky "no such table" errors as soon as database/sql opens a second
// connection. Each call gets a uniquely named DSN so tests don't share state.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:testdb%d?mode=memory&cache=shared", testDBCounter.Add(1))
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	schema := []string{
		`CREATE TABLE entry_rows (id INTEGER PRIMARY KEY, title TEXT, body TEXT, created_at TEXT, updated_at TEXT, active BOOLEAN NOT NULL DEFAULT TRUE)`,
		`CREATE TABLE spotify_tokens (id INTEGER PRIMARY KEY, access_token TEXT, refresh_token TEXT, expires_at date, updated_at date)`,
	}
	for _, stmt := range schema {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("create schema: %v", err)
		}
	}
	return db
}

// recoverAsFailure converts a panic into a normal test failure so one broken
// case doesn't take down the whole test binary and hide every other result.
// Only safe to use around calls that don't leak a DB connection on panic —
// see withTimeout for calls that do (anything reaching ScanRow/ScanRows).
func recoverAsFailure(t *testing.T) {
	t.Helper()
	if r := recover(); r != nil {
		t.Fatalf("panicked: %v\n%s", r, debug.Stack())
	}
}

// runEntityRepositoryCRUD exercises the full EntityRepository contract against
// a real (in-memory) DB. It is instantiated once per concrete entity type so
// the same assertions prove the generic implementation behaves consistently
// across types, not just for whichever one was hand-tested.
//
// Each subtest gets its own fresh DB via newRepo (rather than one shared
// instance) so a wedged connection from one subtest's panic can't hang every
// subtest that runs after it.
func runEntityRepositoryCRUD[T interfaces.EntityInterface](t *testing.T, newRepo func(t *testing.T) *EntityRepository[T], create func() T, touch func(T)) {
	t.Helper()

	t.Run("Create_then_FindByID", func(t *testing.T) {
		repo := newRepo(t)

		saved, err := withTimeout(t, "Create", func() (T, error) { return repo.Create(create()) })
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if saved.GetID() == 0 {
			t.Fatalf("Create did not assign an id")
		}

		found, err := withTimeout(t, "FindByID", func() (T, error) { return repo.FindByID(saved.GetID()) })
		if err != nil {
			t.Fatalf("FindByID: %v", err)
		}
		if found.GetID() != saved.GetID() {
			t.Fatalf("FindByID id = %d, want %d", found.GetID(), saved.GetID())
		}
	})

	t.Run("FindAll_and_Count_agree", func(t *testing.T) {
		repo := newRepo(t)

		before, err := withTimeout(t, "Count", func() (int, error) { return repo.Count() })
		if err != nil {
			t.Fatalf("Count: %v", err)
		}

		if _, err := withTimeout(t, "Create", func() (T, error) { return repo.Create(create()) }); err != nil {
			t.Fatalf("Create: %v", err)
		}

		after, err := withTimeout(t, "Count", func() (int, error) { return repo.Count() })
		if err != nil {
			t.Fatalf("Count: %v", err)
		}
		if after != before+1 {
			t.Fatalf("Count = %d, want %d", after, before+1)
		}

		all, err := withTimeout(t, "FindAll", func() ([]T, error) { return repo.FindAll() })
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		if len(all) != after {
			t.Fatalf("FindAll returned %d rows, want %d", len(all), after)
		}
	})

	t.Run("GetLatestOrNull_on_empty_table", func(t *testing.T) {
		repo := newRepo(t)

		latest, err := withTimeout(t, "GetLatestOrNull", func() (T, error) { return repo.GetLatestOrNull() })
		if err != nil {
			t.Fatalf("GetLatestOrNull: %v", err)
		}
		if !reflect.ValueOf(latest).IsNil() {
			t.Fatalf("GetLatestOrNull on empty table = %+v, want zero value", latest)
		}
	})

	t.Run("GetLatestOrNull_returns_most_recent", func(t *testing.T) {
		repo := newRepo(t)

		if _, err := withTimeout(t, "Create", func() (T, error) { return repo.Create(create()) }); err != nil {
			t.Fatalf("Create: %v", err)
		}
		saved, err := withTimeout(t, "Create", func() (T, error) { return repo.Create(create()) })
		if err != nil {
			t.Fatalf("Create: %v", err)
		}

		latest, err := withTimeout(t, "GetLatestOrNull", func() (T, error) { return repo.GetLatestOrNull() })
		if err != nil {
			t.Fatalf("GetLatestOrNull: %v", err)
		}
		if latest.GetID() != saved.GetID() {
			t.Fatalf("GetLatestOrNull id = %d, want %d (most recently created)", latest.GetID(), saved.GetID())
		}
	})

	t.Run("Update_persists_change", func(t *testing.T) {
		repo := newRepo(t)

		saved, err := withTimeout(t, "Create", func() (T, error) { return repo.Create(create()) })
		if err != nil {
			t.Fatalf("Create: %v", err)
		}

		touch(saved)
		if err := withTimeoutErr(t, "Update", func() error { return repo.Update(saved) }); err != nil {
			t.Fatalf("Update: %v", err)
		}

		found, err := withTimeout(t, "FindByID", func() (T, error) { return repo.FindByID(saved.GetID()) })
		if err != nil {
			t.Fatalf("FindByID after update: %v", err)
		}
		if !reflect.DeepEqual(found, saved) {
			t.Fatalf("after Update, FindByID = %+v, want %+v", found, saved)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		repo := newRepo(t)

		saved, err := withTimeout(t, "Create", func() (T, error) { return repo.Create(create()) })
		if err != nil {
			t.Fatalf("Create: %v", err)
		}

		ok, err := withTimeout(t, "Exists", func() (bool, error) { return repo.Exists(saved.GetID()) })
		if err != nil {
			t.Fatalf("Exists: %v", err)
		}
		if !ok {
			t.Fatalf("Exists(%d) = false, want true", saved.GetID())
		}

		ok, err = withTimeout(t, "Exists", func() (bool, error) { return repo.Exists(saved.GetID() + 1000) })
		if err != nil {
			t.Fatalf("Exists: %v", err)
		}
		if ok {
			t.Fatalf("Exists on unknown id = true, want false")
		}
	})

	t.Run("Delete_removes_row", func(t *testing.T) {
		repo := newRepo(t)

		saved, err := withTimeout(t, "Create", func() (T, error) { return repo.Create(create()) })
		if err != nil {
			t.Fatalf("Create: %v", err)
		}

		if err := withTimeoutErr(t, "Delete", func() error { return repo.Delete(saved.GetID()) }); err != nil {
			t.Fatalf("Delete: %v", err)
		}

		ok, err := withTimeout(t, "Exists", func() (bool, error) { return repo.Exists(saved.GetID()) })
		if err != nil {
			t.Fatalf("Exists after delete: %v", err)
		}
		if ok {
			t.Fatalf("row still exists after Delete")
		}
	})
}

func TestEntryRowRepository_CRUD(t *testing.T) {
	runEntityRepositoryCRUD(t,
		func(t *testing.T) *EntityRepository[*entities.EntryRow] {
			return NewEntityRepository[*entities.EntryRow](newTestDB(t))
		},
		func() *entities.EntryRow {
			return &entities.EntryRow{Title: "title", Body: "body", CreatedAt: "2026-01-01", UpdatedAt: "2026-01-01", Active: true}
		},
		func(e *entities.EntryRow) { e.Title = "updated title" },
	)
}

func TestSpotifyTokenRepository_CRUD(t *testing.T) {
	runEntityRepositoryCRUD(t,
		func(t *testing.T) *EntityRepository[*entities.SpotifyToken] {
			return NewEntityRepository[*entities.SpotifyToken](newTestDB(t))
		},
		func() *entities.SpotifyToken {
			return &entities.SpotifyToken{
				AccessToken:  "access",
				RefreshToken: "refresh",
				ExpiresAt:    time.Now().Add(time.Hour).UTC().Truncate(time.Second),
				UpdatedAt:    time.Now().UTC().Truncate(time.Second),
			}
		},
		func(s *entities.SpotifyToken) { s.AccessToken = "rotated-access" },
	)
}

// Regression test for a bug found while writing this suite: Create discards
// the error from db.Exec (it gets shadowed by the err from LastInsertId), so
// a failed insert panics on a nil sql.Result instead of returning an error.
// This one doesn't leak a connection (Exec fully releases its connection
// before returning, unlike QueryRow), so a plain recover is enough.
func TestEntityRepository_Create_ReturnsErrorOnConstraintViolation(t *testing.T) {
	defer recoverAsFailure(t)

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE entry_rows (id INTEGER PRIMARY KEY, title TEXT UNIQUE, body TEXT, created_at TEXT, updated_at TEXT, active BOOLEAN)`); err != nil {
		t.Fatalf("create schema: %v", err)
	}

	repo := NewEntityRepository[*entities.EntryRow](db)

	if _, err := repo.Create(&entities.EntryRow{Title: "dup"}); err != nil {
		t.Fatalf("first Create: %v", err)
	}

	if _, err := repo.Create(&entities.EntryRow{Title: "dup"}); err == nil {
		t.Fatalf("second Create with duplicate title: want a UNIQUE constraint error, got nil error")
	}
}

// Regression test for a second, worse effect of the ScanRow nil-pointer bug:
// database/sql.DB.QueryRow checks out a connection immediately, and it is
// normally Row.Scan's job to release it (via a deferred rows.Close()).
// ScanRow panics while still building its argument list, before Scan is ever
// entered, so that connection is never released. On a size-limited pool,
// enough leaked connections exhaust it and every later query hangs forever
// instead of returning an error.
func TestEntityRepository_FindByID_PanicLeaksConnection(t *testing.T) {
	db, err := sql.Open("sqlite", "file:leaktest?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(`CREATE TABLE entry_rows (id INTEGER PRIMARY KEY, title TEXT, body TEXT, created_at TEXT, updated_at TEXT, active BOOLEAN)`); err != nil {
		t.Fatalf("create schema: %v", err)
	}

	repo := NewEntityRepository[*entities.EntryRow](db)
	saved, err := repo.Create(&entities.EntryRow{Title: "t"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	func() {
		defer func() { recover() }() // the known ScanRow nil-receiver panic
		repo.FindByID(saved.GetID())
	}()

	if _, err := withTimeout(t, "Count", func() (int, error) { return repo.Count() }); err != nil {
		t.Fatalf("Count after a panicking FindByID: %v", err)
	}
}

func TestEntryRowRepository_GetActiveRows(t *testing.T) {
	db := newTestDB(t)
	repo := NewEntryRowRepository(db)

	if _, err := repo.Create(&entities.EntryRow{Title: "active one", Active: true}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := repo.Create(&entities.EntryRow{Title: "inactive one", Active: false}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	active, err := withTimeout(t, "GetActiveRows", func() ([]*entities.EntryRow, error) { return repo.GetActiveRows() })
	if err != nil {
		t.Fatalf("GetActiveRows: %v", err)
	}
	if len(active) != 1 {
		t.Fatalf("GetActiveRows returned %d rows, want 1", len(active))
	}
	if active[0].Title != "active one" {
		t.Fatalf("GetActiveRows returned %q, want %q", active[0].Title, "active one")
	}
}

func TestSpotifyTokenRepository_CreateFromDTO(t *testing.T) {
	defer recoverAsFailure(t)

	db := newTestDB(t)
	repo := NewSpotifyTokenRepository(db)

	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	saved, err := repo.CreateFromDTO(&dtos.CreateSpotifyTokenDTO{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		t.Fatalf("CreateFromDTO: %v", err)
	}

	if saved.AccessToken != "access" || saved.RefreshToken != "refresh" {
		t.Fatalf("CreateFromDTO did not map DTO fields: %+v", saved)
	}
	if !saved.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("ExpiresAt = %v, want %v", saved.ExpiresAt, expiresAt)
	}
	if saved.UpdatedAt.IsZero() {
		t.Fatalf("CreateFromDTO did not stamp UpdatedAt")
	}
}
