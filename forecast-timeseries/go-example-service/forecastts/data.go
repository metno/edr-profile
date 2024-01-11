package forecastts

import (
	"fmt"
	"time"
)

var referenceTime, _ = time.Parse(time.RFC3339, "2024-01-01T03:00:00Z")
var startTime, _ = time.Parse(time.RFC3339, "2024-01-01T03:00:00Z")
var endTime, _ = time.Parse(time.RFC3339, "2024-01-01T06:00:00Z")

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

	verticalRefSysID := "????"
	verticalRefSys := ReferenceSystem{
		Type: "VerticalCRS",
	}
	verticalRefSys.FromReferenceSystem1(ReferenceSystem1{
		Id: &verticalRefSysID,
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
			AxisNames: &[]string{"t"},
			Shape:     &[]float32{3},
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
