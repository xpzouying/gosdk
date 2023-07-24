package gosdk

import "context"

type App interface {
	Start() error
	Stop(ctx context.Context) error
}
