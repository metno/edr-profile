package forecastts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	baseURL string
}

// implements ServerInterface
func NewHandler(baseURL string) *Handler {
	return &Handler{baseURL: baseURL}
}

// landing page of this API
// (GET /)
func (h *Handler) GetLandingPage(ctx echo.Context, params GetLandingPageParams) error {
	ret := LandingPage{
		Title:       ptr("Sample edr service"),
		Description: ptr("A sample edr service compliant with forecast timeseries profile, implemented using ogen."),
		Links: []Link{
			{
				Href:  h.baseURL,
				Rel:   "self",
				Type:  ptr("application/json"),
				Title: ptr("this document"),
			},
			{
				Href:  h.baseURL + "/api",
				Rel:   "service-desc",
				Type:  ptr("application/openapi+json;version=3.0"),
				Title: ptr("the API definition"),
			},
			{
				Href:  h.baseURL + "/conformance",
				Rel:   "conformance",
				Type:  ptr("application/json"),
				Title: ptr("OGC conformance classes implemented by this API"),
			},
			{
				Href:  h.baseURL + "/collections",
				Rel:   "data",
				Title: ptr("Metadata about the resource collections"),
			},
		},
		Keywords: &[]string{"meteorology", "test"},
		Provider: &struct {
			Name *string "json:\"name,omitempty\""
			Url  *string "json:\"url,omitempty\""
		}{
			Name: ptr("MET Norway"),
			Url:  ptr("https://met.no"),
		},
		// Contact
	}

	return writer(params.F, ctx)(http.StatusOK, &ret)
}

// List the available collections from the service
// (GET /collections)
func (h *Handler) ListCollections(ctx echo.Context, params ListCollectionsParams) error {
	var collections []Collection
	for _, collectionName := range []string{"MEPS"} {
		collection, err := h.getQueries(collectionName)
		if err != nil {
			return writer(params.F, ctx)(http.StatusInternalServerError,
				Exception{
					Code:        "internal error",
					Description: ptr("Internal server error"),
				},
			)
		}
		collections = append(collections, *collection)
	}

	ret := Collections{
		Links: []Link{
			{
				Href: h.baseURL + "/edr/collections/",
				Rel:  "self",
				Type: ptr("application/json"),
			},
		},
		Collections: collections,
	}

	return writer(params.F, ctx)(http.StatusOK, &ret)
}

// List query types supported by the collection
// (GET /collections/{collectionId})
func (h *Handler) GetQueries(ctx echo.Context, collectionId CollectionId, params GetQueriesParams) error {
	const id = "MEPS"

	if collectionId != id {
		return writer(params.F, ctx)(http.StatusNotFound,
			Exception{
				Code:        "not found",
				Description: ptr("No such collection id."),
			},
		)
	}

	ret, err := h.getQueries(collectionId)
	if err != nil {
		return writer(params.F, ctx)(http.StatusInternalServerError,
			Exception{
				Code:        "internal error",
				Description: ptr("Internal server error"),
			},
		)
	}

	return writer(params.F, ctx)(http.StatusOK, &ret)
}

func (h *Handler) getQueries(collectionId CollectionId) (*Collection, error) {
	return &Collection{
		// Links []GetQueriesOKApplicationJSONLinksItem `json:"links"`
		Id:       collectionId,
		Title:    ptr("MEPS"),
		Keywords: &[]string{"forecast", "timeseries", "nordic", "air_temperature"},
		Links: []Link{
			{
				Href: fmt.Sprintf("%s/collections/%s", h.baseURL, collectionId),
				Rel:  "data",
				Type: ptr("application/geo+json"),
			},
			{
				Href: fmt.Sprintf("%s/collections/%s/position", h.baseURL, collectionId),
				Rel:  "data",
				Type: ptr("application/vnd.cov+json"),
			},
		},
		Extent: Extent{
			Spatial: &struct {
				Bbox []Extent_Spatial_Bbox_Item "json:\"bbox\""
				Crs  string                     "json:\"crs\""
				Name *string                    "json:\"name,omitempty\""
			}{
				Bbox: []Extent_Spatial_Bbox_Item{
					func() Extent_Spatial_Bbox_Item {
						var ret Extent_Spatial_Bbox_Item
						ret.FromExtentSpatialBbox0([]float32{-180, 90, 180, -90})
						return ret
					}(),
				},
				Crs: "GEOGCS[\"WGS 84\",DATUM[\"WGS_1984\",SPHEROID[\"WGS 84\",6378137,298.257223563,AUTHORITY[\"EPSG\",\"7030\"]],AUTHORITY[\"EPSG\",\"6326\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.01745329251994328,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4326\"]]",
			},
			Temporal: &struct {
				// Interval RFC3339 compliant Date and Time
				Interval [][]time.Time `json:"interval"`

				// Name Name of the temporal coordinate reference system
				Name *string `json:"name,omitempty"`
				Trs  string  `json:"trs"`

				// Values Provides information about the time intervals available in the collection
				// as ISO8601 compliant dates, either as a time range specified
				// as start time / end time  (e.g. 2017-11-14T09:00Z/2017-11-14T21:00Z)  or
				// as number of repetitions / start time / interval (e.g. R4/2017-11-14T21:00Z/PT3H)
				// or a list of time values (e.g.
				// 2017-11-14T09:00Z,2017-11-14T12:00Z,2017-11-14T15:00Z,2017-11-14T18:00Z,2017-11-14T21:00Z)
				Values *[]time.Time `json:"values,omitempty"`
			}{
				Interval: [][]time.Time{{startTime, endTime}},
				Values:   &[]time.Time{startTime, endTime},
				Trs:      "http://www.opengis.net/def/uom/ISO-8601/0/Gregorian",
			},
			Vertical: &struct {
				Interval [][]string "json:\"interval\""
				Name     *string    "json:\"name,omitempty\""
				Values   *[]string  "json:\"values,omitempty\""
				Vrs      string     "json:\"vrs\""
			}{
				Vrs:      `PARAMETRICCRS["WMO standard atmosphere layer 0",PDATUM["Mean Sea Level",ANCHOR["101325 Pa at 15°C"]],CS[parametric,1],AXIS["pressure (Pa)",up],PARAMETRICUNIT["Pascal",1.0]]`,
				Interval: [][]string{{"100000", "50000"}},
				Values:   &[]string{"100000", "50000"},
			},
		},
		// // Detailed information relevant to individual query types.
		DataQueries: struct {
			Items *struct {
				Link *ItemsLink `json:"link,omitempty"`
			} `json:"items,omitempty"`
			Locations *struct {
				Link *LocationsLink `json:"link,omitempty"`
			} `json:"locations,omitempty"`
			Position struct {
				Link *PositionLink `json:"link,omitempty"`
			} `json:"position"`
		}{
			Position: struct {
				Link *PositionLink `json:"link,omitempty"`
			}{
				Link: &Link{
					Href:      fmt.Sprintf("%s/collections/%s/position?coords={coords}", h.baseURL, collectionId),
					Rel:       "data",
					Templated: ptr(true),
				},
			},
		},

		Crs:           []string{"CRS84"},
		OutputFormats: []string{"CoverageJSON"},

		ParameterNames: map[string]interface{}{
			"air_temperature": []byte("{}"),
		},
	}, nil
}

// List data instances of {collectionId}
// (GET /collections/{collectionId}/instances)
func (h *Handler) GetCollectionInstances(ctx echo.Context, collectionId CollectionId, params GetCollectionInstancesParams) error {
	var collections []Collection
	for _, collectionName := range []string{"MEPS"} {
		collection, err := h.getQueries(collectionName)
		collection.Id = "2024-01-01T03:00:00Z"
		if err != nil {
			return writer(params.F, ctx)(http.StatusInternalServerError,
				Exception{
					Code:        "internal error",
					Description: ptr("Internal server error"),
				},
			)
		}
		collections = append(collections, *collection)
	}

	ret := &Instances{
		Links: []Link{
			{
				Href: h.baseURL + "/edr/collections/" + collectionId + "/",
				Rel:  "self",
				Type: ptr("application/json"),
			},
		},
		Instances: collections,
	}

	return writer(params.F, ctx)(http.StatusOK, &ret)
}

// Query end point for position queries  of collection {collectionId}
// (GET /collections/{collectionId}/position)
func (h *Handler) GetDataForPoint(ctx echo.Context, collectionId CollectionId, params GetDataForPointParams) error {
	// point, err := wkt.ParsePoint(params.Coords)
	// if err != nil {
	// 	return writer(params.F, ctx)(http.StatusBadRequest,
	// 		Exception{
	// 			Code:        "bad request",
	// 			Description: ptr(err.Error()),
	// 		},
	// 	)
	// }
	// if params.Z != nil {
	// 	z, err := strconv.ParseFloat(*params.Z, 64)
	// 	if err != nil {
	// 		return writer(params.F, ctx)(http.StatusBadRequest,
	// 			Exception{
	// 				Code:        "bad request",
	// 				Description: ptr("invalid value for Z"),
	// 			},
	// 		)
	// 	}
	// 	point.Z = &z
	// }

	// coords := []float32{float32(point.X), float32(point.Y)}
	// if point.Z != nil {
	// 	coords = append(coords, float32(*point.Z))
	// }
	// ret := PointGeoJSON{
	// 	Type:        Point,
	// 	Coordinates: []float32{10, 60},
	// }

	return writer(params.F, ctx)(http.StatusOK, covJsonForPoint())
}

// Query end point for position queries of instance {instanceId} of collection {collectionId}
// (GET /collections/{collectionId}/instances/{instanceId}/position)
func (h *Handler) GetInstanceDataForPoint(ctx echo.Context, collectionId CollectionId, instanceId InstanceId, params GetInstanceDataForPointParams) error {
	return writer(params.F, ctx)(http.StatusOK, covJsonForPoint())
}

// Information about standards that this API conforms to
// (GET /conformance)
func (h *Handler) GetRequirementsClasses(ctx echo.Context, params GetRequirementsClassesParams) error {
	ret := ConfClasses{
		ConformsTo: []string{
			"http://www.opengis.net/spec/ogcapi-edr-1/1.0/conf/core",
			"http://www.opengis.net/spec/ogcapi-common-1/1.0/conf/core",
			"http://www.opengis.net/spec/ogcapi-common-2/1.0/conf/collections",
			"http://www.opengis.net/spec/ogcapi-edr-1/1.0/conf/oas30",
			"http://www.opengis.net/spec/ogcapi-edr-1/1.0/conf/html",
			"http://www.opengis.net/spec/ogcapi-edr-1/1.0/conf/geojson",
		},
	}
	return writer(params.F, ctx)(http.StatusOK, &ret)
}

func writer(f *F, ctx echo.Context) func(code int, i interface{}) error {
	if f == nil || *f == "json" {
		return ctx.JSON
	}
	// Fallback is json
	return ctx.JSON
}

func ptr[T any](t T) *T {
	return &t
}

func covJsonForPoint() *CoverageJSON {

	observedPropertyID := "http://vocab.nerc.ac.uk/standard_name/air_temperature/"

	parameters := map[string]Parameter{
		"air_temperature": {
			Type: "Parameter",
			Description: &I18n{
				"en": "air_temperature, cf-convention?",
			},
			Unit: &Unit{
				Label: &I18n{
					"en": "Kelvin",
				},
			},
			ObservedProperty: ObservedProperty{
				Id: &observedPropertyID,
				Label: I18n{
					"en": "Air temperature is the bulk temperature of the air, not the surface (skin) temperature.",
				},
			},
		},
	}

	temporalRefSys := ReferenceSystem{
		Type: "TemporalRS",
	}
	temporalRefSys.FromReferenceSystem0(ReferenceSystem0{
		Calendar: "Gregorian",
	})

	spatialRefSys := ReferenceSystem{
		Type: "GeographicCRS",
	}
	spatialRefSysID := "http://www.opengis.net/def/crs/OGC/1.3/CRS84"
	spatialRefSys.FromReferenceSystem1(ReferenceSystem1{
		Id: &spatialRefSysID,
	})

	domainType := DomainDomainType("PointSeries")
	domain := Domain{
		DomainType: &domainType,
		Type:       "Domain",
		Axes: struct {
			// T Simple axis with string values (e.g. time strings)
			T *StringValuesAxis `json:"t,omitempty"`

			// X Simple axis with numeric values
			X MettsnumericValuesAxis `json:"x"`

			// Y Simple axis with numeric values
			Y MettsnumericValuesAxis `json:"y"`

			// Z Simple axis with numeric values
			Z *MettsnumericValuesAxis `json:"z,omitempty"`
		}{
			X: MettsnumericValuesAxis{Values: &[]float32{11.0}},
			Y: MettsnumericValuesAxis{Values: &[]float32{60.0}},
			T: &StringValuesAxis{
				Values: []string{"2024-01-01T03:00:00Z", "2024-01-01T06:00:00Z"},
			},
			Z: &MettsnumericValuesAxis{Values: &[]float32{100000, 50000}},
		},
		Referencing: &[]ReferenceSystemConnection{
			{
				Coordinates: []string{"x", "y"},
				System:      spatialRefSys,
			},
			{
				Coordinates: []string{"t"},
				System:      temporalRefSys,
			},
		},
	}

	ranges := map[string]NdArray{
		"air_temperature": {
			Type:      "NdArray",
			DataType:  "float",
			AxisNames: &[]string{"z", "t"},
			Shape:     &[]float32{2, 3},
			Values:    []float32{-20.8, -20.1, -19.5, -25.8, -25.1, -24.5},
		},
	}

	return &CoverageJSON{
		Type:       "Coverage",
		Parameters: &parameters,
		Domain:     &domain,
		Ranges:     &ranges,
	}
}
