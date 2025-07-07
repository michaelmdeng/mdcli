I want to make `mdcli` installable/available as a docker image. The image should contain
the `mdcli` executable and can be used similar to the busybox image, where the executable
is copied in as part of a multi-stage build.

```dockerfile
COPY --from ghcr.io/michaelmdeng/mdcli/mdcli:latest /bin/mdcli .
```

# Requirements

* Image statically linked and built with `CGO_ENABLED=0`
* Image uses a minimal distroless layer with only the `mdcli` executable
* Leverages multi-stage builds to separate build and runtime deps
* Supports multi-platform, specifically linux/amd64 and linux/arm64
* Images should be tagged w/ the `mdcli` version, (ex. `v0.0.1`) and `latest`
* Image is built and pushed on updates to the main branch, as part of existing Github
* Add Makefile targets that support docker functionality locally:
    * `build-image` target that builds the docker image, takes an optional tag argument
    * `run-image` target that builds the docker image, takes an optional tag argument
    * `publish-image` target that pushes the tag to the relevant registry, take required
      registry tag as an argument
* Image is hosted on ghcr
