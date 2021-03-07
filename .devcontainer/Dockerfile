# ubuntu 20.04 (lts) amd64
FROM ubuntu@sha256:e3d7ff9efd8431d9ef39a144c45992df5502c995b9ba3c53ff70c5b52a848d9c

RUN apt update && apt upgrade -y

# common tools
RUN \
    apt install -y \
        bash-completion \
        curl \
        git \
        jq \
        net-tools \
        tar \
        time \
        unzip \
        zip \
        wget \
        build-essential \
        sudo

# go1.16(latest)
RUN \
    apt install -y && \
    curl -fsSL "https://golang.org/dl/go1.16.linux-amd64.tar.gz" | tar -xz -C /usr/local
ENV \
    PATH $PATH:/usr/local/go/bin

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN \
    groupadd --gid $USER_GID ${USERNAME} \
    && useradd -s /bin/bash --uid ${USER_UID} --gid ${USER_GID} -m ${USERNAME} \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME\
    && chmod 0440 /etc/sudoers.d/$USERNAME

USER vscode

ENTRYPOINT [ "bash", "-eux", "/entrypoint.sh" ]