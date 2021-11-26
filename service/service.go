package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "imageelev/api/v1"
	img_db "imageelev/db"
)

type CapturedImageService struct {
	pb.UnimplementedImageByElevationServer
	lastRoutePath     pb.RoutePath
	lastCapturedImage pb.CapturedImage

	dbImages *gorm.DB
}

func HandleError(err interface{}, method string) {
	if pqc, ok := err.(*pgconn.PgError); ok {
		log.Printf("%v PostgreSQL connection error:%v\n", method, pqc)
	} else if pqerr, ok := err.(*pq.Error); ok {
		log.Printf("%v PostgreSQL error:%v\n", method, pqerr)
	}
}

func (ser *CapturedImageService) InitializeDB() error {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v",
		img_db.ServerIP, img_db.User, img_db.Pwd, img_db.DBName, img_db.Port)
	var err error
	ser.dbImages, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	if !ser.dbImages.Migrator().HasTable(&img_db.CapturedImage{}) {
		return errors.New("table " + img_db.CapturedImage{}.TableName() + " dosnt exist")
	}
	return nil
}

func (ser *CapturedImageService) GetRoutePath(req *pb.RoutePathRequest, serv pb.ImageByElevation_GetRoutePathServer) error {
	log.Printf("Receive %T : %v \n", req, req)
	if req.BeginTime == 0 && req.EndTime == 0 {
		return status.Errorf(codes.InvalidArgument, "GetRoutePath must be initialize request params")
	}
	var products []img_db.CapturedImage
	tBeg, tEnd := req.BeginTime, req.EndTime
	if tEnd == 0 {
		tEnd = uint32(time.Now().Unix())
	}
	err := ser.dbImages.Model(&img_db.CapturedImage{}).Where("time > ? AND time < ?", tBeg, tEnd).
		Find(&products).Error
	if err != nil {
		return err
	}
	route := pb.RoutePath{
		BeginTime: tBeg,
		EndTime:   tEnd,
	}
	if len(products) > 0 {
		route.Points = make([]*pb.RoutePath_RoutePoint, len(products))
		for i, img := range products {
			route.Points[i] = &pb.RoutePath_RoutePoint{
				GpsCoord:     &pb.GeoPoint{Lat: img.Coordinate.X, Long: img.Coordinate.Y},
				CapImageId:   uint32(img.ID),
				CapImageTime: uint32(img.Time),
			}
		}

	}

	if err = serv.Send(&route); err != nil {
		return err
	}
	return nil
}

func (ser *CapturedImageService) GetCapturedImage(ctx context.Context, req *pb.CapturedImageRequest) (*pb.CapturedImage, error) {
	// check exist metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Printf("Receive %T : %v with metadata: %v\n", req, req, md)
	} else {
		log.Printf("Receive %T : %v  without metadata \n", req, req)
	}
	var mp img_db.CapturedImage
	err := ser.dbImages.Model(&img_db.CapturedImage{}).First(&mp, "id = ?", req.Id).Error
	if err != nil {
		pqerr, ok := err.(*pq.Error)
		if ok {
			// the PostgreSQL error
			log.Printf("GetCapturedImage PostgreSQL error:%v\n", pqerr)
		}
		log.Printf("GetCapturedImage error:%v\n", err)

		st := status.New(codes.InvalidArgument, fmt.Sprintf("db-error: %v", err))
		badReq := errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "Id",
					Description: fmt.Sprintf("CapturedImage doesn`t exist with id:%v", req.Id),
				},
			},
		}
		st, _ = st.WithDetails(&badReq)
		return nil, st.Err()

	} else {
		res := &pb.CapturedImage{
			Id:   uint32(mp.ID),
			Time: uint32(mp.Time),
			GpsCoord: &pb.GeoPoint{
				Lat:  mp.Coordinate.X,
				Long: mp.Coordinate.Y,
			},
			CountAngle: uint32(mp.CountAngles),
			//Av:        make([]float32, len(mp.AV)),
			ElevationAngles: pb.ConverElevationsFromDB(mp.ElevationAngles),
			PathRaw:         mp.PathRaw,
			PathImage:       mp.PathImage,
		}
		ser.lastCapturedImage = *res
		return res, nil
	}
}

func (ser *CapturedImageService) GetRawImage(req *pb.RawImageRequest, serv pb.ImageByElevation_GetRawImageServer) error {
	log.Printf("Receive %T : %v \n", req, req)

	return nil
}

func (ser *CapturedImageService) CreateCapturedImage(ctx context.Context, req *pb.CreateCapturedImageRequest) (*pb.CapturedImage, error) {
	log.Printf("Receive %T : %v \n", req, req)

	if req.CapImage.GetTime() <= 0 || req.CapImage.GetGpsCoord() == nil {
		st := status.New(codes.InvalidArgument, "CreateCapturedImageRequest has incorrect params `time` or `coordinate`!\n")
		badReq := errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "time",
					Description: fmt.Sprintf("time must be greate 0, value:%v", req.CapImage.Time),
				},
				{
					Field:       "coordinate",
					Description: fmt.Sprintf("coordinate must be set, value:%v", req.CapImage.GetGpsCoord()),
				},
			},
		}
		st, _ = st.WithDetails(&badReq)
		return nil, st.Err()
	}
	//req.CapImage
	newImg := img_db.CapturedImage{
		ID:   0,
		Time: int64(req.CapImage.Time),
		Coordinate: img_db.Point{
			X: req.CapImage.GpsCoord.Lat,
			Y: req.CapImage.GpsCoord.Long,
		},
		CountAngles:     int64(req.CapImage.CountAngle),
		ElevationAngles: pb.ConverElevationsToDB(req.CapImage.ElevationAngles),
		PathRaw:         req.CapImage.PathRaw,
		PathImage:       req.CapImage.PathImage,
	}
	imgNew := ser.dbImages.Create(&newImg)
	if imgNew.Error != nil {
		stat := status.New(codes.InvalidArgument, imgNew.Error.Error())
		//TODO! add reuest body into stat.WithDetails
		return nil, stat.Err()
	} else {
		res := &pb.CapturedImage{
			Id:   uint32(newImg.ID),
			Time: uint32(newImg.Time),
			GpsCoord: &pb.GeoPoint{
				Lat:  newImg.Coordinate.X,
				Long: newImg.Coordinate.Y,
			},
			CountAngle:      uint32(newImg.CountAngles),
			ElevationAngles: pb.ConverElevationsFromDB(newImg.ElevationAngles),
			PathRaw:         newImg.PathRaw,
			PathImage:       newImg.PathImage,
		}
		return res, nil
	}
}

func NewService() *CapturedImageService {
	s := &CapturedImageService{}
	s.InitializeDB()
	return s
}
