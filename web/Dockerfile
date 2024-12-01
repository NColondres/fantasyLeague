FROM node:alpine

WORKDIR /app

COPY web/ /app

RUN rm -rf dist && npm ci && npm run build

EXPOSE 4173

ENTRYPOINT [ "npx", "vite", "preview", "--host"]