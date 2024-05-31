# Override Requests

**InfiniteMITM** enables you to intercept and modify the game's requests and responses on the fly. To customize them, you can edit the `mitm.yaml` file in the root of the generated folder located in your home directory (e.g., `C:\Users\<username>\InfiniteMITM`). The `mitm.yaml` file uses a specific configuration that lets you match various paths based on a service (`blobs` | `authoring` | `discovery` | `stats` | `settings` | `gamecms` | `economy` | `root`) and desired REST methods (`GET` | `POST` | `PATCH` | `PUT` | `DELETE`).

**Note:** When changing the request `body`, the `Content-Length` header will be automatically calculated.

## Example

```yaml
blobs: # blobs-infiniteugc.svc.halowaypoint.com
    - path: /ugcstorage/map/:guid/:guid/:map-mvar # Path pattern to match, will catch all .mvar files
      methods: # HTTP methods that this configuration will handle
          - GET
      response:
          body: :mitm-dir/maps/design_21.mvar # Path to the file that will be used as the response body
          headers: # Additional headers to include in the response
              x-infinite-mitm: :infinite-mitm-version
              content-type: :ct-bond
    - path: /ugcstorage/enginegamevariant/:guid/:guid/customgamesuimarkup/Slayer_CustomGamesUIMarkup_en.bin # Path pattern for specific "CustomGamesUIMarkup", for any assetID and assetVersionID
      methods:
          - GET
      response:
          body: :mitm-dir/enginegamevariant/cgui-markups/Slayer_8Teams.bin
          headers:
              x-infinite-mitm: :infinite-mitm-version
              content-type: :ct-bond
    - path: /ugcstorage/enginegamevariant/:guid/:guid/FFA.bin # Path pattern for specific "EngineGameVariant", for any assetID and assetVersionID
      methods:
          - GET
      response:
          body: :blobs-svc/enginegamevariant/$1/9b0d3fd4-2027-4dca-96f5-899b449408e2/FFA.bin # Path to the external file that will be used as the response body, with a specific assetVersionID
          headers:
              x-infinite-mitm: :infinite-mitm-version
              content-type: :ct-bond
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
root: # Must be one of "blobs | authoring | discovery | settings | root" (root = all)
    # Each item in the list represents a specific endpoint configuration.
    - path: /example/path # Targeted path (case insensitive)
      methods: # List of HTTP methods to catch (GET, POST, PATCH, PUT, DELETE)
          - GET
          - POST
      request: # Used to alter the request
          body: :mitm-dir/request/body/file # URI to the file submitted for PUT, POST, and PATCH requests instead of the initial payload
          headers: # Override request headers (case insensitive)
              custom-header: customValue
      response: # Used to alter the response
          body: :mitm-dir/response/body/file # URI to the overridden file
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
-   `:ct-bond`
    -   Represents the content type of binary files consumed by the game.
    -   Output: `application/x-bond-compact-binary`
-   `:ct-json`
    -   Represents a JSON content type.
    -   Output: `application/json
-   `:ct-xml`
    -   Represents a XML content type.
    -   Output: `application/xml
-   `:*`
    -   Will match everything else.
    -   Example: `/foo/bar:*` will match `/foo/bar/baz`
-   `:$`
    -   Ends the match expression.
    -   Example: `/foo/bar:$` will not match `/foo/bar/baz`
-   `:blobs-svc`
    -   Returns blobs service URL.
    -   Output: `https://blobs-infiniteugc.svc.halowaypoint.com`
-   `:authoring-svc`
    -   Returns authoring service URL.
    -   Output: `https://authoring-infiniteugc.svc.halowaypoint.com`
-   `:discovery-svc`
    -   Returns discovery service URL.
    -   Output: `https://discovery-infiniteugc.svc.halowaypoint.com`
-   `:stats-svc`
    -   Returns stats service URL.
    -   Output: `https://halostats.svc.halowaypoint.com`
-   `:settings-svc`
    -   Returns stats service URL.
    -   Output: `https://settings.svc.halowaypoint.com`
-   `:gamecms-svc`
    -   Returns gamecms service URL.
    -   Output: `https://gamecms-hacs.svc.halowaypoint.com`
-   `:economy-svc`
    -   Returns economy service URL.
    -   Output: `https://economy.svc.halowaypoint.com`
-   `:mitm-dir`
    -   Represents the root folder of local files (only suitable for `response.body`).
    -   Output: `~/InfiniteMITM`
-   `:mitm-version`
    -   Represents the InfiniteMITM version (only suitable for `response.headers`).
    -   Example: `0.1.0`

## Response Match Parameters

In some cases, you might need to reuse a parameter that was matched during the request in your response. To do this, you can use `${pos}` where `{pos}` is the position of your route parameter.

### Example

#### Request Path

```
/ekur/97fd2ab9-ece0-41c1-91a8-f0382f24e6d2/olympus/xuid(1234)
```

#### MITM Config

```yaml
blobs:
    - path: /ekur/:guid/olympus/:xuid
      response:
          body: :mitm-dir/example/xuid($2)/test_$1
```

#### Output

```
~/InfiniteMITM/xuid(1234)/test_97fd2ab9-ece0-41c1-91a8-f0382f24e6d2
```
