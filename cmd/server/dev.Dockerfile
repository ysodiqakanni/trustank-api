FROM golang:1.19-alpine
WORKDIR /app

# COPY go.mod, go.sum and download the dependencies
COPY go.* ./
RUN go mod download

# COPY All things inside the project and build
COPY . .

WORKDIR /app/cmd/server
RUN go build -o /trustankbizapi

WORKDIR /app

EXPOSE 8080
CMD ["/trustankbizapi"]