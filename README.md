Install
----------
```
git clone https://0xacab.org/leap/bitmask-vpn.git
cd bitmask-vpn
make generate
go build -buildmode=c-archive -o lib/libgoshim.a gui/backend.go
make build
```
