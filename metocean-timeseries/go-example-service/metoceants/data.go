package metoceants

import (
	"fmt"
	"time"
)

func getTimes() []time.Time {
	var times []time.Time
	for _, ts := range []string{"2024-01-01T03:00:00Z", "2024-01-01T04:00:00Z", "2024-01-01T05:00:00Z"} {
		tmp, _ := time.Parse(time.RFC3339, ts)
		times = append(times, tmp)
	}
	return times
}

func (h *Handler) baseCollection(collectionId CollectionId, collectionPath string) (*Collection, error) {
	parameterNameID := "https://vocab.nerc.ac.uk/standard_name/air_temperature/"
	parameterNameLabel := ObservedPropertyCollection_Label{}
	parameterNameLabel.FromObservedPropertyCollectionLabel0(
		ObservedPropertyCollectionLabel0("Temperature"),
	)
	symbolValue := "K"
	symbolType := "https://qudt.org/vocab/unit/K"

	defaultFormat := "CoverageJSON"
	positionDescription := "Query data by a specified geographical point."
	locationsDescription := "Query list of geographical locations and data for one of those location."

	times := getTimes()

	return &Collection{
		// Links []GetQueriesOKApplicationJSONLinksItem `json:"links"`
		Id:       collectionId,
		Keywords: &[]string{"forecast", "timeseries", "nordic", "air_temperature"},
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
				Interval: [][]time.Time{{times[0], times[len(times)-1]}},
				Values:   &times,
				Trs:      "https://tools.ietf.org/html/rfc3339#section-5.6", // "http://www.opengis.net/def/uom/ISO-8601/0/Gregorian",
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
			Instances *struct {
				Link *InstancesLink "json:\"link,omitempty\""
			} "json:\"instances,omitempty\""
			Items *struct {
				Link *Link "json:\"link,omitempty\""
			} "json:\"items,omitempty\""
			Locations *struct {
				Link *LocationsLink "json:\"link,omitempty\""
			} "json:\"locations,omitempty\""
			Position struct {
				Link *PositionLink "json:\"link,omitempty\""
			} "json:\"position\""
		}{
			Position: struct {
				Link *PositionLink `json:"link,omitempty"`
			}{
				Link: &PositionLink{
					Href:      fmt.Sprintf("%s/%s/position?coords={coords}", h.baseURL, collectionPath),
					Rel:       "data",
					Templated: ptr(true),
					Variables: &PositionDataQuery{
						OutputFormats:       &[]string{"CoverageJSON"},
						DefaultOutputFormat: &defaultFormat,
						CrsDetails: &[]CrsObject{
							{
								Crs: "EPSG:4326",
								Wkt: "GEOGCS[\"WGS 84\",DATUM[\"WGS_1984\",SPHEROID[\"WGS 84\",6378137,298.257223563,AUTHORITY[\"EPSG\",\"7030\"]],AUTHORITY[\"EPSG\",\"6326\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.01745329251994328,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4326\"]]",
							},
						},
						Description: &positionDescription,
						QueryType:   "position",
					},
				},
			},
			Locations: &struct {
				Link *LocationsLink `json:"link,omitempty"`
			}{
				Link: &LocationsLink{
					Href: fmt.Sprintf("%s/%s/locations", h.baseURL, collectionPath),
					Rel:  "data",
					Variables: &LocationsDataQuery{
						OutputFormats:       &[]string{"CoverageJSON"},
						DefaultOutputFormat: &defaultFormat,
						CrsDetails: &[]CrsObject{
							{
								Crs: "EPSG:4326",
								Wkt: "GEOGCS[\"WGS 84\",DATUM[\"WGS_1984\",SPHEROID[\"WGS 84\",6378137,298.257223563,AUTHORITY[\"EPSG\",\"7030\"]],AUTHORITY[\"EPSG\",\"6326\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.01745329251994328,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4326\"]]",
							},
						},
						Description: &locationsDescription,
						QueryType:   "locations",
					},
				},
			},
		},

		Crs:           []string{"CRS84"},
		OutputFormats: []string{"CoverageJSON"},

		ParameterNames: map[string]ParameterNames{
			"air_temperature": {
				Type: "Parameter",
				ObservedProperty: ObservedPropertyCollection{
					Id:    &parameterNameID,
					Label: parameterNameLabel,
				},
				Unit: &CollectionUnit{
					Label: "Kelvin",
					Symbol: CollectionUnitSymbol{
						Value: &symbolValue,
						Type:  &symbolType,
					},
				},
			},
		},
	}, nil
}

func (h *Handler) getInstanceCollection(reftime time.Time, collectionName string) (*Collection, error) {
	instanceID := "d5b6e12e-fd85-447e-b35d-5aa429ad5e5e"
	collectionPath := fmt.Sprintf("collections/%s/instances/%s", collectionName, instanceID)

	collection, _ := h.baseCollection(instanceID, collectionPath)
	collection.Title = ptr(fmt.Sprintf("%s %s", collectionName, reftime.Format(time.RFC3339)))

	return collection, nil
}

func (h *Handler) getCollection(collectionId CollectionId, collectionPath string) (*Collection, error) {
	collection, _ := h.baseCollection(collectionId, collectionPath)
	collection.Title = ptr("MEPS")

	collection.DataQueries.Instances = &struct {
		Link *InstancesLink `json:"link,omitempty"`
	}{
		Link: &InstancesLink{
			Href: fmt.Sprintf("%s/%s/instances", h.baseURL, collectionPath),
			Rel:  "data",
			Variables: &InstancesDataQuery{
				QueryType: "instances",
			},
		},
	}

	collection.Links = append(collection.Links,
		Link{
			Href: fmt.Sprintf("%s/collections/%s/instances", h.baseURL, collectionId),
			Rel:  "collection",
			Type: ptr("application/json"),
		})
	collection.Links = append(collection.Links,
		Link{
			Href: fmt.Sprintf("%s/%s", h.baseURL, collectionPath),
			Rel:  "self",
			Type: ptr("application/geo+json"),
		})

	return collection, nil
}

func covJsonForPoint() *CoverageJSON {
	observedPropertyID := "http://vocab.nerc.ac.uk/standard_name/air_temperature/"

	parameters := map[string]Parameter{
		"air_temperature": {
			Type: "Parameter",
			Description: &I18n{
				"en": "Air temperature is the bulk temperature of the air, not the surface (skin) temperature.",
			},
			Unit: &GeoJSONunit{
				Label: I18n{
					"en": "Kelvin",
				},
				Symbol: &struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				}{
					Type:  "https://qudt.org/vocab/unit/K",
					Value: "K",
				},
			},
			ObservedProperty: ObservedProperty{
				Id: &observedPropertyID,
				Label: I18n{
					"en": "air_temperature",
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
		Description: map[string]string{
			"en": "WGS84 geographical coordinate system using longitude,latitude as values.",
		},
	})

	verticalRefSys := ReferenceSystem{
		Type: "VerticalCRS",
	}
	verticalRefSys.FromReferenceSystem1(ReferenceSystem1{
		Description: map[string]string{
			"en": "Vertical coordinate system using pressure(Pa) as values.",
		},
	})

	var timeStamps []string
	for _, t := range getTimes() {
		timeStamps = append(timeStamps, t.Format(time.RFC3339))
	}
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
				Values: timeStamps,
			},
			Z: &MettsnumericValuesAxis{Values: &[]float32{100000}},
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
			{
				Coordinates: []string{"z"},
				System:      verticalRefSys,
			},
		},
	}

	ranges := map[string]NdArray{
		"air_temperature": {
			Type:      "NdArray",
			DataType:  "float",
			AxisNames: []string{"t"},
			Shape:     []float32{3},
			Values:    []float32{-20.8, -20.1, -19.5},
		},
	}

	return &CoverageJSON{
		Type:       "Coverage",
		Parameters: &parameters,
		Domain:     &domain,
		Ranges:     &ranges,
	}
}

func getLocations(basePath string) *EdrFeatureCollectionGeoJSON {
	locationId := "oslo"
	return &EdrFeatureCollectionGeoJSON{
		Links: &[]Link{},
		Features: []FeatureGeoJSON{
			{
				Type: "Feature",
				Geometry: PointGeoJSON{
					Type:        "Point",
					Coordinates: []float32{11, 60},
				},
				Properties: EdrPropertiesForecastTS{
					Name: &I18n{
						"en": "Oslo",
					},
					Edrqueryendpoint: fmt.Sprintf("%s/locations/oslo", basePath),
					ParameterName:    &[]string{"air_temperature"},
				},
				Id: &locationId,
			},
		},
		Type: "FeatureCollection",
	}
}
