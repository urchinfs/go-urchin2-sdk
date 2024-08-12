package cid

import (
	"context"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	"github.com/urchinfs/go-urchin2-sdk/utils"
	"os"
)

var log = logging.Logger("cid")

func GetCid(path string) (cid.Cid, error) {
	_, err := os.Stat(path)
	if err != nil {
		log.Errorf("get path stat err:%v", err)
		return cid.Cid{}, err
	}

	wrapDataDir, err := utils.WarpPath(path)
	if err != nil {
		return cid.Cid{}, err
	}

	rootCid, err := AddAndBuildCid(context.Background(), wrapDataDir)
	if err != nil {
		return cid.Cid{}, err
	}

	log.Infof("get cid:%s", rootCid)
	return rootCid, nil
}
