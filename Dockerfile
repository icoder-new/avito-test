FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o tma_api ./cmd/

FROM alpine
WORKDIR /app
COPY --from=build /app/.env .
COPY --from=build /app/config/config.yml /app/config/.
COPY --from=build /app/tma_api .
ENTRYPOINT [ "./tma_api" ]
