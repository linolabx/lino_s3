package lino_s3_test

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/linolabx/lino_s3/internal/test"
)

func TestIoPut(t *testing.T) {
	text := "local-" + uuid.NewString()
	filename := "local.txt"
	object := test.GetS3Object("text-v1", filename)

	t.Cleanup(func() {
		os.Remove(filename)
		object.Delete()
	})

	if err := os.WriteFile(filename, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}

	reader, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	if err := object.WriteFrom(reader); err != nil {
		t.Fatal(err)
	}

	if str, err := object.ReadString(); str != text {
		t.Fatal(err)
	}
}

func TestIoGet(t *testing.T) {
	text := "remote:" + uuid.NewString()
	filename := "remote.txt"
	object := test.GetS3Object("text-v1", filename)

	t.Cleanup(func() {
		os.Remove(filename)
		object.Delete()
	})

	if err := object.WriteString(text); err != nil {
		t.Fatal(err)
	}

	writer, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}

	if err := object.ReadTo(writer); err != nil {
		t.Fatal(err)
	}

	if str, err := os.ReadFile(filename); string(str) != text {
		t.Fatal(err)
	}
}

func TestJSON(t *testing.T) {

	violetObj := test.GetS3Object("json-v1", "violet.json")
	gilbertObj := test.GetS3Object("json-v1", "gilbert.json")

	t.Cleanup(func() {
		violetObj.Delete()
		gilbertObj.Delete()
	})

	violetObj.WriteJSON(test.Violet)
	gilbertObj.WriteJSON(test.Gilbert)

	if str, _ := violetObj.ReadString(); str != `{"name":"Violet Evergarden","mail":"violet@ch-postal.com","password":"password","Entropy":"18446744073709551615"}` {
		t.Fatal("Violet JSON write failed")
	}

	if str, _ := gilbertObj.ReadString(); str != `{"name":"Gilbert Bougainvillea","mail":"","password":""}` {
		t.Fatal("Gilbert JSON write failed")
	}

	var _violet test.User
	var _gilbert test.User

	if violetObj.ReadJSON(&_violet); _violet.Entropy != test.Violet.Entropy {
		t.Fatal("Violet JSON read failed")
	}

	if gilbertObj.ReadJSON(&_gilbert); _gilbert.Entropy != test.Gilbert.Entropy {
		t.Fatal("Gilbert JSON read failed")
	}
}

func TestCBOR(t *testing.T) {
	violetObj := test.GetS3Object("cbor-v1", "violet.cbor")

	t.Cleanup(func() {
		violetObj.Delete()
	})

	violetObj.WriteCBOR(test.Violet)

	if buf, _ := violetObj.ReadBuffer(); hex.EncodeToString(buf) != `a4017156696f6c6574204576657267617264656e027476696f6c65744063682d706f7374616c2e636f6d036870617373776f7264041bffffffffffffffff` {
		t.Fatal("Violet CBOR write failed")
	}

	var _violet test.User
	if violetObj.ReadCBOR(&_violet); _violet.Entropy != 18_446_744_073_709_551_615 {
		t.Fatal("Violet CBOR read failed")
	}
}

func TestCSV(t *testing.T) {
	usersObj := test.GetS3Object("csv-v1", "users.csv")

	t.Cleanup(func() {
		usersObj.Delete()
	})

	usersObj.WriteCSV([]*test.User{test.Violet, test.Gilbert})

	if str, _ := usersObj.ReadString(); str != `Name,Email,Password,Entropy
Violet Evergarden,violet@ch-postal.com,password,18446744073709551615
Gilbert Bougainvillea,,,0
` {
		t.Fatal("CSV write failed")
	}

	var users []test.User
	if err := usersObj.ReadCSV(&users); err != nil {
		t.Fatal(err)
	}

	if users[0].Entropy != 18_446_744_073_709_551_615 || users[1].Entropy != 0 {
		t.Fatal("CSV read failed")
	}
}
