package db

import (
	"testing"
	"github.com/crypto-url/utils"
	"github.com/crypto-url/utils/testutil"
)

var clusterNetWork  = utils.ClusterNetWork{
	OrderersNum: 4,
	Version: utils.FabricVersion(utils.NETWORK_TYPE_FABRIC_V1_2),
	Orgs:[]*utils.Organization{
		&utils.Organization{
			OrgName:"org1",
			Domain: "org1.example.com",
			PeersNum: 2,
		},
		&utils.Organization{
			OrgName:"org2",
			Domain: "org2.example.com",
			PeersNum: 2,
		},
	},
}

func TestInitDBHandler(t *testing.T) {
	defer func() {
		if err := recover(); err != nil{
			t.Fatal(err)
		}
	}()

	InitDBHandler("./")
}

func TestPutClusterNetWork(t *testing.T) {
	TestInitDBHandler(t)
	err := PutClusterNetWork("test", &clusterNetWork)
	if err != nil{t.Fatal(err.Error())}
}

func TestGetClusterNetWork(t *testing.T) {
	TestInitDBHandler(t)
	v, err := GetClusterNetWork("test")
	if err != nil{t.Fatal(err.Error())}

	testutil.AssertEquals(t, v.Version, clusterNetWork.Version)
	testutil.AssertEquals(t, v.OrderersNum, clusterNetWork.OrderersNum)
	testutil.AssertEquals(t, v.Orgs, clusterNetWork.Orgs)
}

func TestDeleteCluster(t *testing.T) {
	TestInitDBHandler(t)

	err := DeleteCluster("test")
	if err != nil{t.Fatal(err.Error())}

	_, err = GetClusterNetWork("test")
	if err == nil{t.Fatal("delete fauilre")}
}