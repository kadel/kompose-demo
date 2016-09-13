FROM centos:7

COPY ./caddy /caddy
COPY goenvtemplator /usr/local/bin/goenvtemplator
COPY Caddyfile /opt/http/
COPY public /opt/http/public
RUN chmod a+w  /opt/http/public/scripts/example.js

WORKDIR /opt/http

EXPOSE 8080

ENV API_URL http://localhost:3000/api

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
