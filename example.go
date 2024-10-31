package ipfs_api

import (
	"context"
	logging "github.com/ipfs/go-log"
	"github.com/urchinfs/go-urchin2-sdk/ipfs_api"
	"github.com/urchinfs/go-urchin2-sdk/utils"
)

func example() {
	var log = logging.Logger("example")
	utils.InitLog()

	log.Infof("start...")
	client := ipfs_api.NewClient("192.168.242.42:5001")
	//client := ipfs_api.NewClient("127.0.0.1:5001")
	//inputFile := "init_model.7z"
	//outputFile := "init_model.7z.car"

	//inputDir := "code"
	//outputDir := "code.car"
	//
	//err := client.SwarmConnect(context.Background(), "/ip4/192.168.1.1/tcp/4001/ipfs/12D3KooWGA2h89gV5sqCH3wPxhgguzBDNsAbaHocxS6z7vmEhB1V")
	//if err != nil {
	//	log.Infof("SwarmConnect err: %v", err)
	//	return
	//}
	//log.Infof("Swarm Connect succeed")

	/*
	*
	* upload.................................................
	*
	 */
	//getCid, err := myCid.GetCid(inputDir)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//log.Printf("myCid.GetCid. CID: %s", getCid)
	//
	//rootCid, err := car.PackCarFormat(inputDir, outputDir)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//log.Printf("pack car file done, CID: %s", rootCid)
	//
	//dagImport, err := client.DagImport(outputDir, false, false)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//log.Printf("dag import to peer done, result: %v", dagImport)

	//- add file
	//cid, err := client.Add(context.Background(), "R-50.pth")
	//cid, err := client.Add(context.Background(), "R-50.pth", ipfs_api.CenterId(2))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Infof("File uploaded successfully. CID: %s", cid)

	//cid, err := client.AddDir(context.Background(), "E:\\Exchange_dir\\tmp\\ipfs\\data\\code", ipfs_api.CenterId(2))
	//cid, err := client.AddDir(context.Background(), "E:\\Exchange_dir\\tmp\\ipfs\\data\\code")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Infof("Folder uploaded successfully. CID: %s", cid)

	/*
	*
	* download.................................................
	*
	 */
	//err := client.DagExport("QmZiM8GY7p9rCiRRJmnfa7KrXoWgJ6NbgMzvvgpz3qWUbc", "tmp.car")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	//err := car.UnpackCarFormat("tmp.car", "./ipfs")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	//err := client.Get(context.Background(), "QmU8UBwwik6iCn99VWKquPEezU9zQrJowTctpaokjtYqDa", "/root/test/down")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	//centers, err := client.CidQuery(context.Background(), "QmdvvpjrAw2nPDnNuCuzsj34dQGSDipYyNvSBLL6dRb154")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//log.Infof("centers:%v", centers)

	//peers, err := client.SwarmPeers(context.Background())
	//if err != nil {
	//	return
	//}
	//log.Infof("peers:%v", peers)

	//peer, err := client.PeerSelf(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//log.Infof("peer:%v", peer)

	//peer, err := client.PeerQuery(context.Background(), "12D3KooWB6MLtPV9rADy6TdCNXuDV6D779pZxiG59hCQaLyqTFFS")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//log.Infof("peer:%v", peer)

	peers, err := client.PeerAll(context.Background())
	if err != nil {
		return
	}
	log.Infof("peers:%v", peers)

	//exists, err := client.CidExistInCenter(context.Background(), "QmdvvpjrAw2nPDnNuCuzsj34dQGSDipYyNvSBLL6dRb154", 2)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//log.Infof("exists:%v", exists)

	//- use add interface first
	//err := client.CidSync(context.Background(), "QmdvvpjrAw2nPDnNuCuzsj34dQGSDipYyNvSBLL6dRb154", 2)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	//status, err := client.CheckSyncStatus(context.Background(), "QmdvvpjrAw2nPDnNuCuzsj34dQGSDipYyNvSBLL6dRb154", 2)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//log.Infof("status:%v", status)

	//err := client.CidMigrate(context.Background(), "QmPsPJuLkrF5F75iebbYXRT35GtrvYzTD8m1uttCMtaVnr", 4, "test/tmp")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	//status, err := client.CheckMigrateStatus(context.Background(), "QmPsPJuLkrF5F75iebbYXRT35GtrvYzTD8m1uttCMtaVnr", 4, "test/tmp")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//log.Infof("status:%v", status)

	log.Infof("end.....................")
}
