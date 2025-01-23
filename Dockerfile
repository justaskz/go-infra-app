##################################################
## BASE
##################################################
FROM golang:1.23.4-bookworm AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

##################################################
## DEVELOPMENT
##################################################
FROM base AS development
RUN go install github.com/air-verse/air@latest
COPY . .
CMD ["air"]

##################################################
## BUILD
##################################################
FROM base AS build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY . .
RUN go build -o /tmp/server cmd/infra-app/main.go

##################################################
## PRODUCTION
##################################################
FROM scratch AS production
COPY --from=build /tmp/server .

CMD ["./server"]
