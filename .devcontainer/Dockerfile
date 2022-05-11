FROM golang:1.17

# Install operator-sdk 
RUN export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac) && \
    export OS=$(uname | awk '{print tolower($0)}') && \
    export OPERATOR_SDK_DL_URL=https://github.com/operator-framework/operator-sdk/releases/download/v1.19.1 && \
    curl -LO ${OPERATOR_SDK_DL_URL}/operator-sdk_${OS}_${ARCH} && \
    chmod +x operator-sdk_${OS}_${ARCH} &&  mv operator-sdk_${OS}_${ARCH} /usr/local/bin/operator-sdk

# Install Golang extension package tools
RUN go install github.com/ramya-rao-a/go-outline@latest         && \
    go install github.com/cweill/gotests/gotests@latest         && \
    go install github.com/fatih/gomodifytags@latest             && \
    go install github.com/josharian/impl@latest                 && \
    go install github.com/haya14busa/goplay/cmd/goplay@latest   && \
    go install github.com/go-delve/delve/cmd/dlv@latest         && \
    go install honnef.co/go/tools/cmd/staticcheck@latest        && \
    go install golang.org/x/tools/gopls@latest



