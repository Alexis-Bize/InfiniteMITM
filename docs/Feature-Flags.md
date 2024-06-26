# Feature Flags

**Halo Infinite** relies on various `featureflags` to enable and disable in-game features and options such as Campaign, Multiplayer, and more. These feature flags are returned by the settings service (`settings.svc.halowaypoint.com/hi/featureflags`) and can be edited at any time, although a game reboot is required for changes to take effect.

## Supported Feature Flags

-   `AcademyDrillsEnabled`
-   `AcademyEnabled`
-   `AcademyTrainingEnabled`
-   `AcademyTutorialEnabled`
-   `ActiveRosterEnabled`
-   `ArmorMenuEnabled`
-   `BodyAndAIMenuEnabled`
-   `Bot_EnforceCountLimit`
-   `CampaignEnabled`
-   `challenges_post_match_enabled`
-   `CustomGameEnabled`
-   `CustomizationEnabled`
-   `CustomsBrowserEnabled`
-   `DebugConsoleEnabled`
-   `DiskCacheOfflinePackageEnabled`
-   `director_enable_third_person_camera_game_option`
-   `enableFilmLatencyEmulation`
-   `enableObserverLatencyEmulation`
-   `ForgeEnabled`
-   `IsPostBattlePassGrindEnabled`
-   `ManageGameEnabled`
-   `matchmaking_uneven_teams_early_exit_enabled`
-   `matchmaking_uneven_teams_participant_changes_enabled`
-   `MissionRestartEnabled`
-   `MultiplayerEnabled`
-   `ObserverTeamNameOverridesEnabled`
-   `PlayerProfilesEnabled`
-   `PresentationMenuEnabled`
-   `RecommendedFilesEnabled`
-   `RecommendedFilmsEnabled`
-   `RecommendedMapsEnabled`
-   `RecommendedModesEnabled`
-   `SavedFilm_EnableLegacyDeathCam`
-   `SavedFilm_EnableLegacyDeathCamInObserverMode`
-   `SavedFilmPlaybackRenderDebugEnabled`
-   `SearchRegionSettingEnabled`
-   `SelectServerEnabled`
-   `ShowMirroredOptionsPracticeMenu`
-   `SpotlightEnabled`
-   `TaskBarEnabled`
-   `TelemetryCombatEncountersEnabled`
-   `TheaterEnabled`
-   `UseIslandGeoInMenus`
-   `VehiclesMenuEnabled`
-   `WeaponsMenuEnabled`

## How to Customize Existing Flags

Thanks to the amazing work by **soupstream**, [InfiniteVariantTool](https://github.com/soupstream/InfiniteVariantTool) allows you to convert bond binary files (format used by the game) into .xml files and vice versa.

### Schema (Example)

The first `set` is for enabled features, and the second `set` is for disabled ones. The following example will enable additional options in Practice and disable the Forge menu.

```xml
<?xml version="1.0" encoding="utf-8"?>
<?InfiniteVariantTool version="0.6.0.0"?>
<struct>
  <set id="0" type="string">
    <string>ShowMirroredOptionsPracticeMenu</string>
  </set>
  <set id="1" type="string">
    <string>ForgeEnabled</string>
  </set>
</struct>
```
