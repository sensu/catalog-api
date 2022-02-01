# Sensu Catalog API

![GoReleaser Build](https://github.com/sensu/catalog-api/actions/workflows/goreleaser.yml/badge.svg)

Sensu Catalog API generates a static API from a Sensu Catalog repository (e.g. https://github.com/sensu/catalog). 

## API Endpoints

### `GET /version.json`

Returns the latest content version (used by the Sensu web app to determine the latest API subpath).

#### Example Response

```json
{
  "release_sha256": "af3c54b86b90fac8977f1bdc80d955002dd3f441bdbb4cc603c94abbb929dcf6",
  "last_updated": 1643664852
}
```

### `GET /<release_sha256>/v1/catalog.json`

Returns the list of integration namespaces & names for the catalog.

#### Example Response

```json
{
  "namespaced_integrations": {
    "nginx": [
      "nginx-monitoring"
    ],
    "system": [
      "host-monitoring"
    ]
  }
}

```

### `GET /<release_sha256>/v1/<integration_namespace>/<integration_name>.json`

Returns the integration configuration for the latest version along with a list of
versions for the requested integration.

#### Example Response

```json
{
  "metadata": {
    "name": "nginx-monitoring",
    "namespace": "nginx"
  },
  "class": "community",
  "contributors": [
    "@nixwiz",
    "@calebhailey"
  ],
  "provider": "agent/check",
  "short_description": "NGINX monitoring",
  "supported_platforms": [
    "darwin",
    "linux",
    "windows"
  ],
  "tags": [
    "http",
    "nginx",
    "webserver"
  ],
  "versions": [
    "20220125.0.0",
    "20220126.0.0"
  ]
}
```

### `GET /<release_sha256>/v1/<integration_namespace>/<integration_name>/versions.json`

Returns the list of available versions for the requested integration.

#### Example Response

```json
[
  "20220125.0.0",
  "20220126.0.0"
]
```

### `GET /<release_sha256>/v1/<integration_namespace>/<integration_name>/<version>.json`

Returns the integration configuration for the requested version of an integration.

#### Example Response

```json
{
  "metadata": {
    "name": "nginx-monitoring",
    "namespace": "nginx"
  },
  "class": "community",
  "contributors": [
    "@nixwiz",
    "@calebhailey"
  ],
  "provider": "agent/check",
  "short_description": "NGINX monitoring",
  "supported_platforms": [
    "darwin",
    "linux",
    "windows"
  ],
  "tags": [
    "http",
    "nginx",
    "webserver"
  ],
  "version": "20220125.0.0"
}
```

### `GET /<release_sha256>/v1/<integration_namespace>/<integration_name>/<version>/sensu-resources.json`

Returns the Sensu resources, in JSON format, for the requested integration version.

#### Example Response

```json
{
  "api_version": "core/v2",
  "metadata": {
    "name": "nginx-healthcheck"
  },
  "spec": {
    "command": "check-nginx-status.rb --url {{ .annotations.check_nginx_status_url | default \"http://localhost:80/nginx_status\" }}",
    "interval": 30,
    "pipelines": [
      {
        "api_version": "core/v2",
        "name": "alerts",
        "type": "Pipeline"
      },
      {
        "api_version": "core/v2",
        "name": "incident-management",
        "type": "Pipeline"
      }
    ],
    "publish": true,
    "runtime_assets": [
      "sensu-plugins/sensu-plugins-nginx:3.1.2",
      "sensu/sensu-ruby-runtime:0.0.10"
    ],
    "subscriptions": [
      "nginx"
    ],
    "timeout": 10
  },
  "type": "CheckConfig"
}
```

### `GET /<release_sha256>/v1/<integration_namespace>/<integration_name>/<version>/README.md`

Returns the README, in markdown format, for the requested integration version.

### `GET /<release_sha256>/v1/<integration_namespace>/<integration_name>/<version>/CHANGELOG.md`

Returns the CHANGELOG, in markdown format, for the requested integration version.

### `GET /<release_sha256>/v1/<integration_namespace>/<integration_name>/<version>/logo.png`

Returns the logo, in PNG format, for the requested integration version.
