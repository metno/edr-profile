# EDR profile: Weather forecast timeseries

A service that is compliant with this profile is also compliant with [OGC-API EDR spec v1.1.0](https://docs.ogc.org/is/19-086r6/19-086r6.html).

Add this profile as a conformance class to your api with the link: [https://github.com/metno/edr-profile/blob/main/profile_weather_forecast_timeseries.md](https://github.com/metno/edr-profile/blob/main/profile_weather_forecast_timeseries.md)

## Overview

A profile for an EDR service that has weather forecast timeseries collection. A Weather forecast timeseries delivers a set for parameters from a weather forecast models for a number of timesteps for one vertical level. The forecast timeseries is encoded in CoverageJSON using the PointSeries domainType.

The docs for the profile contains the following:

- `Conformance` class with a list of requirements for the profile.
- [OpenAPI spec](openapi/forecast-ts.yaml) compliant with the profile.
- A [golang example service](go-example-service) compliant with the profile.

OpenAPI specifications in this profile was copied from https://github.com/opengeospatial/ogcapi-environmental-data-retrieval, and then modified to fit the profile.

## Conformance

**Conformance Class**: https://github.com/metno/edr-profile/blob/main/profile_weather_forecast_timeseries.md#Conformance 

**Target type**: Web API

**Requirements Class**: https://github.com/metno/edr-profile/blob/main/profile_weather_forecast_timeseries.md#Requirements

**Dependency**: http://www.opengis.net/spec/ogcapi-edr-1/1.0/conf/core

**Dependency**: http://www.opengis.net/spec/ogcapi-edr-1/1.0/conf/coveragejson

## Requirements

### Requirement A.1

#### /response_format

- A: The response format for a point data query SHALL be a CoverageJSON document with media type  application/prs.coverage+json.
- B: The response format for a 5xx response SHALL be....
- C: DomainType SHOULD be PointSeries.

### Requirement A.2

#### /collections

- A: temporal extent SHALL either be null or specify the start and end time that cover all instances of this collection.
- B: parameter_names includes all parameters mentioned in at least one of the instances. No guarantee that a parameter will be available for all instances.

#### /collections/<collectionid>/instances/<instanceid>/

- A: id SHALL represent the reference time of the forecast model, the value of the id parameter MUST be on the ISO-8601 format.
- B: CRS must be WGS 84: http://www.opengis.net/def/crs/OGC/1.3/CRS84.
- C: All temporal values are on ISO-8601 format: http://www.opengis.net/def/uom/ISO-8601/0/Gregorian.
- D: temporal extent MUST describe the start and end of the model forecast run that this instance represent.
- E: parameter_names must list the exact parameters available in this instance of the collection.

### Requirement A.3

#### /collections/data_queries

- A: Only position and location data queries are supported.
- B: Each collection supports a maximum of one type of vertical level(meter, model level, pressure). This vertical level is described in the vertical property of the extent in a collection response and in the CoverageJSON data query response as axis z in the domain object.

### Requirement B.1

#### /collections/response

- A: The response to a data query request SHALL be CoverageJSON with domainType PointSeries.
