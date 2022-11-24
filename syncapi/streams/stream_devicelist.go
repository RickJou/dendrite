package streams

import (
	"context"

	keyapi "github.com/RickJou/dendrite/keyserver/api"
	"github.com/RickJou/dendrite/roomserver/api"
	"github.com/RickJou/dendrite/syncapi/internal"
	"github.com/RickJou/dendrite/syncapi/storage"
	"github.com/RickJou/dendrite/syncapi/types"
)

type DeviceListStreamProvider struct {
	DefaultStreamProvider
	rsAPI  api.SyncRoomserverAPI
	keyAPI keyapi.SyncKeyAPI
}

func (p *DeviceListStreamProvider) CompleteSync(
	ctx context.Context,
	snapshot storage.DatabaseTransaction,
	req *types.SyncRequest,
) types.StreamPosition {
	return p.LatestPosition(ctx)
}

func (p *DeviceListStreamProvider) IncrementalSync(
	ctx context.Context,
	snapshot storage.DatabaseTransaction,
	req *types.SyncRequest,
	from, to types.StreamPosition,
) types.StreamPosition {
	var err error
	to, _, err = internal.DeviceListCatchup(context.Background(), snapshot, p.keyAPI, p.rsAPI, req.Device.UserID, req.Response, from, to)
	if err != nil {
		req.Log.WithError(err).Error("internal.DeviceListCatchup failed")
		return from
	}
	err = internal.DeviceOTKCounts(req.Context, p.keyAPI, req.Device.UserID, req.Device.ID, req.Response)
	if err != nil {
		req.Log.WithError(err).Error("internal.DeviceListCatchup failed")
		return from
	}

	return to
}
