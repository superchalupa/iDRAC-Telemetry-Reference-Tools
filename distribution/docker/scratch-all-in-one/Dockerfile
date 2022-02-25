ARG USER_ID
ARG GROUP_ID

FROM golang:1.17 as builder

# yes, this has to be here twice, afaict. Once globally before first FROM. Once in this build step after FROM
ARG USER_ID
ARG GROUP_ID

RUN DEBIAN_FRONTEND=noninteractive apt-get update && DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade

WORKDIR /build
COPY internal /build/internal
COPY cmd /build/cmd
COPY go.* /build/
RUN mkdir /out

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"  -o /out/dbdiscauth  ./cmd/dbdiscauth
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"  -o /out/configui    ./cmd/configui
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"  -o /out/redfishread ./cmd/redfishread
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"  -o /out/simpleauth  ./cmd/simpleauth
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"  -o /out/simpledisc  ./cmd/simpledisc

RUN echo ADDING telemetry user; \
      userdel -f telemetry;\
      getent group ${GROUP_ID}  && groupdel $(getent group ${GROUP_ID} | cut -d: -f1) ; \
      getent group telemetry && groupdel telemetry; \
      groupadd -g ${GROUP_ID} telemetry ;\
      useradd -l -u ${USER_ID} -g telemetry telemetry

# better to use scratch, but our entry point is a shell script now so need alpine for /bin/sh
#FROM scratch

FROM alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/
COPY --from=builder /out /app
COPY distribution/docker/scratch-all-in-one/run.sh /run.sh

USER telemetry:telemetry
ENTRYPOINT ["/run.sh"]
