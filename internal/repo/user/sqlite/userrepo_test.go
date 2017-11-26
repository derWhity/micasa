package sqlite_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/derWhity/micasa/internal/fsutils"
	"github.com/derWhity/micasa/internal/log"
	"github.com/derWhity/micasa/internal/migrate"
	"github.com/derWhity/micasa/internal/models"
	"github.com/derWhity/micasa/internal/repo"
	"github.com/derWhity/micasa/internal/repo/user/sqlite"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Just needed for the sqlite driver
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testDbName = filepath.Join(os.TempDir(), uuid.NewV4().String(), "test.db")
	testUsers  = []models.User{
		{FullName: "Jack Harkness", Name: "boe"},
		{FullName: "Time And Relative Dimensions In Space", Name: "tardis"},
		{FullName: "K9", Name: "k9"},
		{FullName: "Donna Noble", Name: "donna"},
		{FullName: "John Smith", Name: "doctor"},
		{FullName: "Rory Williams", Name: "rory"},
		{FullName: "Amy Pond", Name: "amy"},
		{FullName: "Rose Tyler", Name: "badwolf"},
		{FullName: "Clara Oswald", Name: "clara"},
		{FullName: "Ashildr", Name: "knightmare"},
	}
	testPassword = "secret"
)

type testLogWriter struct {
	t      *testing.T
	active bool // Activate this to log the output
}

// Write is the implementation for io.Writer we need to use the test output for writing log data
func (t *testLogWriter) Write(p []byte) (n int, err error) {
	if t.active {
		t.t.Logf("%s", p)
	}
	return len(p), nil
}

func createTestUsers(r *sqlite.UserRepo) error {
	for index, user := range testUsers {
		user.SetPassword(testPassword)
		if err := r.Create(&user); err != nil {
			return err
		}
		testUsers[index].ID = user.ID // Make sure that we know the new generated ID
	}
	return nil
}

func createTestLogger(t *testing.T) log.Logger {
	var logger log.Logger
	writer := testLogWriter{t: t, active: false}
	logger = log.New(kitlog.NewLogfmtLogger(&writer), log.LvlDebug)
	logger = logger.With(
		log.FldTimestamp, kitlog.DefaultTimestampUTC,
	)
	return logger
}

func setupTestDB(logger log.Logger) (*sqlx.DB, error) {
	teardownTestDB(nil, logger)
	fsutils.CheckAndCreateDir(filepath.Dir(testDbName), logger)
	db, err := sqlx.Open("sqlite3", testDbName)
	if err != nil {
		return nil, err
	}
	// Perform DB migrations
	if err = migrate.ExecuteMigrationsOnDb(db, logger); err != nil {
		logger.Crit("Database migration has failed", log.FldError, err)
		return nil, err
	}
	return db, nil
}

func teardownTestDB(db *sqlx.DB, logger log.Logger) error {
	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
	}
	dir := filepath.Dir(testDbName)
	logger.Info(fmt.Sprintf("Deleting test DB directory '%s'", dir))
	if err := os.RemoveAll(dir); err != nil {
		return errors.Wrap(err, "Failed to delete temporary test DB dir")
	}
	return nil
}

func getAllUsers(db *sqlx.DB) []*models.User {
	query := "SELECT * FROM Users"
	users := []*models.User{}
	err := db.Select(&users, query)
	So(err, ShouldBeNil)
	return users
}

func TestCreate(t *testing.T) {
	Convey("Having a test database instance", t, func() {
		logger := createTestLogger(t)
		db, err := setupTestDB(logger)
		So(err, ShouldBeNil)

		Convey("Having a UserRepo instance", func() {
			r := sqlite.New(db)
			So(r, ShouldNotBeNil)
			Convey("Creating a new user should work without error", func() {
				const id = "youshouldnotseeme"
				const pw = "testitest"
				user := models.User{
					FullName: "Teddy Tester",
					Name:     "teDdy",
					ID:       id,
				}
				user.SetPassword(pw)
				err := r.Create(&user)
				So(err, ShouldBeNil)
				So(user.ID, ShouldNotEqual, id)
				// Check if the user was created in the DB correctly
				users := getAllUsers(db)
				So(users, ShouldHaveLength, 1)
				compareUsers(*users[0], user, pw)
			})

			Convey("Creating a user with an already-existing username should fail", func() {
				user := models.User{
					Name: "teddy",
				}
				err := r.Create(&user)
				So(err, ShouldBeNil)
				// And now the second one
				user.FullName = "Blah"
				err = r.Create(&user)
				So(err, ShouldEqual, repo.ErrDuplicate) // Specific duplication error
				users := getAllUsers(db)
				So(users, ShouldHaveLength, 1)
				So(users[0].FullName, ShouldNotEqual, "Blah")

			})
		})

		Reset(func() {
			if err := teardownTestDB(db, logger); err != nil {
				panic(err)
			}
		})
	})
}

func compareUsers(a models.User, b models.User, passwd string) {
	So(a.ID, ShouldEqual, b.ID)
	So(a.FullName, ShouldEqual, b.FullName)
	So(strings.ToLower(a.Name), ShouldEqual, b.Name)
	So(b.CheckPassword(passwd), ShouldBeNil)
}

func getUser(db *sqlx.DB, id string) models.User {
	query := `SELECT * FROM Users WHERE userid = ?`
	user := models.User{}
	So(db.Get(&user, query, id), ShouldBeNil)
	return user
}

func TestUpdate(t *testing.T) {
	Convey("Having a test database instance", t, func() {
		logger := createTestLogger(t)
		db, err := setupTestDB(logger)
		So(err, ShouldBeNil)

		Convey("Having a UserRepo instance", func() {
			r := sqlite.New(db)
			So(r, ShouldNotBeNil)

			Convey("Having a set of users", func() {
				So(createTestUsers(r), ShouldBeNil)

				Convey("Updating users should work", func() {
					for i := 0; i < len(testUsers); i++ {
						user := testUsers[i]
						txt := fmt.Sprintf("Update_%d", i)
						user.Name = txt
						user.FullName = txt
						user.SetPassword(txt)
						So(r.Update(&user), ShouldBeNil)
						// Now load and check the user
						updatedUser := getUser(db, string(user.ID))
						compareUsers(user, updatedUser, txt)
					}
				})

				Convey("Updating a nonexistent user should fail", func() {
					u := testUsers[0]
					u.ID = "IDoNotExist"
					u.FullName = "xxx"
					u.Name = "yyy"
					u.SetPassword("Hurz")
					err := r.Update(&u)
					So(err, ShouldEqual, repo.ErrNotExisting)
					// Check all the test users
					for _, u = range testUsers {
						user := getUser(db, string(u.ID))
						compareUsers(u, user, testPassword)
					}
				})
			})
		})

		Reset(func() {
			if err := teardownTestDB(db, logger); err != nil {
				panic(err)
			}
		})
	})
}

func TestDelete(t *testing.T) {
	Convey("Having a test database instance", t, func() {
		logger := createTestLogger(t)
		db, err := setupTestDB(logger)
		So(err, ShouldBeNil)

		Convey("Having a UserRepo instance", func() {
			r := sqlite.New(db)
			So(r, ShouldNotBeNil)

			Convey("Having users in the database", func() {
				So(createTestUsers(r), ShouldBeNil)

				Convey("Deleting existing users should work", func() {
					for i := 0; i < len(testUsers); i++ {
						So(r.Delete(testUsers[i].ID), ShouldBeNil)
						var row int64
						So(db.Get(&row, "SELECT COUNT(*) AS count FROM Users"), ShouldBeNil)
						So(row, ShouldEqual, int64(len(testUsers)-(i+1)))
					}
				})

				Convey("Deleting non-existing users should return successfully but have no impact on the users", func() {
					So(r.Delete(models.UserID("IDoNotExist")), ShouldBeNil)
					var row int64
					So(db.Get(&row, "SELECT COUNT(*) AS count FROM Users"), ShouldBeNil)
					So(row, ShouldEqual, int64(len(testUsers)))
				})
			})
		})

		Reset(func() {
			if err := teardownTestDB(db, logger); err != nil {
				panic(err)
			}
		})

	})
}

func TestGetByID(t *testing.T) {
	Convey("Having a test database instance", t, func() {
		logger := createTestLogger(t)
		db, err := setupTestDB(logger)
		So(err, ShouldBeNil)

		Convey("Having a UserRepo instance", func() {
			r := sqlite.New(db)
			So(r, ShouldNotBeNil)

			Convey("Having users in the database", func() {
				So(createTestUsers(r), ShouldBeNil)

				Convey("Searching for existing users should return exactly those users", func() {
					for _, user := range testUsers {
						result, err := r.GetByID(user.ID)
						So(err, ShouldBeNil)
						compareUsers(user, *result, testPassword)
					}
				})

				Convey("Searching for non-existing users should yield an ErrNotExisting error", func() {
					for _, id := range []string{"UnknownID", "Nonexistin", "AlsoNotExisting"} {
						result, err := r.GetByID(models.UserID(id))
						So(err, ShouldEqual, repo.ErrNotExisting)
						So(result, ShouldBeNil)
					}
				})
			})
		})

		Reset(func() {
			if err := teardownTestDB(db, logger); err != nil {
				panic(err)
			}
		})

	})
}

func TestGetByCredentials(t *testing.T) {
	t.Skip("Not implemented")
}

func TestFind(t *testing.T) {
	t.Skip("Not implemented")
}

func TestExists(t *testing.T) {
	Convey("Having a test database instance", t, func() {
		logger := createTestLogger(t)
		db, err := setupTestDB(logger)
		So(err, ShouldBeNil)

		Convey("Having a UserRepo instance", func() {
			r := sqlite.New(db)
			So(r, ShouldNotBeNil)

			Convey("Having users in the database", func() {
				So(createTestUsers(r), ShouldBeNil)

				Convey("Checking the existence of created users should yield `true`", func() {
					for _, user := range testUsers {
						result, err := r.Exists(user.ID)
						So(err, ShouldBeNil)
						So(result, ShouldBeTrue)
					}
				})

				Convey("Checking the existence of non-created users should yield `false`", func() {
					for _, id := range []models.UserID{"YouDontKnowMe", "MeNeither", "NotMe", "Bob"} {
						result, err := r.Exists(id)
						So(err, ShouldBeNil)
						So(result, ShouldBeFalse)
					}
				})
			})
		})

		Reset(func() {
			if err := teardownTestDB(db, logger); err != nil {
				panic(err)
			}
		})

	})
}
