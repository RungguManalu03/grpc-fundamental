-Instal dependensi
go get github.com/gin-gonic/gin
go get google.golang.org/grpc
go get google.golang.org/protobuf/proto
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get github.com/jackc/pgx/v5/pgxpool

-Buat direktori untuk Protobuf
mkdir proto

-Buat file user.proto di dalam folder proto

-Kompilasi file Protobuf
protoc --go_out=. --go-grpc_out=. proto/user.proto
