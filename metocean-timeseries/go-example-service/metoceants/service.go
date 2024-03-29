package metoceants

import (
	"encoding/json"
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
		Title:       ptr("Example edr service"),
		Description: ptr("An example edr service compliant with metocean timeseries profile. Service is based on autogenerated go-code from openapi."),
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
			"https://pages.github.com/metno/edr-profile/forecast-timeseries-0.1",
		},
	}
	return writer(params.F, ctx)(http.StatusOK, &ret)
}

// List the available collections from the service
// (GET /collections)
func (h *Handler) ListCollections(ctx echo.Context, params ListCollectionsParams) error {
	var collections []Collection
	for _, collectionName := range []string{"MEPS"} {
		collectionPath := fmt.Sprintf("collections/%s", collectionName)
		collection, err := h.getCollection(collectionName, collectionPath)
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

	collectionPath := fmt.Sprintf("collections/%s", id)
	ret, err := h.getCollection(id, collectionPath)
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
	refTime := getTimes()[0]
	for _, collectionName := range []string{"MEPS"} {
		collection, err := h.getInstanceCollection(refTime, collectionName)
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
				Href: h.baseURL + "/collections/" + collectionId + "/instances/",
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
	return ctx.Blob(http.StatusOK, "application/vnd.cov+json", payload)
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
	return ctx.Blob(http.StatusOK, "application/vnd.cov+json", payload)
	//return writer(params.F, ctx)(http.StatusOK, covJsonForPoint())
}

// List available location identifers for the collection
// (GET /collections/{collectionId}/locations)
func (h *Handler) ListCollectionDataLocations(ctx echo.Context, collectionId CollectionId, params ListCollectionDataLocationsParams) error {
	locationsData := getLocations(fmt.Sprintf("%s/collections/%s", h.baseURL, collectionId))
	payload, err := json.Marshal(locationsData)
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, "application/geo+json", payload)
}

// List available location identifers for the instance
// (GET /collections/{collectionId}/instances/{instanceId}/locations)
func (h *Handler) ListDataInstanceLocations(ctx echo.Context, collectionId CollectionId, instanceId InstanceId, params ListDataInstanceLocationsParams) error {
	locationsData := getLocations(fmt.Sprintf("%s/collections/%s/instances/%s",
		h.baseURL, collectionId, instanceId))
	payload, err := json.Marshal(locationsData)
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, "application/geo+json", payload)
}

// Query end point for queries of collection {collectionId} defined by a location id
// (GET /collections/{collectionId}/locations/{locationId})
func (h *Handler) GetCollectionDataForLocation(ctx echo.Context, collectionId CollectionId, locationId LocationId, params GetCollectionDataForLocationParams) error {
	pointData := covJsonForPoint()
	payload, err := json.Marshal(pointData)
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, "application/vnd.cov+json", payload)
}

// Query end point for queries of instance {instanceId} of collection {collectionId} defined by a location id
// (GET /collections/{collectionId}/instances/{instanceId}/locations/{locationId})
func (h *Handler) GetInstanceDataForLocation(ctx echo.Context, collectionId CollectionId, instanceId InstanceId, locationId LocationId, params GetInstanceDataForLocationParams) error {
	pointData := covJsonForPoint()
	payload, err := json.Marshal(pointData)
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, "application/vnd.cov+json", payload)
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
