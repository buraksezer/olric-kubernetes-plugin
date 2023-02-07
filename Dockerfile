FROM golang:latest as build

WORKDIR /src/
COPY . /src/
RUN go mod download
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -buildmode=plugin -o /usr/lib/olric-kubernetes-plugin.so

FROM olricio/olricd:v0.5.4
COPY --from=build /usr/lib/olric-kubernetes-plugin.so /usr/lib/olric-kubernetes-plugin.so