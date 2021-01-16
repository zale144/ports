FROM golang:1.15.6

WORKDIR /app

COPY ./portdomainservice/go.mod ./portdomainservice/go.sum ./
RUN go mod download

COPY ./portdomainservice ./

EXPOSE 4000
