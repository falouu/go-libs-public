# go libraries

## CONTRIBUTE

Setup the repo for contribution:
```
./tools/install-git-config
```

### How to push to main without a PR

```
g ps origin main:refs/heads/main-rc && while ! g ps; do sleep 20; done
```
