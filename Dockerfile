FROM golang:1.22-alpine as builder

WORKDIR /app
COPY . .

RUN go mod tidy && go mod verify
RUN CGO_ENABLED=0 go build -o worker-service cmd/worker-service/main.go

FROM gcr.io/distroless/static-debian12 as runner

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/worker-service .
COPY --from=builder --chown=nonroot:nonroot /app/.env .
COPY --from=builder --chown=nonroot:nonroot /app/assets/privacy-policy.html .

EXPOSE 8080

ENTRYPOINT ["./worker-service"]