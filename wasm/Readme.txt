Go- Web assembly
-------------------
    1. path\wasm\web>set GOOS=js
	2. path\wasm\web>set GOARCH=wasm
	3. path\wasm\web>go build -o c:\build\dp.wasm



CC=/android-ndk-r25c/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android28-clang"
CXX="/android-ndk-r25c/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android28-clang++"
GOOS=android GOARCH=arm64 CGO_ENABLED=1 go build

Go- Android
---------------------
    android-armv7a:
    	path\wasm\mobile>set CGO_ENABLED=1
    	path\wasm\mobile>set GOOS=android
    	path\wasm\mobile>set GOARCH=arm
    	path\wasm\mobile>set GOARM=7
        path\wasm\mobile>set CC=PATH\AppData\Local\Android\Sdk\ndk\23.1.7779620\toolchains\llvm\prebuilt\windows-x86_64\bin\armv7a-linux-androideabi21-clang
    	path\wasm\mobile>go build -buildmode=c-shared -o PATH\AndroidStudioProjects\dpwallet\app\src\main\jniLibs\armeabi-v7a\libgodp.so main.go

    android-arm64:
    	path\wasm\mobile>set CGO_ENABLED=1
    	path\wasm\mobile>set GOOS=android
    	path\wasm\mobile>set GOARCH=arm64
        path\wasm\mobile>set CC=PATH\AppData\Local\Android\Sdk\ndk\23.1.7779620\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang
    	path\wasm\mobile>go build -buildmode=c-shared -o PATH\AndroidStudioProjects\dpwallet\app\src\main\jniLibs\arm64-v8a\libgodp.so main.go

    android-x86:
    	path\wasm\mobile>set CGO_ENABLED=1
    	path\wasm\mobile>set GOOS=android
    	path\wasm\mobile>set GOARCH=386
        path\wasm\mobile>set CC=PATH\AppData\Local\Android\Sdk\ndk\23.1.7779620\toolchains\llvm\prebuilt\windows-x86_64\bin\i686-linux-android21-clang
    	path\wasm\mobile>go build -buildmode=c-shared -o PATH\AndroidStudioProjects\dpwallet\app\src\main\jniLibs\x86\libgodp.so main.go

    android-x86_64:
    	path\wasm\mobile>set CGO_ENABLED=1
    	path\wasm\mobile>set GOOS=android
    	path\wasm\mobile>set GOARCH=amd64
        path\wasm\mobile>set CC=PATH\AppData\Local\Android\Sdk\ndk\23.1.7779620\toolchains\llvm\prebuilt\windows-x86_64\bin\x86_64-linux-android21-clang
    	path\wasm\mobile>go build -buildmode=c-shared -o PATH\AndroidStudioProjects\dpwallet\app\src\main\jniLibs\x86_64\libgodp.so main.go

        android build.gradle
        ------------------------
        ndk {
            abiFilters "armeabi-v7a" , "arm64-v8a", "x86", "x86_64"
        }
