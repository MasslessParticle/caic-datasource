package plugin

import (
	"context"
	"encoding/json"
	"errors"

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
	zones, err := h.getZones(req)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	filter := struct {
		Zone string `json:"zone"`
	}{}

	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries { //I'm unsure of multiple queries from my datasource?
		err := json.Unmarshal(q.JSON, &filter)
		if err != nil {
			return nil, err
		}

		zoneData := h.zonesToResponse(zones, filter.Zone)
		response.Responses[q.RefID] = zoneData
	}

	return response, nil
}

func (h *Handler) getZones(req *backend.QueryDataRequest) ([]caic.Zone, error) {
	ds, err := h.datasource(req.PluginContext)
	if err != nil {
		return nil, err
	}

	return ds.Client.StateSummary()
}

func (h *Handler) zonesToResponse(zones []caic.Zone, requestedZone string) backend.DataResponse {
	var names []string
	var rating []int64
	for _, z := range zones {
		if requestedZone == "entire-state" || z.ID == requestedZone {
			names = append(names, z.Name)
			rating = append(rating, int64(z.Rating))
		}
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
