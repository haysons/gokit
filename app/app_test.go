package app

import (
	"context"
	"testing"
	"time"

	"github.com/haysons/gokit/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

type TestService struct {
	Name string
}

func NewTestService() *TestService {
	return &TestService{Name: "test-service"}
}

var InitCalled bool

func InitTest(s *TestService) {
	if s != nil && s.Name == "test-service" {
		InitCalled = true
	}
}

func TestConfigOptions(t *testing.T) {
	cfg := Config{}

	WithConfig(Config{Name: "app"})(&cfg)
	assert.Equal(t, "app", cfg.Name)

	WithName("my-app")(&cfg)
	WithVersion("1.2.3")(&cfg)
	WithCommit("abcdef")(&cfg)

	assert.Equal(t, "my-app", cfg.Name)
	assert.Equal(t, "1.2.3", cfg.Version)
	assert.Equal(t, "abcdef", cfg.Commit)
}

func TestNewApp(t *testing.T) {
	log.SetDefaultSlog(&log.Config{
		Level:        "debug",
		ConsoleFmt:   true,
		ConsoleColor: true,
	})
	app := New(
		WithName("my-test-app"),
		WithVersion("0.1.0"),
		WithCommit("123456"),
		WithProvides(fx.Provide(NewTestService)),
		WithInvokes(fx.Invoke(InitTest)),
	)

	assert.NotNil(t, app)
	assert.Equal(t, "my-test-app", app.cfg.Name)
	assert.Equal(t, "0.1.0", app.cfg.Version)
	assert.Equal(t, "123456", app.cfg.Commit)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := app.Start(ctx)
	assert.NoError(t, err)

	err = app.Stop(ctx)
	assert.NoError(t, err)

	assert.True(t, InitCalled)
}

func TestWithModules(t *testing.T) {
	module := fx.Options(
		fx.Provide(NewTestService),
		fx.Invoke(InitTest),
	)

	app := New(
		WithName("module-app"),
		WithModules(module),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := app.Start(ctx)
	assert.NoError(t, err)
	err = app.Stop(ctx)
	assert.NoError(t, err)
	assert.True(t, InitCalled)
}
