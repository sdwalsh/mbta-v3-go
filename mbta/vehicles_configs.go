package mbta

import (
	"net/http"
)

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
	Fields            []string                 // Fields to include with the response. Multiple fields MUST be a comma-separated (U+002C COMMA, “,”) list. Note that fields can also be selected for included data types
	IncludeTrip       bool                     // Include Trip data in response (The trip which the vehicle is currently operating)
	IncludeStop       bool                     // Include Stop data in response (The vehicle’s current (when current_status is StoppedAt) or next stop)
	IncludeRoute      bool                     // Include Route data in response (The one route that is designated for that trip, as in GTFS trips.txt. A trip might also provide service on other routes, identified by the MBTA’s multi_route_trips.txt GTFS extension. filter[route] does consider the multi_route_trips GTFS extension, so it is possible to filter for one route and get a different route included in the response)
	FilterIDs         []string                 // Filter by multiple IDs
	FilterTripIDs     []string                 // Filter by trip IDs
	FilterLabels      []string                 // Filter by label
	FilterRouteIDs    []string                 // Filter by route IDs. If the vehicle is on a multi-route trip, it will be returned for any of the routes
	FilterDirectionID string                   // Filter by Direction ID (Either "0" or "1")
	FilterRouteTypes  []string                 // Filter by route type(s)
}

func (config *GetAllVehiclesRequestConfig) addHTTPParamsToRequest(req *http.Request) {
	// Add fields and includes params to request
	includes := GetVehicleRequestConfig{
		Fields:       config.Fields,
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

// GetVehicleRequestConfig extra options for the GetVehicle request
type GetVehicleRequestConfig struct {
	Fields       []string // Fields to include with the response. Multiple fields MUST be a comma-separated (U+002C COMMA, “,”) list. Note that fields can also be selected for included data types
	IncludeTrip  bool     // Include Trip data in response (The trip which the vehicle is currently operating)
	IncludeStop  bool     // Include Stop data in response (The vehicle’s current (when current_status is StoppedAt) or next stop)
	IncludeRoute bool     // Include Route data in response (The one route that is designated for that trip, as in GTFS trips.txt. A trip might also provide service on other routes, identified by the MBTA’s multi_route_trips.txt GTFS extension. filter[route] does consider the multi_route_trips GTFS extension, so it is possible to filter for one route and get a different route included in the response)
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

	addCommaSeparatedListToQuery(q, "include", includes)
	addCommaSeparatedListToQuery(q, "fields[vehicle]", config.Fields)
	req.URL.RawQuery = q.Encode()
}
