package service

import (
	"context"
	"fmt"
	"net"
	"strings"

	proto_geoip_service "github.com/FREEGREAT/protos/gen/go/geoip"
	"github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"
)

type GeoIpService struct {
	proto_geoip_service.UnimplementedGeoIpServiceServer
}

func CreateNewGeoIpService() *GeoIpService {

	return &GeoIpService{}
}

var country, city string

func (s *GeoIpService) GetLocationByIP(ctx context.Context, req *proto_geoip_service.GetLocationRequest) (*proto_geoip_service.GetLocationResponse, error) {
	logrus.Info("Alyo@!")
	ip := req.Ip
	logrus.Infof("IP address: %s", ip)
	if ip == "127.0.0.1" || strings.HasPrefix(ip, "192.168.") {
		country = "Local"
		city = "Local"
	}
	db, err := geoip2.Open("db/GeoLite2-City.mmdb")
	if err != nil {
		logrus.Errorf("Failed to open GeoIP database: %v", err)
		return nil, fmt.Errorf("failed to open GeoIP database: %v", err)
	}
	defer db.Close()

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		logrus.Errorf("Invalid IP address: %s", ip)
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}
	record, err := db.City(parsedIP)
	if err != nil {
		logrus.Errorf("Failed to get city info: %v", err)
		return nil, fmt.Errorf("failed to get city info: %v", err)
	}
	if record != nil {
		city = record.City.Names["en"]
		country = record.Country.Names["en"]
	}
	resp := proto_geoip_service.GetLocationResponse{
		Location: country + ", " + city,
	}
	logrus.Info("Konechnaya@!")
	return &resp, nil

}
