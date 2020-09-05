FROM alpine:3.12

ENV HELM_VERSION=3.3.1

RUN apk add --no-cache bash gawk sed grep bc coreutils git curl openssl jq tar && \
    curl -fSL --connect-timeout 10 "https://get.helm.sh/helm-v$HELM_VERSION-linux-amd64.tar.gz" | tar zxv --wildcards --strip 1 -C "/usr/local/bin" "*/helm" && \
    chmod 755 /usr/local/bin/helm && \
    helm plugin install https://github.com/chartmuseum/helm-push

ADD chartmuseum.sh /usr/local/bin/chartmuseum.sh
    
ENTRYPOINT /usr/local/bin/chartmuseum.sh
