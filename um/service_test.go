package um

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gosom/kouti/dbdriver"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/testutils"
	"github.com/rs/zerolog"
)

var (
	testconn   dbdriver.DB
	testLogger zerolog.Logger
	srv        *Service
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	testLogger = logger.New(logger.Config{Debug: true})

	postgresContainer, p := testutils.SpinPostgresContainer(ctx)
	defer postgresContainer.Terminate(ctx)
	var err error
	testconn, err = dbdriver.New(ctx, dbdriver.Config{
		Logger:     testLogger,
		ConnString: testutils.GetTestDSN(p),
	})
	if err != nil {
		panic(err)
	}
	defer testconn.Close()

	_, err = testconn.RawConn().Exec(ctx, "CREATE EXTENSION pgcrypto")
	if err != nil {
		panic(err)
	}

	srv = NewService(Config{
		Log:         testLogger,
		DB:          testconn,
		SystemRoles: []string{"admin", "member"},
	})

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestInitSchema(t *testing.T) {
	defer resetDb()
	if err := srv.InitSchema(context.Background()); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestRegisterUser(t *testing.T) {
	if err := srv.InitSchema(context.Background()); err != nil {
		panic(err)
	}
	defer resetDb()
	rp := RegisterUserOpts{
		Identity: "giorgos@example.com",
		Passwd:   "password",
		Roles:    []string{"admin"},
	}
	u, err1, err2 := srv.RegisterUser(context.Background(), rp)
	if err1 != nil {
		t.Error(err1)
		return
	}
	if err2 != nil {
		t.Error(err2)
		return
	}
	if u.ID != 1 {
		t.Errorf("expected user.ID to be 1 but got %d", u.ID)
		return
	}
	if u.Identity != rp.Identity {
		t.Errorf("expected user.Identity = %s but got %s", rp.Identity, u.Identity)
		return
	}
	if len(u.Roles) != len(rp.Roles) {
		t.Errorf("wrong roles num")
		return
	}
	if u.Roles[0].Name != "admin" {
		t.Errorf("user expected to have role %s but got %s", "admin", u.Roles[0].Name)
		return
	}
}

func TestGetUserByUID(t *testing.T) {
	if err := srv.InitSchema(context.Background()); err != nil {
		panic(err)
	}
	defer resetDb()
	rp := RegisterUserOpts{
		Identity: "giorgos@example.com",
		Passwd:   "password",
		Roles:    []string{"admin"},
	}
	u, err1, err2 := srv.RegisterUser(context.Background(), rp)
	if err1 != nil {
		t.Error(err1)
		return
	}
	if err2 != nil {
		t.Error(err2)
		return
	}

	u2, err := srv.GetUserByUID(context.Background(), u.UID.String())
	if err != nil {
		t.Error(err)
		return
	}
	if u2.ID != u.ID {
		t.Errorf("expected user.ID to be %d but got %d", u.ID, u2.ID)
		return
	}
}

func TestSelectUsers(t *testing.T) {
	if err := srv.InitSchema(context.Background()); err != nil {
		panic(err)
	}
	defer resetDb()
	var usrs []User
	for i := 0; i < 100; i++ {
		rp := RegisterUserOpts{
			Identity: fmt.Sprintf("giorgos+%d@example.com", i),
			Passwd:   "password",
			Roles:    []string{"admin"},
		}
		u, err1, err2 := srv.RegisterUser(context.Background(), rp)
		if err1 != nil {
			t.Error(err1)
			return
		}
		if err2 != nil {
			t.Error(err2)
			return
		}
		usrs = append(usrs, u)
	}

	params := SelectUserParams{
		AfterID:  0,
		UseLimit: false,
	}
	retrieved, err := srv.SelectUsers(context.Background(), params)
	if err != nil {
		t.Error(err)
		return
	}
	if len(retrieved) != len(usrs) {
		t.Errorf("expected to retrieve %d users but got %d", len(usrs), len(retrieved))
		return
	}
	for i := 0; i < len(usrs); i++ {
		expected := usrs[i].Identity
		got := retrieved[i].Identity
		if expected != got {
			t.Errorf("expected identity %s but got %s", expected, got)
			return
		}
	}
}

func TestDeleteUser(t *testing.T) {
	if err := srv.InitSchema(context.Background()); err != nil {
		panic(err)
	}
	defer resetDb()
	rp := RegisterUserOpts{
		Identity: "giorgos@example.com",
		Passwd:   "password",
		Roles:    []string{"admin"},
	}
	u, err1, err2 := srv.RegisterUser(context.Background(), rp)
	if err1 != nil {
		t.Error(err1)
		return
	}
	if err2 != nil {
		t.Error(err2)
		return
	}
	if u.ID != 1 {
		t.Errorf("expected user.ID to be 1 but got %d", u.ID)
		return
	}
	if u.Identity != rp.Identity {
		t.Errorf("expected user.Identity = %s but got %s", rp.Identity, u.Identity)
		return
	}
	if len(u.Roles) != len(rp.Roles) {
		t.Errorf("wrong roles num")
		return
	}
	if u.Roles[0].Name != "admin" {
		t.Errorf("user expected to have role %s but got %s", "admin", u.Roles[0].Name)
		return
	}

	if err := srv.DeleteUser(context.Background(), GetUserParams{
		UseID: true,
		ID:    u.ID,
	}); err != nil {
		t.Error(err)
		return
	}
	var exists bool
	q := "select exists(select 1 from users_roles where user_id = $1)"
	if err := testconn.RawConn().QueryRow(context.Background(), q, u.ID).Scan(&exists); err != nil {
		t.Error(err)
		return
	}
	if exists {
		t.Errorf("expected user_roles not exist but they do")
		return
	}
}

// ==========================================================================

func resetDb() {
	tables := []string{"users_roles", "users", "roles"}
	q := "drop table %s"
	for _, t := range tables {
		_, _ = testconn.RawConn().Exec(context.Background(), fmt.Sprintf(q, t))
	}
}
