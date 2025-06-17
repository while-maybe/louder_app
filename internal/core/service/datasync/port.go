package datasync

import "context"

type DataSynchroniser interface {
	SyncCountries(ctx context.Context) (int, int, error)
}
