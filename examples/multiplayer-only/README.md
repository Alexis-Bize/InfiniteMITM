# Multiplayer Only

Disable all options (campaign, customization, store, etc.) and keep only Multiplayer (Matchmaking and Custom game) and Academy.

## Preview

<p align="center">
    <img alt="InfiniteMITM - No Progression Assets" title="InfiniteMITM - No Progression Assets" src="./preview.jpg?v=1" width="720" />
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

Please refer to our [documentation](/blob/main/docs/Feature-Flags.md) for more details about these values.

## Installation

-   Create a `resources/bin/flags` directory in `~/InfiniteMITM`.
-   Download and move [`multiplayer-only.bin`](./resources/bin/flags/multiplayer-only.bin) into the `flags` directory.
-   Copy and paste the content of `mitm.yaml` into your own file, adapting it to your current configuration.

## Notice

By overriding existing flags, you may also disable ones that have been recently added server-side.