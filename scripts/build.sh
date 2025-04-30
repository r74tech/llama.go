#!/usr/bin/env bash

set -e

if [ ! -d "./core" ]; then
    echo "Must run in the root directory of the project"
    exit 1
fi

if [ ! -d "./core/llama.cpp" ]; then
    git submodule update --init --recursive
    echo "Update llama.cpp"
fi

cmake --version


coreDir=$(pwd)/core
buildDir=$(pwd)/build

echo "core dir:" ${coreDir}
echo "build dir:" ${buildDir}

cmake -DCMAKE_BUILD_TYPE=Release -G "Unix Makefiles" -S $coreDir -B $buildDir
cmake --build $buildDir --target llama_core -- -j 9

# go
go version

GITVER=$(git rev-parse --short=7 HEAD)
GITDIRTY=$(git diff --quiet || echo '-dirty')
GITVERSION="${GITVER}${GITDIRTY}"
versionBuild="github.com/Qitmeer/llama.go/version.Build=dev-${GITVERSION}"

export CGO_ENABLED=1
export LD_LIBRARY_PATH=./build/lib
go build -ldflags "-X ${versionBuild}" -o ./build/bin/llama

echo "Output executable file:${buildDir}/bin/llama"
./build/bin/llama --version


