# llama_build.ps1
$ErrorActionPreference = "Stop"

if (-Not (Test-Path "./core" -PathType Container)) {
    Write-Host "Must run in the root directory of the project"
    exit 1
}

if (-Not (Test-Path "./core/llama.cpp/src" -PathType Container)) {
    git submodule update --init --recursive
    Write-Host "Update llama.cpp"
}

cmake --version

$coreDir = (Get-Location).Path + "\core"
$buildDir = (Get-Location).Path + "\build"

Write-Host "core dir: $coreDir"
Write-Host "build dir: $buildDir"

# cuda
$nvcc = Get-Command nvcc -ErrorAction SilentlyContinue
$cudaPath = "C:\Program Files\NVIDIA GPU Computing Toolkit\CUDA"

$cudaTag=""
$cudaCmake=""

if ($nvcc -or (Test-Path "$cudaPath")) {
    Write-Host "Try use CUDA"
    $cudaCmake="-DGGML_CUDA=ON"
    $cudaTag="-tags=cuda"
}

cmake -DCMAKE_BUILD_TYPE=Release $cudaCmake -G "Unix Makefiles" -S $coreDir -B $buildDir
cmake --build $buildDir --target llama_core -- -j 9

# Go
go version

$GITVER = git rev-parse --short=7 HEAD
$GITDIRTY = if (git diff --quiet) { "" } else { "-dirty" }
$GITVERSION = "${GITVER}${GITDIRTY}"
$versionBuild = "github.com/Qitmeer/llama.go/version.Build=dev-$GITVERSION"
$env:CGO_ENABLED = "1"
$env:LD_LIBRARY_PATH = "./build/lib"
go build $cudaTag -ldflags "-X $versionBuild" -o ./build/bin/llama.exe

Write-Host "Output executable file: $buildDir/bin/llama.exe"
& "$buildDir/bin/llama.exe" --version







