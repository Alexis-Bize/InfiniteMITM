# Server Selection

Override the available servers by keeping only one active to force its selection.

## Installation

-   Create a `resources/json` directory in `~/InfiniteMITM`.
-   Download and move [`server-selection.json`](./resources/json/server-selection.json) into the `json` directory.
-   Remove any unwanted servers from the file.
-   Copy and paste the content of `mitm.yaml` into your own file, adapting it to your current configuration.

### Schema (Example)

```json
[
  {
    "region": "EastUs",
    "serverUrl": "pfmsqosprod2-0.eastus.cloudapp.azure.com"
  },
  {
    "region": "EastUs2",
    "serverUrl": "pfmsqosprod2-0.eastus2.cloudapp.azure.com"
  }
]
```

### Server Regions and URLs

- **South Africa North**
  - Server URL: `pfmsqosprod2-0.southafricanorth.cloudapp.azure.com`
  
- **West Europe**
  - Server URL: `pfmsqosprod2-0.westeurope.cloudapp.azure.com`
  
- **Australia East**
  - Server URL: `pfmsqosprod2-0.australiaeast.cloudapp.azure.com`
  
- **East Asia**
  - Server URL: `pfmsqosprod2-0.eastasia.cloudapp.azure.com`
  
- **Southeast Asia**
  - Server URL: `pfmsqosprod2-0.southeastasia.cloudapp.azure.com`
  
- **Brazil South**
  - Server URL: `pfmsqosprod2-0.brazilsouth.cloudapp.azure.com`
  
- **East US**
  - Server URL: `pfmsqosprod2-0.eastus.cloudapp.azure.com`
  
- **East US 2**
  - Server URL: `pfmsqosprod2-0.eastus2.cloudapp.azure.com`
  
- **Central US**
  - Server URL: `pfmsqosprod2-0.centralus.cloudapp.azure.com`
  
- **North Central US**
  - Server URL: `pfmsqosprod2-0.northcentralus.cloudapp.azure.com`
  
- **South Central US**
  - Server URL: `pfmsqosprod2-0.southcentralus.cloudapp.azure.com`
  
- **West US**
  - Server URL: `pfmsqosprod2-0.westus.cloudapp.azure.com`
  
- **West US 2**
  - Server URL: `pfmsqosprod2-0.westus2.cloudapp.azure.com`
  
- **North Europe**
  - Server URL: `pfmsqosprod2-0.northeurope.cloudapp.azure.com`
  
- **Japan East**
  - Server URL: `pfmsqosprod2-0.japaneast.cloudapp.azure.com`
