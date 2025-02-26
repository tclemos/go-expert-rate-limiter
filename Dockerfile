FROM golang:1.23.5 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/

FROM scratch as server
WORKDIR /app
COPY --from=build /app/server .
COPY --from=build /app/cmd/server/.env .
ENTRYPOINT ["./server"]

