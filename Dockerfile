FROM golang:latest 

ARG SERVER_BUILD_VERSION=1.0.0
ARG SERVER_BUILD_COMMIT=

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -ldflags "-X main.buildVersion=${SERVER_BUILD_VERSION} -X main.buildCommit=${SERVER_BUILD_COMMIT} -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" -o main ./cmd/server/
CMD ["/app/main"]