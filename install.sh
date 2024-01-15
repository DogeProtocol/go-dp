#!/bin/bash

apt-get install -qq -y pkg-config
apt-get install -qq -y unzip
apt-get update -qq -y

mkdir templibs
mkdir templibs/pkg-config

mkdir templibs/hybrid-pqc
mkdir templibs/liboqs

chmod 777 $PWD/.config/template/*

cp $PWD/.config/template/libhybridpqc-template.pc $PWD/templibs/pkg-config/libhybridpqc.pc
cp $PWD/.config/template/liboqs-template.pc $PWD/templibs/pkg-config/liboqs.pc

sed -i -e "s|\[INCLUDE_DIR\]|$PWD/templibs/hybrid-pqc/build/include|" $PWD/templibs/pkg-config/libhybridpqc.pc
sed -i -e "s|\[LIB_DIR\]|$PWD/templibs/hybrid-pqc|" $PWD/templibs/pkg-config/libhybridpqc.pc

sed -i -e "s|\[INCLUDE_DIR\]|$PWD/templibs/liboqs/build/include|" $PWD/templibs/pkg-config/liboqs.pc
sed -i -e "s|\[LIB_DIR\]|$PWD/templibs/liboqs|" $PWD/templibs/pkg-config/liboqs.pc

curl -Lo $PWD/templibs/hybrid-pqc/includes.zip https://github.com/DogeProtocol/hybrid-pqc/releases/download/v0.1.12/includes.zip
unzip $PWD/templibs/hybrid-pqc/includes.zip -d $PWD/templibs/hybrid-pqc
echo "dd3001667a7199d2d3faa2bbf1b848a832c71caf04b61c023b34545a142c08f7 $PWD/templibs/hybrid-pqc/includes.zip" | sha256sum --check  - || exit 1

curl -Lo $PWD/templibs/liboqs/includes.zip https://github.com/DogeProtocol/liboqs/releases/download/v0.0.4/includes.zip
unzip $PWD/templibs/liboqs/includes.zip -d $PWD/templibs/liboqs
echo "e04a39e332b169aad8370fbcd99aa8ab03ab5d0e621d711e78c6c9f6aa341d56 $PWD/templibs/liboqs/includes.zip" | sha256sum --check  - || exit 1

curl -Lo $PWD/templibs/hybrid-pqc/libhybridpqc.so.2 https://github.com/DogeProtocol/hybrid-pqc/releases/download/v0.1.12/libhybridpqc.so.2
echo "bc2cafd3f281bc2443ae00fd1a7daf79aa214d8b8839219b61314bab033ee5ff $PWD/templibs/hybrid-pqc/libhybridpqc.so.2" | sha256sum --check  - || exit 1

curl -Lo $PWD/templibs/liboqs/liboqs.so.5 https://github.com/DogeProtocol/liboqs/releases/download/v0.0.4/liboqs.so.5
echo "6694aaff32255faafab324011b7f5ea5ca0f527e0b901265597871dfb01ddf72 $PWD/templibs/liboqs/liboqs.so.5" | sha256sum --check  - || exit 1

echo " "
echo "Installation complete. To start building:"
echo "1) Switch to the go-dp folder."
echo "2) Set the following environment variable; you would want to add it to your bash profile."
echo " "
echo "   export PKG_CONFIG_PATH=$PWD/templibs/pkg-config"
echo " "
echo "3) Then run the following command: "
echo "4) go build -o YOUR_BUILD_DIR ./..."


