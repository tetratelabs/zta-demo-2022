FROM maven:3-jdk-8 as javabuilder
RUN mkdir /build/
COPY exploit/ /build
RUN cd /build && mvn package

FROM golang:1.16 as gobuilder
WORKDIR /
COPY go.* ./
COPY main.go .
RUN go mod download
COPY --from=javabuilder /build/target/log4shell-exploit-1.0-SNAPSHOT.jar .
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags app -o app -ldflags '-w' .

FROM scratch
COPY --from=gobuilder /app /app
ENTRYPOINT ["/app"]
EXPOSE 3000
EXPOSE 1389
