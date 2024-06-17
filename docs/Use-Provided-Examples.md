# Use Provided Examples

**InfiniteMITM** offers various [examples](/examples) of rewriting incoming and outgoing traffic. To make your task easier, this documentation will explain how to incorporate these examples into your own `mitm.yaml` file.

## Step-by-Step Guide

### Step 1:

-   Using a file editor, open your `mitm.yaml` file.
    -   This file is located in the **InfiniteMITM** directory within your home directory (e.g., `C:\Users\<username>\InfiniteMITM`).
    -   The following content represents the default template for this file.

```yaml
# do not edit
version: 1

domains:
  # *.svc.halowaypoint.com
  root:
  # authoring-infiniteugc.svc.halowaypoint.com
  authoring:
  # blobs-infiniteugc.svc.halowaypoint.com
  blobs:
  # discovery-infiniteugc.svc.halowaypoint.com
  discovery:
  # economy.svc.halowaypoint.com
  economy:
  # gamecms-hacs.svc.halowaypoint.com
  gamecms:
  # lobby-hi.svc.halowaypoint.com
  lobby:
  # settings.svc.halowaypoint.com
  settings:
  # halostats.svc.halowaypoint.com
  stats:
  # skill.svc.halowaypoint.com
  skill:

options:
  # cache static content in memory or on the disk to minimize network usage and enhance game's performance
  smart_cache:
    enabled: false
    strategy: persistent
    ### ├── memory:     will write cached responses in memory
    ### └── persistent: will write cached responses on the disk (~/InfiniteMITM/cache)
  traffic_display: overrides
  ## ├── all:           will show all requests/responses in the network table
  ## ├── overrides:     will only display overridden requests/responses in the network table
  ## ├── smart_cached:  will only display smart cached requests/responses in the network table
  ## └── silent:        will silent (hide) all requests/responses
```

### Step 2:

-   Open the desired `mitm.yaml` file in one of our provided [examples](/examples).
-   Copy the configuration (specified here between `copy ↓↓↓` and `↑↑↑ copy`).

```yaml
domains:
  discovery:
    # copy ↓↓↓
    - path: /:title/films/matches/:guid/spectate
      methods:
        - GET
      response:
        code: 200
        body: :discovery-svc/$1/films/matches/e04e566e-834f-452a-8764-6fea1cd9dfa3/spectate
        headers:
          content-type: :ct-bond
    # ↑↑↑ copy
```

### Step 3:

-   Paste the copied content under the same example node as follows (`domains` → `discovery`).
    -   Make sure to respect the indentation (spaces).

```yaml
# do not edit
version: 1

domains:
  # *.svc.halowaypoint.com
  root:
  # authoring-infiniteugc.svc.halowaypoint.com
  authoring:
  # blobs-infiniteugc.svc.halowaypoint.com
  blobs:
  # discovery-infiniteugc.svc.halowaypoint.com
  discovery:
    # pasted ↓↓↓
    - path: /:title/films/matches/:guid/spectate
      methods:
        - GET
      response:
        code: 200
        body: :discovery-svc/$1/films/matches/e04e566e-834f-452a-8764-6fea1cd9dfa3/spectate
        headers:
          content-type: :ct-bond
    # ↑↑↑ pasted
  # economy.svc.halowaypoint.com
  economy:
  # gamecms-hacs.svc.halowaypoint.com
  gamecms:
  # lobby-hi.svc.halowaypoint.com
  lobby:
  # settings.svc.halowaypoint.com
  settings:
  # halostats.svc.halowaypoint.com
  stats:
  # skill.svc.halowaypoint.com
  skill:

options:
  # cache static content in memory or on the disk to minimize network usage and enhance game's performance
  smart_cache:
    enabled: false
    strategy: persistent
    ### ├── memory:     will write cached responses in memory
    ### └── persistent: will write cached responses on the disk (~/InfiniteMITM/cache)
  traffic_display: overrides
  ## ├── all:           will show all requests/responses in the network table
  ## ├── overrides:     will only display overridden requests/responses in the network table
  ## ├── smart_cached:  will only display smart cached requests/responses in the network table
  ## └── silent:        will silent (hide) all requests/responses
```

-   Save the file, and you are ready to go.

## Notes

-   All changes are applied upon saving, so there is no need to restart **InfiniteMITM**.
-   Some changes may require to restart **Halo Infinite**.
