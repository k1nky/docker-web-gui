FROM alpine
WORKDIR /opt
ADD ./ui /opt/ui
ADD ./dockerboard /opt/dockerboard
COPY ./*.pem /opt/
RUN apk add --no-cache \
        libc6-compat
CMD ["dockerboard"]
EXPOSE 8000
