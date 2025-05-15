# Lino S3

High level S3 wrapper for Go.

## Usage

```plaintext
[Must]LoadS3(dsn string) -> LinoS3
[Must]LoadS3Bucket(dsn string) -> LinoS3Bucket
[Must]LoadS3Object(dsn string) -> LinoS3Object
[Must]LoadS3Path(dsn string) -> LinoS3Path

NewLinoS3(client *s3.Client) -> LinoS3

LinoS3
  -> UseInterceptors(...interceptors) -> LinoS3 (sub entities will inherit these interceptors)
  -> Bucket(bucketName) -> LinoS3Bucket

LinoS3Bucket
  -> UseInterceptors(...interceptors) -> LinoS3Bucket
  -> Object(key) -> LinoS3Object
  -> SubPath(subPath) -> LinoS3Path
  -> Piece(objectKey, start, end) -> LinoS3Piece
  -> PieceByKey(pieceKey) -> LinoS3Piece

LinoS3Path
  -> UseInterceptors(...interceptors) -> LinoS3Path
  -> Object(key) -> LinoS3Object
  -> SubPath(subPath) -> LinoS3Path

LinoS3Object
  -> UseInterceptors(...interceptors) -> LinoS3Path
  -> HasInterceptors() -> bool
  -> Piece(start, end) -> LinoS3Piece

  -> Get() -> *s3.GetObjectOutput
  -> Put(s3.PutObjectInput) -> *s3.PutObjectOutput
     // key and bucket in input will be overrided,so it is not necessary to set them
  -> Upload(s3manager.UploadInput) -> *s3manager.UploadOutput
  -> Delete() => *s3.DeleteObjectOutput

  -> ReadTo(io.Writer)
  -> ReadBuffer() -> buffer
  -> ReadString() -> string
  -> ReadJSON(&value) // encoding/json
  -> ReadCBOR(&value) // github.com/fxamacker/cbor/v2
  -> ReadCSV(&value)  // github.com/gocarina/gocsv

  -> WriteFrom(io.Reader, contentType? string)
  -> WriteBuffer(buffer, contentType? string)
  -> WriteString(string, contentType? string)
  -> WriteJSON(&value)
  -> WriteCBOR(&value)
  -> WriteCSV(&value)

  -> Key() -> string // pieceKey

LinoS3Piece
  -> Get() -> *s3.GetObjectOutput

  -> ReadTo(io.Writer)
  -> ReadBuffer() -> buffer
  -> ReadString() -> string
  -> ReadJSON(&value)
  -> ReadCBOR(&value)
  -> ReadCSV(&value)

  -> Key() -> string // pieceKey
```

### interceptors

```go
import (
  "github.com/linolabx/lino_s3"
  "github.com/linolabx/lino_s3/interceptors"
)

s3 := lino_s3.NewLinoS3(client).Bucket("bucket")
blocksDir := bucket.SubPath("blocks-v1").UseInterceptors(interceptors.Gzip)

block SomeStruct
blockObject := blocksDir.Object("somekey.gz").ReadCBOR(&block);

block.SomeField = "new value"
blockObject.WriteCBOR(&block)
```

### utils

```go
import "github.com/linolabx/lino_s3/utils"

utils.HashSplit("123456") // "e1/0a/dc"
utils.HashSplit("123456", 5) // "e1/0a/dc/39/49"

utils.HashPrefix("hello.txt") // "2e/54/14/hello.txt"
utils.HashPrefix("hello.txt", 3) // "2e/54/14/hello.txt"
```

### structures

LinoS3Map

```go
import (
  "github.com/linolabx/lino_s3"
  "github.com/linolabx/lino_s3/structures/lino_s3_map"
)

bucket := s3.Bucket("bucket")
obj := bucket.Object("my-map")

someMap := lino_s3_map.NewMap(obj)
someMap.Set("foo1", []byte{"bar1"})
someMap.Set("foo2", []byte{"bar2"})
// save index of items in the first bytes, and then items
someMap.Save()

// read bytes from memory
someMap.Get("foo1") // []byte{"bar1"}

// directly read required bytes from s3
pieceKey, _ := someMap.PieceKey("foo1")
bucket.PieceByKey(pieceKey).ReadBuffer() // []byte{"bar1"}

// load map index, and then read required bytes from s3
secondMap := lino_s3_map.LoadMap(obj)
secondMap.Get("foo1") // []byte{"bar1"}

// load entire map, and then read required bytes from memory
thirdMap := lino_s3_map.NewMap(obj)
thirdMap.Get("foo1") // []byte{"bar1"}
```

## Develop

start minio, and default access key and secret key is `minioadmin:minioadmin`

```bash
docker run --name lino-s3-test --rm -p 9000:9000 -e MINIO_ROOT_USER=minioadmin -e MINIO_ROOT_PASSWORD=minioadmin minio/minio server /data

docker exec -it lino-s3-test mc alias set minio http://localhost:9000 minioadmin minioadmin
docker exec -it lino-s3-test mc mb minio/lino-stor
```
