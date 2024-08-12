package main

import (
	logging "github.com/ipfs/go-log"
	"github.com/urchinfs/go-urchin2-sdk/ipfs_api"
	"github.com/urchinfs/go-urchin2-sdk/utils"
)

func test() {
	var log = logging.Logger("main")
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

	//cid, err := client.Add("R-50.pkl")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("File uploaded successfully. CID: %s", cid)

	//cid, err := client.AddDir("code")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("Folder uploaded successfully. CID: %s", cid)

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
	//
	//err = car.UnpackCarFormat("tmp.car", "download")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	err := client.Get("QmU8UBwwik6iCn99VWKquPEezU9zQrJowTctpaokjtYqDa", "/root/test/down")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Infof("end.....................")
}
