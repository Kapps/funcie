package receiver

import "context"

type ClientBastion interface {
	RegisterApplication(ctx context.Context, application *Application) error
}
