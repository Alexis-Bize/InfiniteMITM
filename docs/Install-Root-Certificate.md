# Install Root Certificate

When starting the application for the first time, you must install the root certificate. This certificate allows the CLI to listen for the game's network traffic. The CLI should provide the option to do this. If you prefer to install the root certificate manually, please follow the steps below.

## Install the Root Certificate manually

### Windows

1. Download the [**InfiniteMITMRootCA.cer**](/cert/InfiniteMITMRootCA.cer) certificate (View RAW) and double-click on it.
2. Click **Install Certificate**, install it for the **Current User** and click **Next**.
3. Select **Place all certificates in the following store**.
4. Select **Trusted Root Certification Authorities**.
5. Click **Next** and then **Finish**.
6. After you close the "**The import was successful**" alert, restart **InfiniteMITM**.
