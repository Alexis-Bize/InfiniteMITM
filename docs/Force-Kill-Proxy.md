# Force Kill Proxy

In some cases, such as during a fatal error generally due to an invalid `mitm.yaml` configuration, the created proxy may remain active even though the application has been shut down. This can cause errors like `ERR_PROXY_CONNECTION_FAILED` during your internet browsing.

To address this issue, you can force the proxy to stop by restarting the application and selecting **Force Kill Proxy**.

Alternatively, you can also run a terminal (e.g., `cmd.exe` on Windows) as an administrator and type the following command: `netsh winhttp reset proxy`
