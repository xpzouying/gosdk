package gosdk

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

type recordApp struct {
	runningState atomic.Bool

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup
}

// NewRecordApp to create app with slow stop time.
// If stopTime is 0, stop immediately.
func NewRecordApp() *recordApp {
	ctx, cancel := context.WithCancel(context.Background())

	return &recordApp{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (r *recordApp) Start() error {
	log.Printf("recordApp.Start() ...")

	wg := new(sync.WaitGroup) // 用于保护异步的 goroutine 运行成功后才能结束 Start() 函数。

	r.wg.Add(1)

	wg.Add(1) // 运行异步的 goroutine
	go func() {
		r.runningState.Store(true)
		log.Printf("the app is running")
		wg.Done() // 异步的 goroutine 配置运行状态成功

		<-r.ctx.Done()
		log.Printf("the app is stop running, stop the async goroutine....")
		r.wg.Done() // 相当于 「通知 App Stop」 可以结束了。
	}()

	wg.Wait()
	return nil
}

func (r *recordApp) Stop(_ context.Context) error {
	r.runningState.Store(false)

	if r.cancel != nil {
		r.cancel() // 为了取消运行中的 goroutine
	}
	log.Printf("stop the app")

	r.wg.Wait() // 为了等待运行中的 goroutine 结束

	return nil
}

func TestServerServeForOneApp(t *testing.T) {

	app := NewRecordApp() // no stop time
	server := NewServer([]App{
		app,
	})

	assert.NoError(t, server.Serve())
	assert.Equal(t, true, app.runningState.Load())

	server.Shutdown()
	assert.Equal(t, false, app.runningState.Load())
}

func TestServerServeForMultipleApp(t *testing.T) {

	for i := 0; i < 10; i++ {
		apps := newMultipleApps(i)
		server := NewServer(apps)

		assert.NoError(t, server.Serve())
		checkAllAppsState(t, true, apps)

		server.Shutdown()
		checkAllAppsState(t, false, apps)
	}
}

func newMultipleApps(n int) []App {

	apps := make([]App, 0, n)

	for i := 0; i < n; i++ {
		apps = append(apps, NewRecordApp())
	}
	return apps
}

func checkAllAppsState(t *testing.T, isRunningState bool, apps []App) {
	for _, app := range apps {
		a := app.(*recordApp)

		assert.Equal(t, isRunningState, a.runningState.Load())
	}
}
