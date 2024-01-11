package forecastts

import (
	"encoding/json"
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
				Href: h.baseURL + "/collections/",
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
				Href: h.baseURL + "/collections/" + collectionId + "/",
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
	pointData := covJsonForPoint()
	payload, err := json.Marshal(pointData)
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, "application/prs.coverage+json", payload)
	//return writer(params.F, ctx)(http.StatusOK, covJsonForPoint())
}

// Query end point for position queries of instance {instanceId} of collection {collectionId}
// (GET /collections/{collectionId}/instances/{instanceId}/position)
func (h *Handler) GetInstanceDataForPoint(ctx echo.Context, collectionId CollectionId, instanceId InstanceId, params GetInstanceDataForPointParams) error {
	pointData := covJsonForPoint()
	payload, err := json.Marshal(pointData)
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, "application/prs.coverage+json", payload)
	//return writer(params.F, ctx)(http.StatusOK, covJsonForPoint())
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
