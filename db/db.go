/// Descriptions types storing in Database
package db

import (
	"github.com/lib/pq"
)

var (
	DBName   string = "baltica_tour"
	ServerIP string = "localhost"
	Port     int32  = 5432
	User     string = "user"
	Pwd      string = "123qwe"
)

type CapturedImage struct {
	ID              uint64 `gorm:"primarykey"`
	Time            int64  `gorm:"primaryKey"`
	Coordinate      Point
	CountAngles     int64           `gorm:"column:countangles"`
	ElevationAngles pq.Float64Array `gorm:"type:float[]"`
	PathRaw         string
	PathImage       string
}

func (m CapturedImage) TableName() string {
	return "captured_images"
}
