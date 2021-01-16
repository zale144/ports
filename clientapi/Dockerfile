FROM golang:1.13.4

WORKDIR /app

COPY ./clientapi/go.mod ./clientapi/go.sum ./
RUN go mod download

COPY ./clientapi ./

EXPOSE 4000
