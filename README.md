# Infinite MITM

<p>
    <img alt="InfiniteMITM" title="InfiniteMITM" src="./assets/logo.png" width="256">
    <br>
    <a href="https://github.com/Alexis-Bize/InfiniteMITM/releases"><img src="https://img.shields.io/github/v/release/Alexis-Bize/InfiniteMITM?include_prereleases" alt="Latest Release"></a>
    <a href="https://www.apache.org/licenses/LICENSE-2.0"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="Apache 2.0"></a>
</p>

**InfiniteMITM** is a MITM (man-in-the-middle) CLI for **Halo Infinite** which enables you to intercept and modify the game's requests and responses on the fly.

<img alt="InfiniteMITM CLI" title="InfiniteMITM CLI" src="./assets/preview.png?v=1" />

## Disclaimer

This application is designed to enhance your experience and should not impact other players' experiences. However, by using this app, you acknowledge and agree that you are solely responsible for any actions taken with this app, including any potential bans or other consequences that may result. The developers of this application are not responsible for any disciplinary actions taken by game administrators or any other parties.

## Installation

Download and unzip one of the files from the [latest release](https://github.com/Alexis-Bize/InfiniteMITM/releases/latest) for your current OS.

## Documentation

-   [Install Root Certificate](/docs/Install-Root-Certificate.md)
-   [Override Requests](/docs/Override-Requests.md)

## Examples (Snippet)

-   [Server Selection](/examples/server-selection)
-   [Flags Override](/examples/flags-override)

## Building From Source

```shell
$ chmod +x ./scripts/build.sh
$ ./scripts/build.sh
```

## Known Issues

-   **Windows** may flag the application as a threat or virus due to a **false positive** (https://go.dev/doc/faq#virus)
-   As the application will create a local server to intercept traffic on a local port (1337), it **must be run as an administrator**.

## Licence

[Apache Version 2.0](/LICENCE)
