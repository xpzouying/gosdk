package gosdk

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	apps []App
}

func NewServer(apps []App) *Server {

	return &Server{
		apps: apps,
	}
}

func (s *Server) Serve() error {

	for _, app := range s.apps {
		if err := app.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Shutdown() {
	parent := context.Background()

	g, ctx := errgroup.WithContext(parent)
	for _, app := range s.apps {
		app := app

		g.Go(func() error {
			return app.Stop(ctx)
		})
	}

	_ = g.Wait()
}
