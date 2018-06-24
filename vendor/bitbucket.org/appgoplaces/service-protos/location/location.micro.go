// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: bitbucket.org/appgoplaces/service-protos/location/location.proto

/*
Package location is a generated protocol buffer package.

It is generated from these files:
	bitbucket.org/appgoplaces/service-protos/location/location.proto

It has these top-level messages:
	QueryLocationRequest
	PlaceResults
	QueryLocationResponse
	ExtractPlaceDataRequest
	ExtractPlaceDataResponse
	LatLng
	LatLngBounds
	Geometry
	OpeningHoursOpenClose
	OpeningHoursPeriod
	OpeningHours
	Photo
	AltID
	PlacesSearchResult
	QueryRequest
	Place
	QueryResponse
	SuggestCityRequest
	SearchCityCountry
	SuggestCityResponse
*/
package location

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Location service

type LocationService interface {
	SuggestCity(ctx context.Context, in *SuggestCityRequest, opts ...client.CallOption) (*SuggestCityResponse, error)
	QueryLocations(ctx context.Context, in *QueryLocationRequest, opts ...client.CallOption) (*QueryLocationResponse, error)
	ExtractPlaceData(ctx context.Context, in *ExtractPlaceDataRequest, opts ...client.CallOption) (*ExtractPlaceDataResponse, error)
}

type locationService struct {
	c           client.Client
	serviceName string
}

func LocationServiceClient(serviceName string, c client.Client) LocationService {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "location"
	}
	return &locationService{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *locationService) SuggestCity(ctx context.Context, in *SuggestCityRequest, opts ...client.CallOption) (*SuggestCityResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Location.SuggestCity", in)
	out := new(SuggestCityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *locationService) QueryLocations(ctx context.Context, in *QueryLocationRequest, opts ...client.CallOption) (*QueryLocationResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Location.QueryLocations", in)
	out := new(QueryLocationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *locationService) ExtractPlaceData(ctx context.Context, in *ExtractPlaceDataRequest, opts ...client.CallOption) (*ExtractPlaceDataResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Location.ExtractPlaceData", in)
	out := new(ExtractPlaceDataResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Location service

type LocationHandler interface {
	SuggestCity(context.Context, *SuggestCityRequest, *SuggestCityResponse) error
	QueryLocations(context.Context, *QueryLocationRequest, *QueryLocationResponse) error
	ExtractPlaceData(context.Context, *ExtractPlaceDataRequest, *ExtractPlaceDataResponse) error
}

func RegisterLocationHandler(s server.Server, hdlr LocationHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&Location{hdlr}, opts...))
}

type Location struct {
	LocationHandler
}

func (h *Location) SuggestCity(ctx context.Context, in *SuggestCityRequest, out *SuggestCityResponse) error {
	return h.LocationHandler.SuggestCity(ctx, in, out)
}

func (h *Location) QueryLocations(ctx context.Context, in *QueryLocationRequest, out *QueryLocationResponse) error {
	return h.LocationHandler.QueryLocations(ctx, in, out)
}

func (h *Location) ExtractPlaceData(ctx context.Context, in *ExtractPlaceDataRequest, out *ExtractPlaceDataResponse) error {
	return h.LocationHandler.ExtractPlaceData(ctx, in, out)
}