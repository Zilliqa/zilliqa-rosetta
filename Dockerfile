# ====================
# Rosetta Build Stage
# ====================
FROM golang:1.15 as rosetta-build-stage

ARG ROSETTA_COMMIT_OR_TAG=v1.0.0

WORKDIR /app
RUN git clone https://github.com/zilliqa/zilliqa-rosetta .
RUN git checkout ${ROSETTA_COMMIT_OR_TAG}
RUN git status
RUN go build main.go
RUN printf "================================================================================\n\n\nRosetta Build Stage Complete\n\n\n================================================================================\n"



# ====================
# Scilla Build Stage
# ====================
FROM ubuntu:18.04 as scilla-build-stage

ARG MAJOR_VERSION=0
ARG SCILLA_COMMIT_OR_TAG=v0.11.0

WORKDIR /scilla/${MAJOR_VERSION}

RUN apt-get update \
    && apt-get install -y software-properties-common \
    && add-apt-repository ppa:avsm/ppa -y \
    && apt-get update && apt-get install -y --no-install-recommends \
    curl \
    wget \
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
    libboost-test-dev \
    libpcre3-dev \
    && rm -rf /var/lib/apt/lists/*

ENV OCAML_VERSION 4.08.1

ENV PATH="/root/.local/bin:${PATH}"

RUN git clone --recurse-submodules https://github.com/zilliqa/scilla .
RUN git checkout ${SCILLA_COMMIT_OR_TAG}
RUN git submodule update --init --recursive
RUN git status
RUN bash scripts/install_cmake_ubuntu.sh \
    && make opamdep-ci \
    && echo '. ~/.opam/opam-init/init.sh > /dev/null 2> /dev/null || true ' >> ~/.bashrc \
    && eval $(opam env) && \
    make

RUN mkdir -p /scilla/0/bin2/ && cp -L /scilla/0/bin/* /scilla/0/bin2/ && strip /scilla/0/bin2/*
RUN printf "================================================================================\n\n\nScilla Build Stage Complete\n\n\n================================================================================\n"



# ====================
# Zilliqa Build Stage
# ====================
FROM ubuntu:18.04 as zilliqa-build-stage

# Format guideline: one package per line and keep them alphabetically sorted
RUN apt-get update \
    && apt-get install -y software-properties-common \
    && apt-get update && apt-get install -y --no-install-recommends \
    autoconf \
    build-essential \
    ca-certificates \
    cmake \
    wget \
    # curl is not a build dependency
    curl \
    git \
    golang \
    # rysnc bydefault gets installed with opam package of scilla.Better to explicitly
    # mention again as installation candidate
    rsync \
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
    python3-dev \
    python3-pip \
    python3-setuptools \
    libsecp256k1-dev \
    && rm -rf /var/lib/apt/lists/*

# Install cmake 3.19
RUN wget https://github.com/Kitware/CMake/releases/download/v3.19.3/cmake-3.19.3-Linux-x86_64.sh \	
    && mkdir -p "${HOME}"/.local \	
    && bash ./cmake-3.19.3-Linux-x86_64.sh --skip-license --prefix="${HOME}"/.local/	

# Include path to refer to latest version of cmake	
ENV PATH="/root/.local/bin:${PATH}"	

RUN cmake --version \	
    && rm cmake-3.19.3-Linux-x86_64.sh

# Manually input tag or commit, can be overwritten by docker build-args
ARG COMMIT_OR_TAG=v8.1.0
ARG REPO=https://github.com/Zilliqa/Zilliqa.git
ARG SOURCE_DIR=/zilliqa
ARG BUILD_DIR=/build
ARG INSTALL_DIR=/usr/local
ARG BUILD_TYPE=RelWithDebInfo
ARG EXTRA_CMAKE_ARGS=
ARG MONGO_INSTALL_DIR=${BUILD_DIR}/mongo

RUN git clone ${REPO} ${SOURCE_DIR} \
    && git -C ${SOURCE_DIR} checkout ${COMMIT_OR_TAG} \
    && cmake -H${SOURCE_DIR} -B${BUILD_DIR} -DCMAKE_BUILD_TYPE=${BUILD_TYPE} \
    -DCMAKE_INSTALL_PREFIX=${INSTALL_DIR} ${EXTRA_CMAKE_ARGS} \
    && cmake --build ${BUILD_DIR} -- -j$(nproc) \
    && cmake --build ${BUILD_DIR} --target install \
    && echo "built files:" && ls -lh ${BUILD_DIR} && echo "installed libs:" && ls -lh ${INSTALL_DIR}/lib \
    && echo "mongo files:" && ls -lh ${MONGO_INSTALL_DIR}/lib \
    # strip all exes
    && strip /usr/local/bin/grepperf \
       /usr/local/bin/zilliqad \
       /usr/local/bin/genkeypair \
       /usr/local/bin/signmultisig \
       /usr/local/bin/verifymultisig \
       /usr/local/bin/getpub \
       /usr/local/bin/getaddr \
       /usr/local/bin/genaccounts \
       /usr/local/bin/sendcmd \
       /usr/local/bin/gentxn \
       /usr/local/bin/restore \
       /usr/local/bin/gensigninitialds \
       /usr/local/bin/validateDB \
       /usr/local/bin/genTxnBodiesFromS3 \
       /usr/local/bin/getnetworkhistory \
       /usr/local/bin/isolatedServer \
       /usr/local/bin/getrewardhistory \
    #   /usr/local/bin/zilliqa \
    #   /usr/local/bin/data_migrate \
       /usr/local/lib/libSchnorr.so \
       /usr/local/lib/libCryptoUtils.so \
       /usr/local/lib/libNAT.so \
       /usr/local/lib/libCommon.so \
       /usr/local/lib/libTrie.so

RUN printf "================================================================================\n\n\nZilliqa Build Stage Complete\n\n\n\n================================================================================\n"



# ====================
# Actual Container
# ====================
FROM ubuntu:18.04


# --------------------
# Zilliqa Deployment
# --------------------
# install all necessary libraries
RUN apt-get update \
    && apt-get install -y software-properties-common \
    && apt-get update && apt-get install -y --no-install-recommends \
    # libs
    ca-certificates \
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
    python3-dev \
    python3-pip \
    python3-setuptools \
    libsecp256k1-dev \
    # tools
    curl \
    dnsutils \
    gdb \
    git \
    less \
    logrotate \
    net-tools \
    rsync \
    rsyslog \
    trickle \
    vim \
    && rm -rf /var/lib/apt/lists/*

# install all necessary libraries for python
COPY --from=zilliqa-build-stage /zilliqa/docker/requirements3.txt /zilliqa/docker/requirements3.txt
RUN pip3 install wheel \
    && pip3 install --no-cache-dir -r /zilliqa/docker/requirements3.txt \
    && update-alternatives --install /usr/bin/python python /usr/bin/python3 10 # set python3 as default instead python2

# make dirs fro scilla and zilliqa
RUN mkdir -p \
    /scilla/0/bin /scilla/0/src/stdlib \
    /zilliqa/scripts

ARG INSTALL_DIR=/usr/local
ARG MONGO_INSTALL_DIR=/build/mongo

# pour in scilla binaries
COPY --from=scilla-build-stage  /scilla/0/bin2            /scilla/0/bin
# pour in scilla conntract stdlibs
COPY --from=scilla-build-stage  /scilla/0/src/stdlib     /scilla/0/src/stdlib
# pour in zilliqa scripts
COPY --from=zilliqa-build-stage /zilliqa/scripts         /zilliqa/scripts
# pour in zilliqa binaries and dynnamic libs
COPY --from=zilliqa-build-stage ${INSTALL_DIR}/bin/*     ${INSTALL_DIR}/bin/
COPY --from=zilliqa-build-stage ${INSTALL_DIR}/lib/*.so* ${INSTALL_DIR}/lib/
COPY --from=zilliqa-build-stage ${MONGO_INSTALL_DIR}/lib/*.so* ${INSTALL_DIR}/lib/

ADD https://github.com/krallin/tini/releases/latest/download/tini /tini

ENV LD_LIBRARY_PATH=${INSTALL_DIR}/lib:${MONGO_INSTALL_DIR}/lib

RUN printf "================================================================================\n\n\nZilliqa Deployment Complete\n\n\n================================================================================\n"


# --------------------
# Rosetta Deployment
# --------------------
COPY --from=rosetta-build-stage /app/main /rosetta/main
COPY --from=rosetta-build-stage /app/testnet.config.local.yaml /rosetta/testnet.config.local.yaml
COPY --from=rosetta-build-stage /app/mainnet.config.local.yaml /rosetta/mainnet.config.local.yaml
EXPOSE 8080
RUN printf "================================================================================\n\n\nRosetta Deployment Complete\n\n\n================================================================================\n"


# --------------------
# Seed node setup
# --------------------
EXPOSE 4201
EXPOSE 4301
EXPOSE 4501
EXPOSE 33133

WORKDIR /run/zilliqa
COPY --from=rosetta-build-stage /app/seed_scripts/rosetta_seed_launch.sh /run/zilliqa/rosetta_seed_launch.sh
ENTRYPOINT ["/bin/bash", "rosetta_seed_launch.sh"]
#CMD ["/bin/bash"]
