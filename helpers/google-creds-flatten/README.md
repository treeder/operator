Converts a Google credentials JSON file to a base64 string using  RFC 4648.

## Usage

```sh
docker run --rm -v $PWD:/envs -w /envs treeder/operator:google-creds-flatten google-creds.json
```

To decode it in Go: 

```go
gcJson, err := base64.StdEncoding.DecodeString(encoded)
if err != nil {
    logrus.WithError(err).Errorln("base64 decode error")
    return nil, err
}
```
