To integrate the proxy configuration within your environment and leverage it for HTTP and HTTPS traffic, you can set environment variables on your system. This configuration directs your HTTP and HTTPS traffic through the proxy server, which then uses the provided tenant name and API key for authentication and access control, based on the ACL (Access Control List) configurations associated with that tenant.

Below is a brief guide on how to set up these environment variables for using the proxy:

---

# Configuring Proxy Environment Variables

To route your HTTP and HTTPS traffic through the proxy server, you can use the `http_proxy` and `https_proxy` environment variables. This setup requires specifying the tenant name and API key within the proxy URL to authenticate and authorize your requests according to the tenant's ACL settings.

## Proxy URL Format

The proxy URL should be formatted as follows:

```
http_proxy=http://tenant_name:api_key@proxy-addr
https_proxy=http://tenant_name:api_key@proxy-addr
```

- **tenant_name:** Your unique tenant identifier.
- **api_key:** The API key provided for authentication.
- **proxy-addr:** The address of the proxy server, including the port number if necessary (e.g., `proxy-server.example.com:8080`).

## Setting Environment Variables

### On Unix/Linux/macOS

Open a terminal and use the export command to set the `http_proxy` and `https_proxy` variables:

```bash
export http_proxy="http://tenant_name:api_key@proxy-addr"
export https_proxy="https://tenant_name:api_key@proxy-addr"
```

To make these changes permanent, you can add them to your shell's initialization script, such as `~/.bashrc` or `~/.zshrc`.

## Verifying the Configuration

After setting the environment variables, you can verify that your traffic is being routed through the proxy by accessing a web service that shows your IP address or by checking the logs of the proxy server to see if your requests are appearing as expected.

## Conclusion

By configuring the `http_proxy` and `https_proxy` environment variables with your tenant name and API key, you can securely and efficiently route your HTTP and HTTPS traffic through the proxy server. This setup ensures that your requests comply with the ACL rules specified for your tenant, enhancing the security and manageability of your network traffic.

---

Replace `tenant_name`, `api_key`, and `proxy-addr` with the actual values provided for your specific setup.

# ACL Manager Usage Guide

The ACL (Access Control List) Manager is a core component of our proxy application, designed to control access based on client requests to specific tenants. Each tenant's access control rules are defined in a JSON file, which includes a whitelist and a blacklist of domain patterns.

## JSON File Structure

Each tenant's JSON file should follow this structure:

```json
{
  "Whitelist": [
    "*example.com"
  ],
  "Blacklist": [
    "restricted.example.com"
  ]
}
```

- **Whitelist:** Requests matching any pattern in the whitelist are allowed, unless also matched by the blacklist.
- **Blacklist:** Requests matching any pattern in the blacklist are blocked, taking precedence over the whitelist.

## Loading Tenant Lists

The `ACLManager` loads these JSON files from the `acl-Data-Path` directory, specified in the application's configuration. Each file should be named after the tenant it represents and contain the JSON structure shown above.

## Request Evaluation Logic

When a request is received, the `ACLManager` performs the following steps to determine if it should be allowed:

1. **Extract Hostname:** The hostname from the incoming request is extracted, including stripping out the port if present.
2. **Blacklist Check:** The request's hostname is checked against the blacklist patterns. If a match is found, the request is immediately blocked, and no further checks are performed.
3. **Whitelist Check:** If the request's hostname does not match any blacklist pattern, it is then checked against the whitelist patterns. If a match is found, the request is allowed.
4. **Default Block:** If the request does not match any whitelist pattern (or if no patterns are defined for the tenant), the request is blocked by default.

## Pattern Matching

- Domain patterns in the lists can include wildcards (`*`) for matching multiple subdomains or specific characters.
- The matching is case-insensitive and expects the entire hostname to match the pattern.

## Examples

- **Whitelist Example:** If the whitelist contains `*.example.com`, then requests to `sub.example.com` and `example.com` are allowed, but `sub.restricted.example.com` is not allowed if `restricted.example.com` is in the blacklist.
- **Blacklist Example:** If the blacklist contains `restricted.example.com`, any request to this domain is blocked, regardless of the whitelist.

## Implementation Details

- The ACLManager compiles the patterns into regular expressions for efficient matching.
- The blacklist takes precedence over the whitelist. If a hostname matches both, it is considered blocked.
- Logging is provided at various stages to aid in debugging and understanding the decision process for allowing or blocking requests.

## Conclusion

The ACLManager provides a flexible and powerful way to manage access control for different tenants in the proxy application. By defining clear and concise rules in the JSON files, administrators can easily control which domains are allowed or blocked, ensuring secure and efficient operation of the proxy service.