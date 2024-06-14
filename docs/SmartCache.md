# SmartCache

**Halo Infinite** tends to repeatedly request the same content (images, binaries, JSON, etc.) as you play, which can be quite frustrating. To reduce network usage and enhance the game's performance, **InfiniteMITM** introduces a solution called **SmartCache**, which, once enabled, will automatically cache this content in memory.

### Cached Services

-   `https://authoring-infiniteugc.svc.halowaypoint.com`
-   `https://discovery-infiniteugc.svc.halowaypoint.com`
-   `https://blobs-infiniteugc.svc.halowaypoint.com`
-   `https://gamecms-hacs.svc.halowaypoint.com`
-   `https://halostats.svc.halowaypoint.com`
-   `https://skill.svc.halowaypoint.com`

**Note:** Some URLs (e.g., your files, favorites, etc.) and overriden requests/responses will not be cached.

## How to Enable It

Open the `mitm.yaml` file located in your home directory (e.g., `C:\Users\<username>\InfiniteMITM`) and enable the `smart_cache` option.

```yaml
options:
  smart_cache:
    enabled: true
    strategy: memory
```

Additionally, you could switch the `traffic_display` option to `silent` or `smart_cached` to reduce network table updates.

## Strategies

- `memory`
    - Will write cached responses in memory, which will be flushed once the CLI is closed.
- `persistent`
    - Will write cached responses to the disk (`~/InfiniteMITM/resources/cache`), making them available after a restart.
