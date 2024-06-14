# Infinite MITM

<p>
    <img alt="InfiniteMITM" title="InfiniteMITM" src="./assets/logo.png" width="256">
    <br>
    <a href="https://github.com/Alexis-Bize/InfiniteMITM/releases"><img src="https://img.shields.io/github/v/release/Alexis-Bize/InfiniteMITM?include_prereleases" alt="Latest Release"></a>
    <a href="https://www.apache.org/licenses/LICENSE-2.0"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="Apache 2.0"></a>
</p>

**InfiniteMITM** is an interactive MITM (Man-In-The-Middle) CLI for **Halo Infinite** which enables you to intercept and modify the game's requests and responses on the fly.

**Note:** While this interactive CLI is primarily designed to work with Halo Infinite and [halowaypoint.com](https://www.halowaypoint.com), other Halo titles may also be supported.

<img alt="InfiniteMITM CLI" title="InfiniteMITM CLI" src="./assets/preview.gif?v=3" width="800" />

## Disclaimer

This application is designed to enhance your experience and should not impact other players' experiences. However, by using this app, you acknowledge and agree that you are solely responsible for any actions taken with this app, including any potential bans or other consequences that may result. The developers of this application are not responsible for any disciplinary actions taken by game administrators or any other parties.

## Installation

Download and unzip one of the files from the [latest release](https://github.com/Alexis-Bize/InfiniteMITM/releases/latest) for your current OS.

## Documentation

-   [Install Root Certificate](/docs/Install-Root-Certificate.md)
-   [Override Requests](/docs/Override-Requests.md)
-   [Force Kill Proxy](/docs/Force-Kill-Proxy.md)

## SmartCache

**Halo Infinite** tends to repeatedly request the same content (images, binaries, JSON, etc.) as you play, which can be quite frustrating. To reduce network usage and enhance the game's performance, **InfiniteMITM** introduces a solution called **SmartCache**, which, once enabled, will automatically cache this content in memory.

-   Documentation: [SmartCache](/docs/SmartCache.md)

## Examples (Snippet)

-   [Server Selection](/examples/server-selection)
-   [Flags Override](/examples/flags-override)

## But Why Not Fiddler?

While **Fiddler** remains a leading MITM (Man-In-The-Middle) proxy tool, it can quickly become overwhelming due to the numerous traffic from various processes, making it quite complex to analyze requests and responses. **InfiniteMITM**, however, focuses solely on the traffic related to **Halo services**, providing an easy way to view and **rewrite everything on the fly** through a simple configuration file (`mitm.yaml`).

## Building From Source

### Requirements:

-   Generate your own `InfiniteMITMRootCA.pem`, `InfiniteMITMRootCA.key` and `InfiniteMITMRootCA.cer` certificates in the `cert/` directory using `openssl` (CN=InfiniteMITMRootCA).
-   Install the generated certificates on your machine.

### Build Script:

```shell
$ chmod +x ./scripts/build.sh
$ ./scripts/build.sh
```

## Known Issues

-   **Windows** may flag the application as a threat or virus due to a **false positive** (https://go.dev/doc/faq#virus)
-   As the application will create a local server to intercept traffic on a local port (1337), it **must be run as an administrator**.
-   The default Windows terminal (`cmd.exe`) won't render this application nicely.
    -   We recommend using the new [Windows Terminal](https://www.microsoft.com/p/windows-terminal-preview/9n0dx20hk701) instead.

## Licence

[Apache Version 2.0](/LICENCE)
