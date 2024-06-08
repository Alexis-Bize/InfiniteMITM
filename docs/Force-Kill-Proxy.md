# Force Kill Proxy

In some cases, such as during a fatal error, the created proxy may remain active even though the application has been shut down, which can cause errors like `NET::ERR_CERT_AUTHORITY_INVALID` during your internet browsing. To address this issue, you can force the proxy to stop by restarting the application and selecting **Force Kill Proxy**.

Alternatively, you can also run a terminal (e.g., `cmd.exe` on Windows) as an administrator and type the following command:

`netsh winhttp reset proxy`
