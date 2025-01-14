##################################################
## BASE
##################################################
FROM golang:1.23.4-bookworm AS base
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .

##################################################
## DEVELOPMENT
##################################################
FROM base AS development
RUN go install github.com/air-verse/air@latest

CMD ["air"]

##################################################
## BUILD
##################################################
FROM base AS build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o out/main ./main.go

##################################################
## PRODUCTION
##################################################
FROM scratch AS production
COPY --from=build /app/out/main .

CMD ["./main"]
