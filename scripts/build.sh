#!/usr/bin/env bash

set -e

if [ ! -d "./llama.cpp" ]; then
    git submodule update --init --recursive
    echo "Update llama.cpp"
fi

cmake --version

cd ./../
projectDir=$(pwd)/llama.cpp
buildDir=$(pwd)/build

echo "project dir:" ${projectDir}
echo "build dir:" ${buildDir}

cmake -DCMAKE_BUILD_TYPE=Release -G "Unix Makefiles" -S $projectDir -B $buildDir
cmake --build $buildDir --target $2 -- -j 9

export CGO_ENABLED=1
go version

# 目录设置
cd /Users/jin/Applications/HalalChain/qitmeer/qng/cmd/qng
workDir=$(pwd)

binDir=$@
if [[ $binDir == "" ]]; then
   binDir="/Users/jin/Applications/HalalChain/qitmeer/bin"
fi

exeName="qng"


# 获取git版本
GITVER=$(git rev-parse --short=7 HEAD)
GITDIRTY=$(git diff --quiet || echo '-dirty')
GITVERSION="${GITVER}${GITDIRTY}"
versionBuild="github.com/Qitmeer/qng/version.Build=dev-${GITVERSION}"


echo "工作目录:" ${workDir} && \
echo "输出目录:" ${binDir} && \
${GOROOT}/bin/go build -ldflags "-X ${versionBuild}" -o ${binDir}/${exeName} && \
echo "编译完成." && \
${binDir}/${exeName} -V
