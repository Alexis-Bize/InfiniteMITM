# Override Requests

**InfiniteMITM** enables you to intercept and modify the game's requests and responses on the fly. To customize them, you can edit the `mitm.yaml` file in the root of the generated folder located in your home directory (e.g., `C:\Users\<username>\InfiniteMITM`). The `mitm.yaml` file uses a specific configuration that lets you match various paths based on a service (`blobs` | `authoring` | `discovery` | `settings`) and desired REST methods (`GET` | `POST` | `PATCH` | `PUT` | `DELETE`).

**Note:** When changing the `body`, the `Content-Length` header will be automatically calculated.

## Example

```yaml
blobs: # blobs-infiniteugc.svc.halowaypoint.com
    - path: /ugcstorage/map/:guid/:guid/:map-mvar # Path pattern to match, will catch all .mvar files
      methods: # HTTP methods that this configuration will handle
          - GET
      response:
          body: :infinite-mitm-root/maps/design_21.mvar # Path to the file that will be used as the response body
          headers: # Additional headers to include in the response
              x-infinite-mitm: :infinite-mitm-version
              content-type: :content-type-bond
    - path: /ugcstorage/enginegamevariant/:guid/:guid/customgamesuimarkup/Slayer_CustomGamesUIMarkup_en.bin # Path pattern for specific "CustomGamesUIMarkup", for any assetID and assetVersionID
      methods:
          - GET
      response:
          body: :infinite-mitm-root/enginegamevariant/cgui-markups/Slayer_8Teams.bin
          headers:
              x-infinite-mitm: :infinite-mitm-version
              content-type: :content-type-bond
    - path: /ugcstorage/enginegamevariant/:guid/:guid/FFA.bin # Path pattern for specific "EngineGameVariant", for any assetID and assetVersionID
      methods:
          - GET
      response:
          body: https://blobs-infiniteugc.svc.halowaypoint.com/enginegamevariant/:guid/9b0d3fd4-2027-4dca-96f5-899b449408e2/FFA.bin # Path to the external file that will be used as the response body, with a specific assetVersionID
          headers:
              x-infinite-mitm: :infinite-mitm-version
              content-type: :content-type-bond
    - path: /ugcstorage/:path* # Match all paths after /ugcstorage/
      methods:
          - GET
          - POST
          - PATCH
          - PUT
          - DELETE
      request: # Headers to override in the request
          headers:
              x-343-authorization-spartan: v4=MyCustomSpartanToken
```

## Definition

```yaml
discovery: # Must be one of blobs | authoring | discovery | settings
    # Each item in the list represents a specific endpoint configuration.
    - path: /example/path # Targeted path (case insensitive)
      methods: # List of HTTP methods to catch (GET, POST, PATCH, PUT, DELETE)
          - GET
          - POST
      request: # Used to alter the request
          body: :infinite-mitm-root/request/body/file # URI to the file submitted for PUT, POST, and PATCH requests instead of the initial payload
          headers: # Override request headers (case insensitive)
              custom-header: customValue
      response: # Used to alter the response
          body: :infinite-mitm-root/response/body/file # URI to the overridden file
          headers: # Override response headers (case insensitive)
              custom-response-header: customValue
```

## Predefined Route Parameters

-   `:guid`
    -   Matches a valid GUID.
    -   Example: `1104ee8f-90d9-409f-a295-c9cd3ce16b40`
-   `:map-mvar`
    -   Matches all `Map` files.
    -   Example: `ctf_bazaar.mvar`
-   `:egv-bin`
    -   Matches all `EngineGameVariant` files.
    -   Example: `FFA.bin`
-   `:cgui-bin`
    -   Matches all `CustomGamesUIMarkup` files.
    -   Example: `Slayer_CustomGamesUIMarkup_en.bin`
-   `:spartan-token`
    -   Represents the value of the `X-343-Authorization-Spartan` header used in requests or responses.
    -   Example: `v4=JSONWebToken`
-   `:flight-id`
    -   Represents the value of the `343-Clearance` header used in requests or responses.
    -   Example: `970df3e1-86ae-4571-8488-6b453876da88`
-   `:telemetry-id`
    -   Represents the value of the `343-Telemetry-Session-Id` header used in requests or responses.
    -   Example: `ab66f54f-1c18-47e2-9486-0070e0ae9cc5`
-   `:content-type-bond`
    -   Represents the content type of binary files consumed by the game.
    -   Output: `application/x-bond-compact-binary`
-   `:path*`
    -   Matches all paths.
-   `:infinite-mitm-root`
    -   Represents the root folder of local files.
    -   Output: `~/InfiniteMITM`
-   `:infinite-mitm-version`
    -   Represents the InfiniteMITM version.
    -   Example: `0.1.0`
