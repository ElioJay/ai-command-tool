param([string]$Target = "all")

$Version = "0.1.0"
$Binary  = "aict"
New-Item -ItemType Directory -Force dist | Out-Null

function Build($goos, $goarch, $suffix) {
    $env:GOOS   = $goos
    $env:GOARCH = $goarch
    $out = "dist/$Binary-$goos-$goarch$suffix"
    Write-Host "Building $out..."
    go build -o $out ./cmd/aict/
    $env:GOOS   = ""
    $env:GOARCH = ""
}

switch ($Target) {
    "windows" { Build "windows" "amd64" ".exe" }
    "macos"   { Build "darwin"  "amd64" ""; Build "darwin" "arm64" "-arm64" }
    "linux"   { Build "linux"   "amd64" "" }
    "test"    { go test ./... }
    "portable-windows" {
        Build "windows" "amd64" ".exe"
        $dir = "dist/portable-windows"
        New-Item -ItemType Directory -Force $dir | Out-Null
        New-Item -ItemType Directory -Force "$dir/.aict" | Out-Null
        Copy-Item "dist/$Binary-windows-amd64.exe" "$dir/$Binary.exe"
        Compress-Archive -Path $dir -DestinationPath "dist/$Binary-portable-windows-amd64.zip" -Force
        Write-Host "便携包已生成：dist/$Binary-portable-windows-amd64.zip"
    }
    default {
        Build "windows" "amd64" ".exe"
        Build "darwin"  "amd64" ""
        Build "darwin"  "arm64" "-arm64"
        Build "linux"   "amd64" ""
    }
}
