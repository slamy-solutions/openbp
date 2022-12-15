FROM node:19-alpine3.15 as build

WORKDIR /app
COPY modules/tools/services/gui/package.json ./
COPY modules/tools/services/gui/yarn.lock ./
RUN yarn

COPY modules/tools/services/gui ./
RUN yarn build

# production environment
FROM --platform=${BUILDPLATFORM:-linux/amd64} nginx:1.23-alpine
COPY modules/tools/services/gui/nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]