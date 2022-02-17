FROM moby/buildkit:v0.9.3
WORKDIR /api
COPY api README.md /api/
ENV PATH=/api:$PATH
ENTRYPOINT [ "/bhojpur/api" ]