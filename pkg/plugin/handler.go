package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/grafana/caic-datasource/pkg/caic"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// Provides everything the framework needs to handle CAIC requests
func DatasourceOpts(im instancemgmt.InstanceManager) datasource.ServeOpts {
	h := &Handler{
		im: im, //handler can instantiate datasource with the instance manager. The instance is whatever caicDataSourceInstance
	}

	return datasource.ServeOpts{
		QueryDataHandler:   h,
		CheckHealthHandler: h,
	}
}

// Handles calls to QueryData and CheckHealth
type Handler struct {
	im instancemgmt.InstanceManager
}

// Handles queries for CAIC Zone data
func (h *Handler) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	filter := struct {
		Zone caic.Region `json:"zone"`
	}{}

	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		log.DefaultLogger.Info(string(q.JSON))
		err := json.Unmarshal(q.JSON, &filter)
		if err != nil {
			return nil, errors.New(fmt.Sprint("bad query: ", err.Error()))
		}

		zoneData, err := h.queryZones(req, filter.Zone)
		if err != nil {
			return nil, err
		}

		response.Responses[q.RefID] = zoneData
	}

	return response, nil
}

func (h *Handler) queryZones(req *backend.QueryDataRequest, r caic.Region) (backend.DataResponse, error) {
	ds, err := h.datasource(req.PluginContext)
	if err != nil {
		return backend.DataResponse{}, err
	}

	if r == caic.EntireState {
		zones, err := ds.Client.StateSummary()
		if err != nil {
			return backend.DataResponse{}, err
		}
		return h.createResponse(zones), nil
	}

	zone, err := ds.Client.RegionSummary(r)
	if err != nil {
		return backend.DataResponse{}, err
	}
	return h.createResponse([]caic.Zone{zone}), nil
}

func (h *Handler) createResponse(zones []caic.Zone) backend.DataResponse {
	var names []string
	var rating []int64
	for _, z := range zones {
		names = append(names, z.Name)
		rating = append(rating, int64(z.Rating))

	}

	response := backend.DataResponse{}
	frame := data.NewFrame("Zones")
	frame.Fields = append(frame.Fields, data.NewField("name", nil, names))
	frame.Fields = append(frame.Fields, data.NewField("rating", nil, rating))
	response.Frames = append(response.Frames, frame)
	return response
}

func (h *Handler) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	ds, err := h.datasource(req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	if !ds.Client.CanConnect() {
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

func (h *Handler) datasource(pc backend.PluginContext) (*CaicDatasource, error) {
	i, err := h.im.Get(pc)
	if err != nil {
		return nil, err
	}

	ds, ok := i.(*CaicDatasource)
	if !ok {
		return nil, errors.New("bad datasource")
	}

	return ds, nil
}
