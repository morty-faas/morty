# Morty Function Registry

This repository contains the source code of the Morty function registry.

## Development

### Requirements

- `Go 1.19`

Run the following command one time to have dependencies ready : 

```
docker compose up -d

export AWS_ACCESS_KEY_ID=mortymorty
export AWS_SECRET_ACCESS_KEY=mortymorty
```

If you have a remote S3 Storage, specify the host with:
    
```bash
export AWS_ENDPOINT_URL=http://custom-s3.storage.com:9000
```

Run the project : 

```bash
make build/registry
./morty-registry
```
