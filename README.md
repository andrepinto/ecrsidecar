# Ecr Side Car

Is a utility that gets a token for ecr. The token is refreshed in a defined time interval

# Why use

If you are using kubernetes and ecr then you should have felt the limitation in using the ecr credentials (expires after 12h). Ecr-side-car keeps your credentials always updated.

# Config file

```yaml
interval: 50s
accessKeyId: "XXXXXXX" # AWS_ACCESS_KEY_ID
secretAccessKey: "XXXXXXX" # AWS_SECRET_ACCESS_KEY
region: "eu-west-1"
tokenFile: "ecr-registry.token"
registryID: "00000000" # AWS ACCOUNT ID
```

# RUN

```shell
docker run -v $(pwd)/tokens:/opt/tokens -v $(pwd)/config.yaml:/opt/config/ecrsidecar/config.yaml andrepinto/ecrsidecar
```