if [ -z "$CHART_HOME" ]; then
    echo "error: CHART_HOME should be initialized"
    exit 1
fi

OUTPUT=${CHART_HOME}/output
OUTPUT_BIN=${OUTPUT}/bin
PULSARCTL_VERSION=v2.8.2.1
PULSARCTL_BIN=${CHART_HOME}/scripts/pulsarctl-amd64-linux/pulsarctl
export PATH=${CHART_HOME}/scripts/pulsarctl-amd64-linux/pulsarctl/plugins:${PATH}

discoverArch() {
  ARCH=$(uname -m)
  case $ARCH in
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
    arm64) ARCH="arm64";;
  esac
}

discoverArch
OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

test -d "$OUTPUT_BIN" || mkdir -p "$OUTPUT_BIN"

function pulsar::verify_pulsarctl() {
    if test -x "$PULSARCTL_BIN"; then
        return
    fi
    return 1
}

function pulsar::ensure_pulsarctl() {
    if pulsar::verify_pulsarctl; then
        return 0
    fi
    echo "Get pulsarctl install.sh script ..."
    ${CHART_HOME}/scripts/download/install.sh --user --version ${PULSARCTL_VERSION}
}


