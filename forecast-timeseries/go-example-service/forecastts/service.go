package forecastts

import (
	"fmt"
	"net/http"

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
		Description: ptr("A sample edr service, implemented using ogen."),
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
	for _, collectionName := range []string{"oceanforecast"} {
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
	const id = "oceanforecast"

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
		Title:    ptr("Ocean forecast"),
		Keywords: &[]string{"forecast", "ocean"},
		Links: []Link{
			{
				Href: fmt.Sprintf("%s/collections/%s", h.baseURL, collectionId),
				Rel:  "data",
				Type: ptr("application/geo+json"),
			},
			{
				Href: fmt.Sprintf("%s/collections/%s/position", h.baseURL, collectionId),
				Rel:  "data",
				Type: ptr("application/geo+json"),
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
		OutputFormats: []string{"GeoJSON"},

		ParameterNames: map[string]interface{}{
			"sea_surface_wave_height": []byte("{}"),
		},
	}, nil
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
	ret := PointGeoJSON{
		Type:        Point,
		Coordinates: []float32{10, 60},
	}

	return writer(params.F, ctx)(http.StatusOK, &ret)
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
