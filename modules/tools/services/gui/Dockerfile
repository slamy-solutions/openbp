FROM --platform=${BUILDPLATFORM:-linux/amd64} node:18-alpine3.18 as builder
RUN yarn global add @quasar/cli
WORKDIR /app
COPY modules/tools/services/gui/package.json package.json
COPY modules/tools/services/gui/yarn.lock yarn.lock
RUN yarn
COPY modules/tools/services/gui/ .
RUN quasar build

FROM nginx:1.25.1-alpine as prod
COPY --from=builder /app/dist/spa /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]