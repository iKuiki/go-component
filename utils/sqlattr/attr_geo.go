package sqlattr

// 此文件下的struct主要为地理相关的struct

// GeoPoint 地理坐标
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
