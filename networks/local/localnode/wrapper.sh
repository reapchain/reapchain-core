#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/reapchain-core/${BINARY:-reapchain-core}
ID=${ID:-0}
LOG=${LOG:-reapchain-core.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'reapchain-core' E.g.: -e BINARY=reapchain-core_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export TMHOME="/reapchain-core/node${ID}"

if [ -d "`dirname ${TMHOME}/${LOG}`" ]; then
  "$BINARY" "$@" | tee "${TMHOME}/${LOG}"
else
  "$BINARY" "$@"
fi

chmod 777 -R /reapchain-core

