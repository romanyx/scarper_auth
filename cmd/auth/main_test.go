package main

import (
	"crypto/rsa"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	txdb "github.com/DATA-DOG/go-txdb"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ory/dockertest"
	"github.com/romanyx/scraper_auth/internal/storage/postgres/schema"
	"github.com/romanyx/scraper_auth/kit/auth"
	"github.com/romanyx/scraper_auth/kit/docker"
	"github.com/romanyx/scraper_auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	caseTiemout = 3 * time.Second
)

var (
	db *sql.DB
	a  *auth.Authenticator
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		os.Exit(m.Run())
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %v", err)
	}

	pgDocker, err := docker.NewPostgres(pool)
	if err != nil {
		log.Fatalf("prepare postgres with docker: %v", err)
	}
	db = pgDocker.DB

	if err := schema.Migrate(db); err != nil {
		log.Fatalf("migrate schema: %v", err)
	}

	txdb.Register("pgsqltx", "postgres", fmt.Sprintf("password=test user=test dbname=test host=localhost port=%s sslmode=disable", pgDocker.Resource.GetPort("5432/tcp")))

	prvKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateRSAKey))
	if err != nil {
		log.Fatalf("private key: %v", err)
	}

	a, err = auth.NewAuthenticator(prvKey, privateRSAKeyID, "RS256", auth.NewSingleKeyFunc(privateRSAKeyID, prvKey.Public().(*rsa.PublicKey)))
	if err != nil {
		log.Fatalf("new authenticator: %v", err)
	}

	code := m.Run()

	db.Close()
	if err := pool.Purge(pgDocker.Resource); err != nil {
		log.Fatalf("could not purge postgres docker: %v", err)
	}

	os.Exit(code)
}

func newServer(db *sql.DB) (string, func()) {
	srv := setupServer(a, db, time.Hour)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterAuthServer(s, srv)
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return lis.Addr().String(), s.Stop
}

func postgresDB(t *testing.T) (db *sql.DB, teardown func() error) {
	dbName := fmt.Sprintf("db_%d", time.Now().UnixNano())
	db, err := sql.Open("pgsqltx", dbName)

	if err != nil {
		log.Fatalf("open postgres tx connection: %s", err)
	}

	return db, db.Close
}

// The key id we would have generated for the private below key
const privateRSAKeyID = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"

// Output of:
// openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
const privateRSAKey = `-----BEGIN PRIVATE KEY-----
MIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQDdiBDU4jqRYuHl
yBmo5dWB1j9aeDrXzUTJbRKlgo+DWDQzIzJQvackvRu8/f7B5cseoqmeJcmBu6pc
4DmQ+puGNHxzCyYVFSMwRtHBZvfWS3P+UqIXCKRAX/NZbLkUEeqPnn5WXjA+YXKk
sfniE0xDH8W22o0OXHOzRhDWORjNTulpMpLv8tKnnLKh2Y/kCL/4vo0SZ+RWh8F9
4+JTZx/47RHWb6fkxkikyTO3zO3efIkrKjfRx2CwFwO2rQ/3T04GQB/Lgr5lfJQU
iofvvVYuj2xBJao+3t9Ir0OeSbw1T5Rz03VLtN8SZhvaxWaBfwkUuUNL1glJO+Yd
LkMxGS0zAgMBAAECggEBAKM6m7RQUPlJE8u8qfOCDdSSKbIefrT9wZ5tKN0dG2Oa
/TNkzrEhXOO8F5Ek0a7LA+Q51KL7ksNtpLS0XpZNoYS8bapS36ePIJN0yx8nIJwc
koYlGtu/+U6ZpHQSoTiBjwRtswcudXuxT8i8frOupnWbFpKJ7H9Vbcb9bHB8N6Mm
D63wSBR08ZMrZXheKHQCQcxSQ2ZQZ+X3LBIOdXZH1aaptU2KpMEU5oyxXPShTVMg
0f748yU2njXCF0ZABEanXgp13egr/MPqHwnS/h0PH45bNy3IgFtMEHEouQFsAzoS
qNe8/9WnrpY87UdSZMnzF/IAXV0bmollDnqfM8/EqxkCgYEA96ThXYGzAK5RKNqp
RqVdRVA0UTT48sJvrxLMuHpyUzg6cl8FZE5rrNxFbouxvyN192Ctv1q8yfv4/HfM
KpmtEjt3fYtITHVXII6O3qNaRoIEPwKT4eK/ar+JO59vI0YvweXvDH5TkS9aiFr+
pPGf3a7EbE24BKhgiI8eT6K0VuUCgYEA5QGg11ZVoUut4ERAPouwuSdWwNe0HYqJ
A1m5vTvF5ghUHAb023lrr7Psq9DPJQQe7GzPfXafsat9hGenyqiyxo1gwClIyoEH
fOg753kdHcy60VVzumsPXece3OOSnd0rRMgfsSsclgYO7z0g9YZPAjt2w9NVw6uN
UDqX3eO2WjcCgYEA015eoNHv99fRG96udsbz+hI/5UQibAl7C+Iu7BJO/CrU8AOc
dYXdr5f+hyEioDLjIDbbdaU71+aCGPMjRwUNzK8HCRfVqLTKndYvqWWhyuZ0O1e2
4ykHGlTLDCHD2Uaxwny/8VjteNEDI7kO+bfmLG9b5djcBNW2Nzh4tZ348OUCgYEA
vIrTppbhF1QciqkGj7govrgBt/GfzDaTyZtkzcTZkSNIRG8Bx3S3UUh8UZUwBpTW
9OY9ClnQ7tF3HLzOq46q6cfaYTtcP8Vtqcv2DgRsEW3OXazSBChC1ZgEk+4Vdz1x
c0akuRP6jBXe099rNFno0LiudlmXoeqrBOPIxxnEt48CgYEAxNZBc/GKiHXz/ZRi
IZtRT5rRRof7TEiDxSKOXHSG7HhIRDCrpwn4Dfi+GWNHIwsIlom8FzZTSHAN6pqP
E8Imrlt3vuxnUE1UMkhDXrlhrxslRXU9enynVghAcSrg6ijs8KuN/9RB/I7H03cT
77mx9eHMcYcRUciY5C8AOaArmMA=
-----END PRIVATE KEY-----`
