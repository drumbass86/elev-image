package service_test

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	pb "imageelev/api/v1"
	img_db "imageelev/db"
	"imageelev/service"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	pbjson "google.golang.org/protobuf/encoding/protojson"
)

var hostDB = flag.String("host", "", "Host DB server")

// Mocking struct for testing server method
type mockImageByElevation_GetRoutePathServer struct {
	grpc.ServerStream
	Result *pb.RoutePath
}

func (_m *mockImageByElevation_GetRoutePathServer) Send(route *pb.RoutePath) error {
	_m.Result = route
	return nil
}

type CapturedImageServiceTestSuite struct {
	suite.Suite
	serv *service.CapturedImageService
}

func (ts *CapturedImageServiceTestSuite) SetupSuite() {
	ts.serv = &service.CapturedImageService{}

	if hostDB != nil && len(*hostDB) > 0 {
		img_db.ServerIP = *hostDB
	}
	if err := ts.createTestDB(); err != nil {
		ts.T().Fatalf("SetupTest InitializeDB return error:%v", err)
	}
}

func (ts *CapturedImageServiceTestSuite) createTestDB() error {
	//saw file starting docker container "./tools/run_test_postgres_docker.sh"
	img_db.DBName = "test"
	img_db.User = "user"
	img_db.Pwd = "test"
	img_db.Port = 5433
	img_db.ServerIP = "localhost"

	if err := ts.serv.InitializeDB(); err != nil {
		return err
	} else {
		return nil
	}
}

func (ts *CapturedImageServiceTestSuite) TestGetRoutePath() {
	req := pb.RoutePathRequest{BeginTime: 0, EndTime: uint32(time.Now().Unix())}
	resp := mockImageByElevation_GetRoutePathServer{}
	err := ts.serv.GetRoutePath(&req, &resp)
	if err != nil {
		ts.T().Errorf("TestGetRoutePath request(timeBeg:%v timeEnd:%v) got unexpected error", req.BeginTime, req.EndTime)
	} else {
		ts.NotNil(resp.Result, "RoutePath must be initialize")
		ts.Condition(func() bool {
			return resp.Result.Points != nil && len(resp.Result.Points) > 0
		},
			"RoutePath must contain point",
		)
	}

}

func (ts *CapturedImageServiceTestSuite) TestGetCapturedImage() {
	req := pb.CapturedImageRequest{
		Id: 1,
	}
	ctx := context.Background()
	resp, err := ts.serv.GetCapturedImage(ctx, &req)
	if err != nil {
		ts.T().Errorf("TestGetCapturedImage req:%v error:%v", req.Id, err)
	} else {
		ts.NotNil(resp, "CapturedImage must exist")
		ts.Equal(req.Id, resp.Id, "CapturedImage.Id must be equal Request.Id")
		ts.NotNil(resp.Time, "CapturedImage.Time must be not nil")
		ts.NotNil(resp.GpsCoord, "CapturedImage.GpsCoord must be not nil")
		ts.NotNil(resp.CountAngle, "CapturedImage.CountAngle must be not nil")
		ts.NotEmpty(resp.ElevationAngles, "CapturedImage.ElevationAngles must be not empty")
	}
}

func (ts *CapturedImageServiceTestSuite) TestCreateCapturedImage() {
	ctx := context.Background()
	req := pb.CreateCapturedImageRequest{
		CapImage: &pb.CapturedImage{},
	}
	testTime := uint32(time.Now().Unix())
	const testCountAngles = uint32(14)
	jsonReq := []byte(
		fmt.Sprintf(`{
			"cap_image":{
				"elevation_angles": [
					-1.96,
					0,
					2.05,
					3.02,
					4.07,
					5.04,
					6.00,
					9.06,
					11.03,
					13.99,
					18.05,
					25.01,
					35.07,
					62.04
				],
				"count_angle": %v,
				"path_raw": "/20210723/img_raw_%v.dat",
				"path_image": "img_%v.png",
				"gps_coord": {
					"lat": 54.8330192565918,
					"long": 19.303146362304688    
				},    
				"time": %v,
				"time_desc": "23:07:2021 11:30:10"
			}    
		}`, testCountAngles, testTime, testTime, testTime),
	)
	unmarsh := pbjson.UnmarshalOptions{
		//AllowPartial:   true,
		DiscardUnknown: true,
	}
	if err := unmarsh.Unmarshal(jsonReq, &req); err != nil {
		ts.T().Errorf("TestCreateCapturedImage can`t unmarshal test json request, err:%v\n", err)
	} else if req.CapImage == nil {
		ts.T().Error("TestCreateCapturedImage after unmarshal test json request is empty!")
	} else {
		resp, err := ts.serv.CreateCapturedImage(ctx, &req)
		if err != nil {
			ts.T().Errorf("CreateCapturedImage return error:%v\n", err)
		} else {
			ts.GreaterOrEqual(resp.Id, uint32(1), "CapturedImage id must init >= 1")
			ts.Equal(resp.Time, testTime, "CapturedImage time not equal initial value")
			ts.Equal(resp.CountAngle, testCountAngles, "CapturedImage CountAngle not equal initial value")

			ts.T().Logf("Create CapturedImage: {%v}\n", resp)
		}
	}

}

func TestCapturedImageServiceSuite(t *testing.T) {
	suite.Run(t, new(CapturedImageServiceTestSuite))
}
