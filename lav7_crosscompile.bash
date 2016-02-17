# lav7 cross-compiling script
# Use golang-crosscompile(https://github.com/davecheney/golang-crosscompile) to install Go crosscompiler for platforms.
# For Android support, compile Go toolchain for Android.
# https://jasonplayne.com/programming-2/how-to-cross-compile-golang-for-android
# And, set $NDK_CC to arm-linux-androideabi-gcc toolchain directory.
#
# Usage:
#   lav7-crosscompile $GOOS-$GOARCH: cross-compiles lav7 to ~/builds/$GOOS/$GOARCH
#      - If you are compiling it for linux-arm/android-arm, each ARMv5, ARMv6, ARMv7 binary will be created.
#   lav7-build-all: Build for all possible platforms to ~/builds
#   lav7-build-publish: Execute lav7-build-all and move to ~/share/lav7/builds/{current date}

type setopt >/dev/null 2>&1 && setopt shwordsplit
PLATFORMS="darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm windows/386 windows/amd64 openbsd/386 openbsd/amd64 android/arm"
ARMS="5 6 7"

function lav7-crosscompile {
    cd ${GOPATH}/src/github.com/L7-MCPE/lav7
    local GOOS=${1%/*}
    local GOARCH=${1#*/}
    mkdir -p ~/builds/${GOOS}/${GOARCH}
    local CMD="go-${GOOS}-${GOARCH}"
    local ACMD=""

    if [ "$GOOS" = "windows" ]; then
        ACMD="${CMD} build -o lav7.exe -ldflags '-X \"github.com/L7-MCPE/lav7.GitCommit=$(git rev-parse --verify HEAD)\" -X \"github.com/L7-MCPE/lav7.BuildTime=$(LC_ALL=en_US.utf8 date -u)\"' l7start/lav7.go"
        echo $ACMD
        eval $ACMD || return 1
        mv lav7.exe ~/builds/${GOOS}/${GOARCH}
        rm -Rf $GOOS
    elif [ "$CMD" = "go-linux-arm" ]; then
        for ARM in $ARMS; do
            mkdir -p ~/builds/${GOOS}/${GOARCH}/ARMv${ARM}
            ACMD="GOARM=${ARM} $CMD build -o lav7 -ldflags '-X \"github.com/L7-MCPE/lav7.GitCommit=$(git rev-parse --verify HEAD)\" -X \"github.com/L7-MCPE/lav7.BuildTime=$(LC_ALL=en_US.utf8 date -u)\"' l7start/lav7.go"
            echo $ACMD
            eval $ACMD || return 1
            mv lav7 ~/builds/${GOOS}/${GOARCH}/ARMv${ARM}
            rm -Rf $ARM
        done
    elif [ "$CMD" = "go-android-arm" ]; then
        for ARM in $ARMS; do
            mkdir -p ~/builds/${GOOS}/${GOARCH}/ARMv${ARM}
            ACMD="CC_FOR_TARGET=$NDK_CC GOARM=${ARM} $CMD build -o lav7 -ldflags '-X \"github.com/L7-MCPE/lav7.GitCommit=$(git rev-parse --verify HEAD)\" -X \"github.com/L7-MCPE/lav7.BuildTime=$(LC_ALL=en_US.utf8 date -u)\"' l7start/lav7.go"
            echo $ACMD
            eval $ACMD || return 1
            mv lav7 ~/builds/${GOOS}/${GOARCH}/ARMv${ARM}
            rm -Rf $ARM
        done
    else
        ACMD="${CMD} build -o lav7 -ldflags '-X \"github.com/L7-MCPE/lav7.GitCommit=$(git rev-parse --verify HEAD)\" -X \"github.com/L7-MCPE/lav7.BuildTime=$(LC_ALL=en_US.utf8 date -u)\"' l7start/lav7.go"
        echo $ACMD
        eval $ACMD || return 1
        mv lav7 ~/builds/${GOOS}/${GOARCH}
        rm -Rf $GOOS
    fi
}

function lav7-build-all {
    local FAILS=""
    for PLATFORM in $PLATFORMS; do
        eval lav7-crosscompile ${PLATFORM} || FAILS="$FAILS $PLATFORM"
    done
    if [ "$FAILS" != "" ]; then
        echo "*** lav7-build-all failed on: ${FAILS}"
    fi
    cd ~/builds
}

function lav7-build-publish {
    mkdir -p ~/share/lav7/builds
    eval "lav7-build-all"
    for i in $(find . -type f -print); do sha256sum "$i"; done > SHA256SUM
    cd ~
    rm -Rf ~/share/lav7/builds/`LC_ALL=en_US.utf-8 date -u +%Y-%m-%d`
    mv builds ~/share/lav7/builds/`LC_ALL=en_US.utf-8 date -u +%Y-%m-%d`
}
