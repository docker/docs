To update nsenter.tar

```
docker pull justincormack/nsenter1
docker tag justincormack/nsenter1 d4w/nsenter
docker save -o .\src\Resources\nsenter.tar d4w/nsenter
```