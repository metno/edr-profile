# edr-profile
NOTE: This repo is no longer updated, and will be removed in the near future.

An edr profile sets constraints and extensions for an OGC-API EDR service providing meteorological and oceanographic data.

The idea is that EDR services distributing a certain data type should be interoperable.
Interoperable in the sense that one client/library should be able to use EDR-services from different projects/institutes to query data for a specific data type without having custom code/plugins for each service.

Each profile is meant to have a precise OpenAPI specification, which means that all collections available from a given landing page will have to be compliant with the profile and its OpenAPI specification.

We build a separate profile for each data type. Each profile is versioned separately using semantic versioning.

The data types are:

- [Timeseries for meteorological and oceanographic data](metocean-timeseries/README.md)
- Aviation
- Radar

## Tentative process for developing the profiles

### Adding a rule to a profile

- Create an issue in this repo.
  - Describe the problem. Specify if its a contraint or an extension.
  - Add label(s) for the data type or use `common` if it applies to all data types.
  - If possible, suggest a solution to the problem.
  - Describe wether this rule is backwards compatible or not.
- Everyone is welcome to discuss and weigh in with suggestions. At least two maintainers from two different insitutions must agree in order for a decision to be reached.

If a project can not wait for a decision, or if it strongly disagrees with the decision beeing made, they can give feedback on this, and proceed on their own. E.g, its better to have an open discussion and accept that in certain areas divergence might be unavoidable.

### Documenting the profiles

- Each data type profile has its own markdown file in this repo.
- Each profile is separated into the sections:
  - Landing page
  - Collection
  - Instances
  - Data queries
  - Response format
- A new file is created when a new major version is necessary.
