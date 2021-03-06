package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type caicClient interface {
	CanConnect() bool
	Summary(caic.Region) ([]caic.Zone, error)
	AspectDanger(caic.Region) (caic.AspectDanger, error)
}

// Handles calls to QueryData and CheckHealth
type Handler struct {
	Client caicClient
}

// Handles queries for CAIC Zone data
func (h *Handler) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	filter := struct {
		Zone caic.Region `json:"zone"`
	}{}

	qr := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		err := json.Unmarshal(q.JSON, &filter)
		if err != nil {
			return nil, errors.New(fmt.Sprint("bad query: ", err.Error()))
		}

		zoneFrame, err := h.queryZones(req, filter.Zone)
		if err != nil {
			return nil, err
		}

		problemFrame, err := h.queryProblems(req, filter.Zone)
		if err != nil {
			return nil, err
		}

		resp := backend.DataResponse{}
		resp.Frames = append(resp.Frames, zoneFrame)
		resp.Frames = append(resp.Frames, problemFrame)

		qr.Responses[q.RefID] = resp
	}

	return qr, nil
}

func (h *Handler) queryZones(req *backend.QueryDataRequest, r caic.Region) (*data.Frame, error) {
	zones, err := h.Client.Summary(r)
	if err != nil {
		return nil, err
	}
	return h.createResponse(zones), nil
}

func (h *Handler) queryProblems(req *backend.QueryDataRequest, r caic.Region) (*data.Frame, error) {
	aspectDanger, err := h.Client.AspectDanger(r)
	if err != nil {
		return nil, err
	}

	ordinals := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	degrees := []int32{0, 45, 90, 135, 180, 225, 270, 315}

	aboveTreeline := []int32{
		toInt(aspectDanger.AboveTreeline.North),
		toInt(aspectDanger.AboveTreeline.NorthEast),
		toInt(aspectDanger.AboveTreeline.East),
		toInt(aspectDanger.AboveTreeline.SouthEast),
		toInt(aspectDanger.AboveTreeline.South),
		toInt(aspectDanger.AboveTreeline.SouthWest),
		toInt(aspectDanger.AboveTreeline.West),
		toInt(aspectDanger.AboveTreeline.NorthWest),
	}

	nearTreeline := []int32{
		toInt(aspectDanger.NearTreeline.North),
		toInt(aspectDanger.NearTreeline.NorthEast),
		toInt(aspectDanger.NearTreeline.East),
		toInt(aspectDanger.NearTreeline.SouthEast),
		toInt(aspectDanger.NearTreeline.South),
		toInt(aspectDanger.NearTreeline.SouthWest),
		toInt(aspectDanger.NearTreeline.West),
		toInt(aspectDanger.NearTreeline.NorthWest),
	}

	belowTreeline := []int32{
		toInt(aspectDanger.BelowTreeline.North),
		toInt(aspectDanger.BelowTreeline.NorthEast),
		toInt(aspectDanger.BelowTreeline.East),
		toInt(aspectDanger.BelowTreeline.SouthEast),
		toInt(aspectDanger.BelowTreeline.South),
		toInt(aspectDanger.BelowTreeline.SouthWest),
		toInt(aspectDanger.BelowTreeline.West),
		toInt(aspectDanger.BelowTreeline.NorthWest),
	}

	frame := data.NewFrame("AspectDanger")
	frame.Fields = append(frame.Fields, data.NewField("ordinals", nil, ordinals))
	frame.Fields = append(frame.Fields, data.NewField("degrees", nil, degrees))
	frame.Fields = append(frame.Fields, data.NewField("aboveTreeline", nil, aboveTreeline))
	frame.Fields = append(frame.Fields, data.NewField("nearTreeline", nil, nearTreeline))
	frame.Fields = append(frame.Fields, data.NewField("belowTreeline", nil, belowTreeline))

	return frame, nil
}

func (h *Handler) createResponse(zones []caic.Zone) *data.Frame {
	var names []string
	var rating []int64
	var aboveTreeline []int64
	var nearTreeline []int64
	var belowTreeline []int64
	for _, z := range zones {
		names = append(names, z.Name)
		rating = append(rating, int64(z.Rating))
		aboveTreeline = append(aboveTreeline, int64(z.AboveTreeline))
		nearTreeline = append(nearTreeline, int64(z.NearTreeline))
		belowTreeline = append(belowTreeline, int64(z.BelowTreeline))
	}

	frame := data.NewFrame("Zones")
	frame.Fields = append(frame.Fields, data.NewField("name", nil, names))
	frame.Fields = append(frame.Fields, data.NewField("rating", nil, rating))
	frame.Fields = append(frame.Fields, data.NewField("aboveTreeline", nil, aboveTreeline))
	frame.Fields = append(frame.Fields, data.NewField("nearTreeline", nil, nearTreeline))
	frame.Fields = append(frame.Fields, data.NewField("belowTreeline", nil, belowTreeline))
	return frame
}

func (h *Handler) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if !h.Client.CanConnect() {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Error reaching CAIC site",
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}

func toInt(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
