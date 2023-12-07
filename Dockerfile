FROM docker.io/library/node:20-alpine3.17 AS ui
WORKDIR /workspace
COPY console/environment-hub .
RUN npm i
RUN npm run build-only

FROM golang:1.21 AS server
WORKDIR /workspace
COPY . .
RUN rm -rf cmd/data/assets
COPY --from=ui /workspace/dist/ cmd/data/
ENV CGO_ENABLED=0
RUN go build -o env-hub .

FROM alpine:3.18.5
WORKDIR /workspace
COPY --from=server /workspace/env-hub .
CMD [ "/workspace/env-hub" ]
