# Server Selection

Override the available servers by keeping only desired ones to force their selection.

## Preview

<p align="center">
    <img alt="InfiniteMITM - Server Selection" title="InfiniteMITM - Server Selection" src="./preview.png?v=1" width="720" />
</p>

## Installation

-   Create a `resources/json` directory in `~/InfiniteMITM`.
-   Download and move [`server-selection.json`](./resources/json/server-selection.json) into the `json` directory.
-   Replace all `serverUrl` by the desired ones.
-   Copy and paste the content of `mitm.yaml` into your own file, adapting it to your current configuration.

### Schema (Example)

In the following example (partial, must not be used as is), all `serverUrl` instances are replaced with the `NorthCentralUs` one

```json
[
  {
    "region": "SouthAfricaNorth",
    "serverUrl": "pfmsqosprod2-0.northcentralus.cloudapp.azure.com"
  },
  {
    "region": "EastUs2",
    "serverUrl": "pfmsqosprod2-0.northcentralus.cloudapp.azure.com"
  }
]
```

## Notice

The returned estimated wait time may not be impacted by your changes, but the match search will be.

### Server Regions and URLs

- **South Africa North**
  - Region: `SouthAfricaNorth`
  - Server URL: `pfmsqosprod2-0.southafricanorth.cloudapp.azure.com`
  
- **West Europe**
  - Region: `WestEurope`
  - Server URL: `pfmsqosprod2-0.westeurope.cloudapp.azure.com`
  
- **Australia East**
  - Region: `AustraliaEast`
  - Server URL: `pfmsqosprod2-0.australiaeast.cloudapp.azure.com`
  
- **East Asia**
  - Region: `EastAsia`
  - Server URL: `pfmsqosprod2-0.eastasia.cloudapp.azure.com`
  
- **Southeast Asia**
  - Region: `SoutheastAsia`
  - Server URL: `pfmsqosprod2-0.southeastasia.cloudapp.azure.com`
  
- **Brazil South**
  - Region: `BrazilSouth`
  - Server URL: `pfmsqosprod2-0.brazilsouth.cloudapp.azure.com`
  
- **East US**
  - Region: `EastUs`
  - Server URL: `pfmsqosprod2-0.eastus.cloudapp.azure.com`
  
- **East US 2**
  - Region: `EastUs2`
  - Server URL: `pfmsqosprod2-0.eastus2.cloudapp.azure.com`
  
- **Central US**
  - Region: `CentralUs`
  - Server URL: `pfmsqosprod2-0.centralus.cloudapp.azure.com`
  
- **North Central US**
  - Region: `NorthCentralUs`
  - Server URL: `pfmsqosprod2-0.northcentralus.cloudapp.azure.com`
  
- **South Central US**
  - Region: `SouthCentralUs`
  - Server URL: `pfmsqosprod2-0.southcentralus.cloudapp.azure.com`
  
- **West US**
  - Region: `WestUs`
  - Server URL: `pfmsqosprod2-0.westus.cloudapp.azure.com`
  
- **West US 2**
  - Region: `WestUs2`
  - Server URL: `pfmsqosprod2-0.westus2.cloudapp.azure.com`
  
- **North Europe**
  - Region: `NorthEurope`
  - Server URL: `pfmsqosprod2-0.northeurope.cloudapp.azure.com`
  
- **Japan East**
  - Region: `JapanEast`
  - Server URL: `pfmsqosprod2-0.japaneast.cloudapp.azure.com`
