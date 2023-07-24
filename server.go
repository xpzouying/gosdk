package gosdk

import "context"

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
	ctx := context.Background()

	for _, app := range s.apps {
		app.Stop(ctx)
	}
}
