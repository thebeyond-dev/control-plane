package launcher

import "context"

type Launcher interface {
	Launch(ctx context.Context)
}
