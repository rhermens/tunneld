# tunneld
A multi-developer http tunneling client/server.  
Centralized tunneld server to receive webhooks and forward them to multiple tunnel clients (developers).  

```mermaid
flowchart LR
        Web["Webhooks"]
        Tunneld
        C1["Tunnel client 1"]
        Dev1["dev server 1"]
        C2["Tunnel client 2"]
        Dev2["dev server 2"]
        C3["Tunnel client 3"]
        Dev3["dev server 3"]

        Web -- HTTP --> Tunneld

        Tunneld -- SSH --> C1 -- HTTP --> Dev1
        Tunneld -- SSH --> C2 -- HTTP --> Dev2
        Tunneld -- SSH --> C3 -- HTTP --> Dev3
```

## Tunnel client usage

```bash
tunnel \
    --registry <tunneld-server-address> \
    --host <local-dev-server-address> \
    --port <local-dev-server-port>
```
Options can be configured via yaml file, defaulting to `~/.config/tunnel/tunnel.yaml`, or `./tunnel.yaml`.  
Refer to `examples/tunnel.yaml` for an example configuration file.


## Tunneld server usage
`tunneld serve` start a tunneld server containing a ssh registry and http proxy.  

refer to `examples/tunneld.yaml`, for an example configuration file.  

The registries host key will be automatically generated on first start and stored in the `registry.ssh.host_key_path` directory.  

### Keystores
Keystores are used to authenticate tunnel clients.
- GitHub organizations, all public keys of members will be allowed to connect.  
- Yaml configuration, a list of allowed public keys.
