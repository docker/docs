FROM debian:unstable

ENV UEFI_REPO		https://github.com/freebsd/uefi-edk2.git
ENV UEFI_BRANCH 	bhyve/UDK2014.SP1
ENV UEFI_COMMIT		cf812813913f0cba82d24efae20bda1641f2e903
			# Really "anything >= 49"
ENV UEFI_TOOLCHAIN	GCC49
			# DEBUG or RELEASE
ENV UEFI_BUILD		RELEASE

RUN apt-get update && apt-get -y upgrade && apt-get -y install \
	build-essential \
	iasl \
	git \
	nasm \
	python \
	uuid-dev \
&& rm -rf /var/lib/apt/lists/*

RUN git clone -b "$UEFI_BRANCH" "$UEFI_REPO" /uefi && \
    cd /uefi && \
    git checkout -q "$UEFI_COMMIT"

RUN cd /uefi && \
    find . -name License.txt | while read l ; do \
        echo $l ; echo ; cat $l ; echo ; echo ; \
    done > /LICENSE

RUN jobs=$(nproc) && \
    cd /uefi && \
    make -C BaseTools

# EDK build requires bash
RUN /bin/bash -c 'jobs=$(nproc) && \
    cd /uefi && \
    . ./edksetup.sh && \
    build -a X64 -b "$UEFI_BUILD" -t "$UEFI_TOOLCHAIN" -p BhyvePkg/BhyvePkgX64.dsc && \
    cp Build/BhyveX64/"$UEFI_BUILD"_"$UEFI_TOOLCHAIN"/FV/BHYVE.fd /UEFI.fd'
