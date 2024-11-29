Install
----------
```
git clone --depth=1 https://0xacab.org/leap/bitmask-vpn.git
cd bitmask-vpn
export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
PROVIDER=riseup make vendor
PROVIDER=riseup QMAKE=qmake6 LRELEASE=/usr/lib64/qt6/bin/lrelease RELEASE=yes make build -j 1
```
![01](https://github.com/artenax/riseup-vpn/assets/107228652/774e79a5-b454-49dd-bbd7-44b0e2b23b61) ![02](https://github.com/artenax/riseup-vpn/assets/107228652/e48f149e-1c90-48b6-b327-6b90c6e0b7d9)
![03](https://github.com/artenax/riseup-vpn/assets/107228652/d7f0027c-d913-440c-a55a-46487899f9f4)  ![riseup](https://github.com/artenax/riseup-vpn/assets/107228652/92a64bbb-b3ef-4ef2-a769-ccbe3c3822fb)
