# Web image: the React build served by an unprivileged nginx that also proxies /api.
# Behind a TLS-inspecting proxy, pass its CA: --secret id=npm_ca,src=ca.crt (and the
# standard HTTPS_PROXY build argument).
# Build from the project root: docker build -f deploy/web.Dockerfile .

FROM node:22-bookworm-slim AS build
WORKDIR /app
COPY package.json package-lock.json ./
COPY apps/api/package.json apps/api/
COPY apps/web/package.json apps/web/
RUN --mount=type=secret,id=npm_ca,required=false \
    if [ -f /run/secrets/npm_ca ]; then export npm_config_cafile=/run/secrets/npm_ca; fi; npm ci -w apps/web --include-workspace-root=false --no-audit --no-fund
COPY apps/web apps/web
RUN npm run build -w apps/web

FROM nginxinc/nginx-unprivileged:1.27-alpine
COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/apps/web/dist /usr/share/nginx/html
EXPOSE 8080
