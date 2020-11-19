FROM golang:1.14 

ENV GO111MODULE=on
ENV STARTWITHDOCKER=true

ENV RAJA_ONGKIR_KEY=<Api-Key>
ENV JWT_SECRET=<Jwt-Secret>
ENV MONGO_URL=mongo:27017
ENV ELASTIC_URL=elastic:9200
ENV REDIS_URL=redis:6379

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN chmod +x ./wait-for-it.sh

RUN go build
