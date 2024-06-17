# Override Requests

**InfiniteMITM** allows you to intercept and modify the game's requests and responses on the fly. To customize them, you can edit the `mitm.yaml` file in the **InfiniteMITM** directory within your home directory (e.g., `C:\Users\<username>\InfiniteMITM`). The `mitm.yaml` file uses a specific configuration that lets you match various paths based on a service (`blobs`, `authoring`, `discovery`, `stats`, `settings`, `gamecms`, `economy`, `lobby`, `skill`, `root`), where `root` is a catch-all, desired REST methods (`GET`, `POST`, `PATCH`, `PUT`, `DELETE`), and **regex** support.

### Notes:

-   All changes are applied upon saving, so there is no need to restart **InfiniteMITM**.
-   When changing the request `body`, the `Content-Length` header will be automatically recalculated.
-   By default, only the overridden traffic will be displayed. This behavior can be changed in the `mitm.yaml` file.
    -   Displaying `all` requests and responses may impact performance.
-   Make sure not to send sensitive information (e.g., `X-343-Authorization-Spartan`) when altering the request `body`.
    - Example: https://github.com/Alexis-Bize/InfiniteMITM/blob/main/examples/surasia/mitm.yaml#L9

## Example

```yaml
domains:
  blobs: # blobs-infiniteugc.svc.halowaypoint.com
    - path: "/ugcstorage/map/:guid/:guid/:map-mvar" # Path pattern to match, will catch all .mvar files
      methods: # HTTP methods that this configuration will handle
        - GET
      response: # Response handler
        body: ":mitm-dir/resources/ugc/maps/design_21.mvar" # Path to the file that will be used as the response body
        headers: # Additional headers to include in the response
          x-infinite-mitm-version: ":mitm-version"
          content-type: ":ct-bond"
    - path: "/ugcstorage/enginegamevariant/:guid/:guid/customgamesuimarkup/Slayer_CustomGamesUIMarkup_en.bin" # Path pattern for specific "CustomGamesUIMarkup", for any assetID and assetVersionID
      methods:
        - GET
      response:
        body: ":mitm-dir/resources/ugc/enginegamevariants/cgui-markups/Slayer_8Teams.bin"
        headers:
          x-infinite-mitm-version: ":mitm-version"
          content-type: ":ct-bond"
    - path: "/ugcstorage/enginegamevariant/:guid/:guid/FFA.bin" # Path pattern for specific "EngineGameVariant", for any assetID and assetVersionID
      methods:
        - GET
      response:
        body: ":blobs-svc/enginegamevariant/$1/9b0d3fd4-2027-4dca-96f5-899b449408e2/FFA.bin" # Path to the external file that will be used as the response body, with a specific assetVersionID
        headers:
          x-infinite-mitm-version: ":mitm-version"
          content-type: ":ct-bond"
    - path: "/ugcstorage/enginegamevariant/:guid/:guid/:egv-bin" # Match any "EngineGameVariant"
      methods:
        - GET
      response:
        before: # Response pre-handler
          commands: # Commands list
            - run: # Run the first command
              - "echo \"First GUID: :guid\""
            - run: # Run the second command
              - "echo \"Second GUID: :guid\""
    - path: "/ugcstorage/*" # Match all after /ugcstorage/
      methods:
        - GET
        - POST
        - PATCH
        - PUT
        - DELETE
      request: # Request handler
        headers: # Headers to override in the request
          x-343-authorization-spartan: "v4=MyCustomSpartanToken"
```

## Definition

```yaml
domains:
  root: # Must be one of "blobs | authoring | discovery | settings | root" (root = all)
    # Each item in the list represents a specific endpoint configuration.
    - path: "/example/path" # Targeted path (case insensitive)
      methods: # List of HTTP methods to catch (GET, POST, PATCH, PUT, DELETE)
        - GET
        - POST
      request: # Used to alter the request
        before: # Used to run various actions before handler execution
          commands: # Used to run desired commands
            - run: # Used to define a run command
              # Will be concatenated into: echo "hello" && echo "world"
              - "echo \"hello\"" # Parameter
              - "&& echo \"world\"" # Parameter
    - path: "/ugcstorage/*" # Match all after /ugcstorage/
        body: ":mitm-dir/request/body/file" # URI to the file submitted for PUT, POST, and PATCH requests instead of the initial payload
        headers: # Override request headers (case insensitive)
          custom-header: "customValue"
      response: # Used to alter the response
        before: # Used to run a command before handler execution
          cmd: "shell command" # Desired command
        code: 200 # Status code (optional), see https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
        body: ":mitm-dir/response/body/file" # URI to the overridden file
        headers: # Override response headers (case insensitive)
          custom-response-header: "customValue"
```

### Before Command

Please refer to our [Command](/docs/Command.md) documentation for further details.

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
-   `:sandbox`
    -   Matches all known sandboxes.
    -   Sandboxes: `retail` | `test` | `beta` | `beta-test`
-   `:title`
    -   Matches all known titles.
    -   Titles: `hi` | `hipnk` | `higrn` | `hired` | `hipur` | `hiorg` | `hiblu` | `hi343`
-   `:ct-bond`
    -   Represents the content type of binary files consumed by the game.
    -   Output: `application/x-bond-compact-binary`
-   `:ct-json`
    -   Represents a JSON content type.
    -   Output: `application/json
-   `:ct-xml`
    -   Represents a XML content type.
    -   Output: `application/xml
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
    -   Output: `<drive>:\Users\<username>\InfiniteMITM`
-   `:mitm-version`
    -   Represents the InfiniteMITM version (only suitable for `response.headers`).
    -   Example: `0.1.0`

## Response Match Parameters

In some cases, you might need to reuse a parameter that was matched during the request in your response. To do this, you can use `${pos}` where `{pos}` is the position of your route parameter or regex.

### Example

#### Request Path

```
/ekur/97fd2ab9-ece0-41c1-91a8-f0382f24e6d2/olympus/xuid(1234)/details
```

#### MITM Config

```yaml
domains:
  blobs:
    - path: "/ekur/:guid/olympus/:xuid/([a-z]+)$"
      response:
        body: ":mitm-dir/example/xuid_$2/test_$1/$3"
```

#### Output

```
~/InfiniteMITM/example/xuid_1234/test_97fd2ab9-ece0-41c1-91a8-f0382f24e6d2/details
```
