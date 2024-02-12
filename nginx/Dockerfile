FROM alpine:latest
RUN apk --update add \nginx  && rm -rf /var/cache/apk/*
ADD nginx.conf /etc/nginx/nginx.conf
RUN mkdir -p /run/nginx
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
