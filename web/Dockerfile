FROM node:alpine
ARG ENV_FILE

WORKDIR /app

COPY web/ /app
COPY ${ENV_FILE} /app/.env

RUN rm -rf dist && npm ci && npm run build

EXPOSE 4173

ENTRYPOINT [ "npx", "vite", "preview", "--host"]
