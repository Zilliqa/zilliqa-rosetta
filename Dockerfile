FROM golang:1.14 as rosetta-build-stage
WORKDIR /app
COPY . ./
RUN go build main.go


FROM ubuntu:16.04

# ====================
# Scilla Prerequisite
# ====================
ARG MAJOR_VERSION=0
ARG SCILLA_COMMIT_OR_TAG=v0.7.1
WORKDIR /scilla/${MAJOR_VERSION}

RUN apt-get update \
    && apt-get install -y software-properties-common \
    && add-apt-repository ppa:tah83/secp256k1 -y \
    && add-apt-repository ppa:avsm/ppa -y \
    && apt-get update && apt-get install -y --no-install-recommends \
    curl \
    git \
    cmake \
    build-essential \
    m4 \
    ocaml \
    opam \
    pkg-config \
    zlib1g-dev \
    libgmp-dev \
    libffi-dev \
    libssl-dev \
    libsecp256k1-dev \
    libboost-system-dev \
    libpcre3-dev \
    && rm -rf /var/lib/apt/lists/*

ENV OCAML_VERSION 4.07.1

RUN git clone --recurse-submodules https://github.com/zilliqa/scilla .
RUN git checkout ${SCILLA_COMMIT_OR_TAG}
RUN git status
RUN make opamdep-ci \
    && echo '. ~/.opam/opam-init/init.sh > /dev/null 2> /dev/null || true ' >> ~/.bashrc \
    && eval $(opam env) && \
    make

# ====================
# Zilliqa Node
# ====================
# Format guideline: one package per line and keep them alphabetically sorted
RUN apt-get update \
    && apt-get install -y software-properties-common \
    && add-apt-repository ppa:tah83/secp256k1 -y \
    && apt-get update && apt-get install -y --no-install-recommends \
    autoconf \
    build-essential \
    ca-certificates \
    cmake \
    # curl is not a build dependency
    curl \
    git \
    golang \
    libboost-filesystem-dev \
    libboost-program-options-dev \
    libboost-system-dev \
    libboost-test-dev \
    libboost-python-dev \
    libcurl4-openssl-dev \
    libevent-dev \
    libjsoncpp-dev \
    libjsonrpccpp-dev \
    libleveldb-dev \
    libmicrohttpd-dev \
    libminiupnpc-dev \
    libsnappy-dev \
    libssl-dev \
    libtool \
    ocl-icd-opencl-dev \
    pkg-config \
    python \
    python-pip \
    python3-dev \
    python3-pip \
    python3-setuptools \
    libsecp256k1-dev \
    && rm -rf /var/lib/apt/lists/*

# Manually input tag or commit, can be overwritten by docker build-args
ARG COMMIT_OR_TAG=17d581f
ARG REPO=https://github.com/Zilliqa/Zilliqa.git
ARG SOURCE_DIR=/zilliqa
ARG BUILD_DIR=/build
ARG INSTALL_DIR=/usr/local
ARG BUILD_TYPE=RelWithDebInfo
ARG EXTRA_CMAKE_ARGS=

RUN git clone ${REPO} ${SOURCE_DIR} \
    && git -C ${SOURCE_DIR} checkout ${COMMIT_OR_TAG}

RUN pip3 install -r /zilliqa/docker/requirements3.txt

RUN cmake -H${SOURCE_DIR} -B${BUILD_DIR} -DCMAKE_BUILD_TYPE=${BUILD_TYPE} \
        -DCMAKE_INSTALL_PREFIX=${INSTALL_DIR} ${EXTRA_CMAKE_ARGS} \
    && cmake --build ${BUILD_DIR} -- -j$(nproc) \
    && cmake --build ${BUILD_DIR} --target install \
    && rm -rf ${BUILD_DIR}

ENV LD_LIBRARY_PATH=${INSTALL_DIR}/lib




# ====================
# Rosetta Deployment
# ====================
WORKDIR /rosetta
COPY --from=rosetta-build-stage /app/main /rosetta/main
COPY --from=rosetta-build-stage /app/config.local.yaml /rosetta/config.local.yaml
EXPOSE 8080
ENTRYPOINT ["/bin/bash"]
