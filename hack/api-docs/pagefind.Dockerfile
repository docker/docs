# Rebuild the existing search tool for Linux hosts with 16 KB or 64 KB pages.
FROM rust:1.94-alpine3.23 AS build
RUN apk add --no-cache build-base
ENV JEMALLOC_SYS_WITH_LG_PAGE=16
RUN cargo install pagefind --version 1.5.2 --locked --root /out
FROM scratch
COPY --from=build /out/bin/pagefind /pagefind
