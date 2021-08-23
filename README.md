# Frost
Basic webservice reverse proxy and process manager.

To build run:

```go build ./cmd/watchdog```

Requires a local instance of Redis, and a copy of config.json (as shown below) to be present and filled with valid data.

```
{
   "baseURL":"<the base hostname you want services to be available from: ie alargerobot.dev>", 
   "vaultKeyID":"",
   "vaultAddress": "<hostname or ip address and port of an instance of Hashicorp Vault>", 
   "vaultAppRoleName": "<the name of a Vault approle for Frost itself to use>", 
   "vaultServicesAppRole": "<the name of a Vault approle for Frost services to use>"
}
```

The UI app (in management/ui) is an Angular 7.3 app and thus requires a NodeJS install.

Vault is an optional requirement, that can be used to provide hosted services with easier access to encryption capabilities. Its use and configuration is best described at https://vaultproject.io
