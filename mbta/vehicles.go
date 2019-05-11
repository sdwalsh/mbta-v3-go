package mbta

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const vehiclesAPIPath = "/vehicles"

// VehicleService service handling all of the vehicle related API calls
type VehicleService service

// GetAllVehiclesSortByType all of the possible ways to sort by for a GetAllVehicles request
type GetAllVehiclesSortByType string

const (
	GetAllVehiclesSortByBearingAscending              GetAllVehiclesSortByType = "bearing"
	GetAllVehiclesSortByBearingDescending             GetAllVehiclesSortByType = "-bearing"
	GetAllVehiclesSortByCurrentStatusAscending        GetAllVehiclesSortByType = "current_status"
	GetAllVehiclesSortByCurrentStatusDescending       GetAllVehiclesSortByType = "-current_status"
	GetAllVehiclesSortByCurrentStopSequenceAscending  GetAllVehiclesSortByType = "current_stop_sequence"
	GetAllVehiclesSortByCurrentStopSequenceDescending GetAllVehiclesSortByType = "-current_stop_sequence"
	GetAllVehiclesSortByDirectionIDAscending          GetAllVehiclesSortByType = "direction_id"
	GetAllVehiclesSortByDirectionIDDescending         GetAllVehiclesSortByType = "-direction_id"
	GetAllVehiclesSortByLabelAscending                GetAllVehiclesSortByType = "label"
	GetAllVehiclesSortByLabelDescending               GetAllVehiclesSortByType = "-label"
	GetAllVehiclesSortByLatitudeAscending             GetAllVehiclesSortByType = "latitude"
	GetAllVehiclesSortByLatitudeDescending            GetAllVehiclesSortByType = "-latitude"
	GetAllVehiclesSortByLongitudeAscending            GetAllVehiclesSortByType = "longitude"
	GetAllVehiclesSortByLongitudeDescending           GetAllVehiclesSortByType = "-longitude"
	GetAllVehiclesSortBySpeedAscending                GetAllVehiclesSortByType = "speed"
	GetAllVehiclesSortBySpeedDescending               GetAllVehiclesSortByType = "-speed"
	GetAllVehiclesSortByUpdatedAtAscending            GetAllVehiclesSortByType = "updated_at"
	GetAllVehiclesSortByUpdatedAtDescending           GetAllVehiclesSortByType = "-updated_at"
)

// GetAllVehiclesRequestConfig extra options for the GetAllVehicles request
type GetAllVehiclesRequestConfig struct {
	PageOffset        string                   // Offset (0-based) of first element in the page
	PageLimit         string                   // Max number of elements to return
	Sort              GetAllVehiclesSortByType // Results can be sorted by the id or any GetAllVehiclesSortByType
	IncludeTrip       bool                     // Include Trip data in response
	IncludeStop       bool                     // Include Stop data in response
	IncludeRoute      bool                     // Include Route data in response
	FilterIDs         []string                 // Filter by multiple IDs
	FilterTripIDs     []string                 // Filter by trip IDs
	FilterLabels      []string                 // Filter by label
	FilterRouteIDs    []string                 // Filter by route IDs. If the vehicle is on a multi-route trip, it will be returned for any of the routes
	FilterDirectionID string                   // Filter by Direction ID (Either "0" or "1")
	FilterRouteTypes  []string                 // Filter by route type(s)
}

func (config *GetAllVehiclesRequestConfig) addHTTPParamsToRequest(req *http.Request) {
	// Add includes params to request
	includes := GetVehicleRequestConfig{
		IncludeTrip:  config.IncludeTrip,
		IncludeStop:  config.IncludeStop,
		IncludeRoute: config.IncludeRoute,
	}
	includes.addHTTPParamsToRequest(req)

	q := req.URL.Query()
	addToQuery(q, "page[offset]", config.PageOffset)
	addToQuery(q, "page[limit]", config.PageLimit)
	addToQuery(q, "sort", string(config.Sort))
	addToQuery(q, "filter[direction_id]", config.FilterDirectionID)
	addCommaSeparatedListToQuery(q, "filter[id]", config.FilterIDs)
	addCommaSeparatedListToQuery(q, "filter[trip]", config.FilterTripIDs)
	addCommaSeparatedListToQuery(q, "filter[label]", config.FilterLabels)
	addCommaSeparatedListToQuery(q, "filter[route]", config.FilterRouteIDs)
	addCommaSeparatedListToQuery(q, "filter[route_type]", config.FilterRouteTypes)

	req.URL.RawQuery = q.Encode()
}

// GetAllVehicles returns all vehicles from the mbta API
func (s *VehicleService) GetAllVehicles(config GetAllVehiclesRequestConfig) ([]Vehicle, error) {
	return s.GetAllVehiclesContext(context.Background(), config)
}

// GetAllVehiclesContext returns all vehicles from the mbta API given a context
func (s *VehicleService) GetAllVehiclesContext(ctx context.Context, config GetAllVehiclesRequestConfig) ([]Vehicle, error) {
	req, err := s.client.newRequest("GET", vehiclesAPIPath, nil)
	config.addHTTPParamsToRequest(req)
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}

	var vehicles []Vehicle
	_, err = s.client.do(req, &vehicles)
	return vehicles, err
}

// GetVehicleRequestConfig extra options for the GetVehicle request
type GetVehicleRequestConfig struct {
	IncludeTrip  bool
	IncludeStop  bool
	IncludeRoute bool
}

func (config *GetVehicleRequestConfig) addHTTPParamsToRequest(req *http.Request) {
	q := req.URL.Query()

	includes := []string{}
	if config.IncludeTrip {
		includes = append(includes, "trip")
	}
	if config.IncludeStop {
		includes = append(includes, "stop")
	}
	if config.IncludeRoute {
		includes = append(includes, "route")
	}
	if len(includes) > 0 {
		includesString := strings.Join(includes, ",")
		q.Add("include", includesString)
		req.URL.RawQuery = q.Encode()
	}
}

// GetVehicle returns a vehicle from the mbta API
func (s *VehicleService) GetVehicle(id string, config GetVehicleRequestConfig) (Vehicle, error) {
	return s.GetVehicleContext(context.Background(), id, config)
}

// GetVehicleContext returns a vehicle from the mbta API given a context
func (s *VehicleService) GetVehicleContext(ctx context.Context, id string, config GetVehicleRequestConfig) (Vehicle, error) {
	path := fmt.Sprintf("/%s/%s", vehiclesAPIPath, id)
	req, err := s.client.newRequest("GET", path, nil)
	config.addHTTPParamsToRequest(req)
	req = req.WithContext(ctx)
	if err != nil {
		return Vehicle{}, err
	}

	var vehicle Vehicle
	_, err = s.client.do(req, &vehicle)
	return vehicle, err
}