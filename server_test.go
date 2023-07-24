package gosdk

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

type recordApp struct {
	count atomic.Int32
}

func NewRecordApp() *recordApp {
	return &recordApp{
		count: atomic.Int32{},
	}
}

func (r *recordApp) Start() error {
	r.count.Add(1)
	return nil
}

func (r *recordApp) Stop(_ context.Context) error {
	r.count.Add(-1)
	return nil
}

func TestServerServeForOneApp(t *testing.T) {

	app := NewRecordApp()
	server := NewServer([]App{
		app,
	})

	assert.NoError(t, server.Serve())
	assert.Equal(t, int32(1), app.count.Load())

	server.Shutdown()
	assert.Equal(t, int32(0), app.count.Load())
}

func TestServerServeForMultipleApp(t *testing.T) {

	for i := 0; i < 100; i++ {
		apps := newMultipleApps(i)
		server := NewServer(apps)

		assert.NoError(t, server.Serve())
		checkAllCountAppsStart(t, apps)

		server.Shutdown()
		checkAllCountAppsStop(t, apps)
	}
}

func newMultipleApps(n int) []App {

	apps := make([]App, 0, n)

	for i := 0; i < n; i++ {
		apps = append(apps, NewRecordApp())
	}
	return apps
}

func checkAllCountAppsStart(t *testing.T, apps []App) {

	for _, app := range apps {

		a := app.(*recordApp)
		want := int32(1)
		got := a.count.Load()
		assert.Equal(t, want, got)
	}
}

func checkAllCountAppsStop(t *testing.T, apps []App) {

	for _, app := range apps {

		a := app.(*recordApp)
		want := int32(0)
		got := a.count.Load()
		assert.Equal(t, want, got)
	}
}
