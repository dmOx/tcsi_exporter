FROM golang:1.14 AS build
ADD . /app
WORKDIR /app/cmd/tcsi_exporter
RUN CGO_ENABLED=0 go build .

FROM alpine:3.12
COPY --from=build /app/cmd/tcsi_exporter/tcsi_exporter /app/tcsi_exporter
EXPOSE 9115
CMD /app/tcsi_exporter