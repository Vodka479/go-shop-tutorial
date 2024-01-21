package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// return เป็น struct ของ Config ออกไป **ฟังก์ชั่นใดๆ สามารถ return type เป็น interface ก็ได้
func LoadConfig(path string) IConfig { //install godotenv เพื่อให้โหลดไฟล์ .env
	envMap, err := godotenv.Read(path) // แปลง type env เป็น map
	if err != nil {                    // ดัก error
		log.Fatalf("load dotenv failed: %v", err) //ให้แอพทำงานพร้อม print err เป็น log
	}
	return &config{
		app: &app{
			host: envMap["APP_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["APP_PORT"]) // port เป็น int เลยต้องแปลง จาก string เป็น int
				if err != nil {
					log.Fatalf("load port failed: %v", err)
				}
				return p
			}(),
			name:    envMap["APP_NAME"],
			version: envMap["APP_VERSION"],
			bodyLimit: func() int {
				b, err := strconv.Atoi(envMap["APP_BODY_LIMIT"])
				if err != nil {
					log.Fatalf("load body limit failed: %v", err)
				}
				return b
			}(),
			readTimeout: func() time.Duration {
				t, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
				if err != nil {
					log.Fatalf("load read timeout failed: %v", err)
				}
				return time.Duration(int64(t) * int64(math.Pow10(9))) //ข้างในหน่วยเป็นวิ เลยต้องคูณกลับเป็นนาที
			}(),
			writeTimeout: func() time.Duration {
				t, err := strconv.Atoi(envMap["APP_WRITE_TIMEOUT"])
				if err != nil {
					log.Fatalf("load write timeout failed: %v", err)
				}
				return time.Duration(int64(t) * int64(math.Pow10(9))) //ข้างในหน่วยเป็นวิ เลยต้องคูณกลับเป็นนาที
			}(),
			fileLimit: func() int {
				f, err := strconv.Atoi(envMap["APP_FILE_LIMIT"])
				if err != nil {
					log.Fatalf("load file limit failed: %v", err)
				}
				return f
			}(),
			gcpbucket: envMap["APP_GCP_BUCKET"],
		},
		db: &db{
			host: envMap["DB_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["DB_PORT"])
				if err != nil {
					log.Fatalf("load port failed: %v", err)
				}
				return p
			}(),
			protocal: envMap["DB_PROTOCAL"],
			username: envMap["DB_USERNAME"],
			password: envMap["DB_PASSWORD"],
			database: envMap["DB_DATABASE"],
			sslMode:  envMap["DB_SSL_MODE"],
			maxConnections: func() int {
				m, err := strconv.Atoi(envMap["DB_MAX_CONNECTIONS"])
				if err != nil {
					log.Fatalf("load db max connections failed: %v", err)
				}
				return m
			}(),
		},
		jwt: &jwt{
			secrectKey: envMap["JWT_SECRET_key"],
			apiKey:     envMap["JWT_API_KEY"],
			adminKey:   envMap["JWT_ADMIN_KEY"],
			accessExpiresAt: func() int {
				t, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRES"])
				if err != nil {
					log.Fatalf("load jwt access expires failed: %v", err)
				}
				return t
			}(),
			refreshExpiresAt: func() int {
				t, err := strconv.Atoi(envMap["JWT_REFRESH_EXPIRES"])
				if err != nil {
					log.Fatalf("load jwt refresh expires failed: %v", err)
				}
				return t
			}(),
		},
	}
}

// สร้าง interface มาเพื่อกระจาย Encapsolution ให้ struct แต่ล่ะตัว **Encapsolution ปกปิดการทำงาน
type IConfig interface { // ใช้พิมพ์ใหญ่เพื่อให้ เป็น global สามารถเข้าถึงจากไฟล์อื่นได้
	App() IAppConfig
	Db() IDbConfig
	Jwt() IJwtConfig
}

type config struct {
	app *app
	db  *db /*private field ภายนอกเข้าถึงไม่ได้*/
	jwt *jwt
}

type IAppConfig interface { // เขียน function แยกทุกตัว
	URL() string // host:port รวมกัน
	Name() string
	Version() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	BodyLimit() int
	FileLimit() int
	Gcpbucket() string
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int //byte
	fileLimit    int //byte
	gcpbucket    string
}

func (c *config) App() IAppConfig { //สร้างเป็น pointer สามารถใช้ address ผ่านเข้าไปใน function ได้ไวกว่า struct **แล้วแต่จะใช้
	return c.app
}

// Url = host:port รวมกัน
func (a *app) URL() string                 { return fmt.Sprintf("%s:%d", a.host, a.port) } //ใช้ string format return กลับไป จะเป็นตาม pattern ที่เราพิมพ์ไป %s[string] %d[int] %v[any]
func (a *app) Name() string                { return a.name }
func (a *app) Version() string             { return a.version }
func (a *app) ReadTimeout() time.Duration  { return a.readTimeout }
func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }
func (a *app) BodyLimit() int              { return a.bodyLimit }
func (a *app) FileLimit() int              { return a.fileLimit }
func (a *app) Gcpbucket() string           { return a.gcpbucket }

type IDbConfig interface {
	Url() string
	MaxOpenConns() int
}

type db struct {
	host           string
	port           int
	protocal       string
	username       string
	password       string
	database       string
	sslMode        string
	maxConnections int
}

func (c *config) Db() IDbConfig {
	return c.db
}

func (d *db) Url() string { //url db ของ postgres
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.host,
		d.port,
		d.username,
		d.password,
		d.database,
		d.sslMode,
	)
}

func (d *db) MaxOpenConns() int { return d.maxConnections }

type IJwtConfig interface {
	SecrectKey() []byte
	AdminKey() []byte
	ApiKey() []byte
	AccessExpiresAt() int
	RefreshExpiresAt() int
	SetJwtAccessExpires(t int) //เผื่อ set วันหมดอายุ
	SetJwtRefreshExpires(t int)
}

type jwt struct {
	secrectKey       string
	adminKey         string
	apiKey           string
	accessExpiresAt  int
	refreshExpiresAt int
}

func (c *config) Jwt() IJwtConfig {
	return c.jwt
}

// imprement ตัว function
func (j *jwt) SecrectKey() []byte         { return []byte(j.secrectKey) }
func (j *jwt) AdminKey() []byte           { return []byte(j.adminKey) }
func (j *jwt) ApiKey() []byte             { return []byte(j.apiKey) }
func (j *jwt) AccessExpiresAt() int       { return j.accessExpiresAt }
func (j *jwt) RefreshExpiresAt() int      { return j.refreshExpiresAt }
func (j *jwt) SetJwtAccessExpires(t int)  { j.accessExpiresAt = t }  // set ค่า paremeter ที่ผ่านเข้ามา
func (j *jwt) SetJwtRefreshExpires(t int) { j.refreshExpiresAt = t } // set ค่า paremeter ที่ผ่านเข้ามา
