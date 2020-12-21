ARG VERSION

FROM centos:${VERSION}

# Needs to be defined after FROM
ARG USERNAME
ARG AWS_REGION

RUN useradd -ms /bin/bash ${USERNAME} && \
# Install all the language specific packages
    yum groupinstall -y "Development tools" && \
    yum install -y golang && \
    yum install -y python3-pip && \
    pip3 install --upgrade pip && \
    pip3 install awscli && \
# Install Go extensions for VSCode
    go get -v github.com/sqs/goreturns && \
    go get -v github.com/rogpeppe/godef && \
    go get -v github.com/ramya-rao-a/go-outline && \
    go get -v github.com/mdempsky/gocode && \
    go get -v github.com/uudashr/gopkgs/v2/cmd/gopkgs && \
    go get -v github.com/stamblerre/gocode && \
    go get -v golang.org/x/lint/golint && \
    go get -v golang.org/x/tools/cmd/gorename && \
    go get -v github.com/go-delve/delve/cmd/dlv && \
    go get -v golang.org/x/tools/gopls && \
# Install project-specific dependencies
    go get -v github.com/aws/aws-lambda-go/lambda && \
    go get -v github.com/aws/aws-sdk-go/aws && \
    go get -v github.com/stretchr/testify/assert && \
    go get -v github.com/golang/mock/mockgen && \
    go get -v github.com/golang/mock/gomock && \
    go get -v github.com/gusaul/go-dynamock && \
# Move files to user directory
    cp -r /root/go /home/${USERNAME} && \
    chown -R ${USERNAME}. /home/${USERNAME}/go && \
    ln -s /home/${USERNAME}/go/bin/mockgen /usr/local/bin
USER ${USERNAME}
WORKDIR /home/${USERNAME}
