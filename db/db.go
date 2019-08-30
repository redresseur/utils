package db

import (
	"github.com/crypto-url/utils"
	"github.com/crypto-url/db/leveldbhelper"
	"path/filepath"
	"encoding/json"
)


var dbHandler *leveldbhelper.DB

func InitDBHandler(path string)  {

	path = filepath.Join(path, "db")
	dbHandler = leveldbhelper.CreateDB(&leveldbhelper.Conf{path})
	dbHandler.Open()
	return
}

func GetClusterNetWork(clusterId string)(*utils.ClusterNetWork ,error) {

	res := &utils.ClusterNetWork{}
	hashId := utils.HashingAlgorithm(utils.HashAlgorithm_SHA256)([]byte(clusterId))
	if v , err := dbHandler.Get(hashId); err != nil{
		return nil, err
	}else {
		if v == nil{
			return nil, nil
		}

		return res, json.Unmarshal(v, res);
	}

	return nil, nil
}

func GetClusterNetWorkV1_1(clusterId string)(*utils.ClusterNetWorkV1_1 ,error) {

	res := &utils.ClusterNetWorkV1_1{}
	hashId := utils.HashingAlgorithm(utils.HashAlgorithm_SHA256)([]byte(clusterId))
	if v , err := dbHandler.Get(hashId); err != nil{
		return nil, err
	}else {
		if v == nil{
			return nil, nil
		}

		return res, json.Unmarshal(v, res);
	}

	return nil, nil
}

func PutClusterNetWork(clusterId string, work *utils.ClusterNetWork) error {
	if v , err := json.Marshal(work); err != nil{
		return err;
	}else {
		hashId := utils.HashingAlgorithm(utils.HashAlgorithm_SHA256)([]byte(clusterId))
		return  dbHandler.Put(hashId, v, true)
	}

	return nil
}

func PutClusterNetWorkV1_1(clusterId string, work *utils.ClusterNetWorkV1_1) error {
	if v , err := json.Marshal(work); err != nil{
		return err;
	}else {
		hashId := utils.HashingAlgorithm(utils.HashAlgorithm_SHA256)([]byte(clusterId))
		return  dbHandler.Put(hashId, v, true)
	}

	return nil
}

func DeleteCluster(clusterId string) error {
	hashId := utils.HashingAlgorithm(utils.HashAlgorithm_SHA256)([]byte(clusterId))
	return dbHandler.Delete(hashId, true)
}

func Close()  {
	dbHandler.Close()
}
