# cueblox

`cueblox` is a tool for creating schemas written in [CUE](https://cuelang.org)

## Server Operations

`cueblox` provides tools to assist in preparation of a schema repository.

### Conventions

Each schema repository will have a manifest file with metadata about the schemas available. It will also contain one or more schema definitions in Cue format.

The manifest file is stored as JSON encoded data, and can be served as a static file. Schemas are stored as UTF-8 encoded text, and are served as static files. Therefore any static hosting service can serve as a schema repository server.

Schemas have a `namespace` and a `name`. The `name` is a short label for the schema, and the `namespace` is the canonical identifier for the schema. By convention, the `namespace` is a URI which includes any version information.

Example:

```yaml
namespace: schemas.cueblox.com/v1
name: devrel
```
