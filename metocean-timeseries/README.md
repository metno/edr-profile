# EDR profile: Timeseries for meteorological and oceanographic (metocean) data.

A service that is compliant with this profile is also compliant with [OGC-API EDR spec v1.1.0](https://docs.ogc.org/is/19-086r6/19-086r6.html).

Add this profile as a conformance class to your api with the link: [https://github.com/metno/edr-profile/blob/main/metocean-timeseries/README.md](https://github.com/metno/edr-profile/blob/main/metocean-timeseries/README.md)

## Overview

A profile for an EDR service that has metocean timeseries collections. Such a timeseries delivers typically a set for parameters from a forecast model for a number of timesteps for one vertical level or in-situ observations as timeseries. The timeseries is encoded in CoverageJSON using the PointSeries domainType.

The docs for the profile contains the following:

- `Conformance` class with a list of requirements for the profile.
- [OpenAPI spec](openapi/metocean-ts.yaml) compliant with the profile.
- [Examples](examples/) list examples of how responses from a compliant service will be.

OpenAPI specifications in this profile was copied from https://github.com/opengeospatial/ogcapi-environmental-data-retrieval, and then modified to fit the profile.

## Conformance

**Conformance Class**: https://github.com/metno/edr-profile/blob/main/metocean-timeseries/README.md#Conformance 

**Target type**: Web API

**Requirements Class**: https://github.com/metno/edr-profile/blob/main/metocean-timeseries/README.md#Requirements

**Dependency**: http://www.opengis.net/spec/ogcapi-edr-1/1.0/conf/core

**Dependency**: http://www.opengis.net/spec/ogcapi-edr-1/1.0/conf/coveragejson

## Requirements

### Requirement 1

#### Collections

- A: A collection represents a metocean data source (e.g NWP model, observations for a geographic area with some shared characteristics).
- B: temporal extent SHALL either be null or specify the start and end time that cover all instances of this collection.
- C: CRS of the spatial extent and the datasets returned by data queries SHOULD be WGS 84: http://www.opengis.net/def/crs/OGC/1.3/CRS84.
- D: parameter_names includes all parameters mentioned in at least one of the instances. No guarantee that a parameter will be available for all instances.
- E: Use the CF-convention(http://cfconventions.org/Data/cf-standard-names/current/build/cf-standard-name-table.html) to describe parameters when possible. The member id of observedProperty SHOULD be a unique resolvable url. If no suitable term exists in the list CF-convention standard names, use a CF-convention standard name as a starting point for constructing a unique id of your parameter.

### Requirement 2

#### Instances

- A: id of an instance SHOULD be a unique string, e.g an uuid.
- B: The title of the instance SHOULD be a short string meant to easily convey what sets this instance apart from all the other instances of the collection. E.g it should be a size and contant that makes it suitable fir populating a list of options a person can be pick from. For forecast model data it is natural to use the reference time of the produced model run as part of the title.
- C: If it is necessary to have the reference time of the instance (e.g for forecast model runs) as structured data then add another custom top-level property `reference_time` with a string value on the format of `20240101T000000Z`.
- D: All temporal values are on the format: https://datatracker.ietf.org/doc/html/rfc3339#section-5.6.
- E: parameter_names SHALL list the exact parameters available in this instance of the collection.

### Requirement 3

#### Data queries

- A: Only position and location data queries are supported.
- B: Each collection supports a maximum of one type of vertical level(meter, model level, pressure). So  each type of vertical level will be represented by a separate collection.
- C: A collection may be queried directly, in addition to queries through its instances. If this is supported, that query SHOULD be done against the most recently created instance of the collection.

### Requirement 4

#### Response format for data queries

- A: The response to a data query request SHALL be CoverageJSON with domainType PointSeries.
