FROM golang:1.14 as build

RUN apt-get update && apt-get install -y ninja-build

# TODO: Змініть на власну реалізацію системи збірки
RUN go get -u github.com/vladShadow/Architecture-LAB2/build/cmd/bood

WORKDIR /go/src/GO23
COPY . .

RUN CGO_ENABLED=0 bood

# ==== Final image ====
FROM alpine:3.11
WORKDIR /opt/GO23
COPY entry.sh ./
COPY --from=build /go/src/GO23/out/bin/* ./
ENTRYPOINT ["/opt/GO23/entry.sh"]
CMD ["server"]
