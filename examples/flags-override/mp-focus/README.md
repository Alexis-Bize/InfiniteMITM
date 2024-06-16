# Flags Override - MP Focus

Disable all options (campaign, customization, store, etc.) and keep only Multiplayer (Matchmaking and Custom game) and Academy.

## Preview

<p align="center">
    <img alt="InfiniteMITM - Flags Override" title="InfiniteMITM - Flags Override" src="./preview.jpg?v=1" width="720" />
</p>

## Disabled Feature Flags

-   `ArmorMenuEnabled`
-   `BodyAndAIMenuEnabled`
-   `BuyCreditsEnabled`
-   `CampaignEnabled`
-   `CommunityEnabled`
-   `CustomizationEnabled`
-   `CustomsBrowserEnabled`
-   `ForgeEnabled`
-   `HCSStoreEnabled`
-   `PresentationMenuEnabled`
-   `RecommendedFilesEnabled`
-   `RecommendedFilmsEnabled`
-   `RecommendedMapsEnabled`
-   `RecommendedModesEnabled`
-   `SpotlightEnabled`
-   `StoreEnabled`
-   `TheaterEnabled`
-   `VehiclesMenuEnabled`
-   `WeaponsMenuEnabled`

Please refer to our [documentation](/docs/Feature-Flags.md) for more details about these values.

## Installation

-   Create a `resources/bin/flags` directory in `~/InfiniteMITM`.
-   Download and move [`mp-focus.bin`](./resources/bin/flags/mp-focus.bin) into the `flags` directory.
-   Copy and paste the content of `mitm.yaml` into your own file, adapting it to your current configuration.
    -   Documentation: [Use Provided Examples](/docs/Use-Provided-Examples.md)

## Notice

By overriding existing flags, you may also disable ones that have been recently added server-side.
