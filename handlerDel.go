package saavuu

import (
	"errors"
)

func (scvCtx *ServiceContext) delHandler() (data []byte, err error) {
	var (
		keyWithMyID string
	)
	if keyWithMyID, err = scvCtx.replaceUID(scvCtx.Key); err != nil {
		return nil, err
	} else if keyWithMyID == "" {
		return nil, errors.New("no key")

		//key must contain @me
	} else if keyWithMyID == scvCtx.Key || !(keyWithMyID[0] >= 'A' && keyWithMyID[0] <= 'Z') {
		return nil, errors.New("Unauthorized deletion!")
	}
	if scvCtx.Field, err = scvCtx.replaceUID(scvCtx.Field); err != nil {
		return nil, err
	}
	if scvCtx.Field == "" {
		cmd := Config.rds.Del(scvCtx.ctx, keyWithMyID)
		if err = cmd.Err(); err != nil {
			return nil, err
		}
		return []byte("key " + scvCtx.Key + " deleted"), nil
	}
	cmd := Config.rds.HDel(scvCtx.ctx, keyWithMyID, scvCtx.Field)
	if err = cmd.Err(); err != nil {
		return nil, err
	}
	return []byte("field " + scvCtx.Field + " deleted"), nil
}
