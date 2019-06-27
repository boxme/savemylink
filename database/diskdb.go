package database

import (
	"database/sql"
	"fmt"
	"savemylink/models"
	article "savemylink/save"
	save "savemylink/save/models"
	"savemylink/services"

	_ "github.com/lib/pq"
)

const (
	USER_TABLE    = "users"
	ARTICLE_TABLE = "articles"
)

func GetDiskDb(config *DbConfig) Db {
	info := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DbName)

	db, err := sql.Open("postgres", info)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Disk Db is now connected")

	return nil
}

type DiskDb struct {
	userDb    *userDiskDb
	articleDb *articleDiskDb
}

func (diskDb *DiskDb) GetUserDb() services.UserDB {
	return diskDb.userDb
}

func (diskDb *DiskDb) GetArticleDb() article.ArticleDB {
	return diskDb.articleDb
}

func (DiskDb *DiskDb) Close() {}

type userDiskDb struct {
	db *sql.DB
}

func (userDb *userDiskDb) CreateUser(
	email, password, userToken string) (*models.User, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s (email, password, userToken) 
		VALUE ($1, $2, $3)
		RETURNING *`,
		USER_TABLE)
	row := userDb.db.QueryRow(query, email, password, userToken)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		panic(err)
	}

	return models.NewUser(id, email, password, userToken), nil
}

func (userDb *userDiskDb) LoginUser(email, userToken string) (*models.User, error) {
	return nil, nil
}

func (userDb *userDiskDb) GetById(id uint64) (*models.User, error) {
	return nil, nil
}

func (userDb *userDiskDb) GetByEmail(email string) (*models.User, error) {
	return nil, nil
}

func (userDb *userDiskDb) GetByUserToken(userToken string) (*models.User, error) {
	return nil, nil
}

func (userDb *userDiskDb) LogoutUser(userToken string) error {
	return nil
}

type articleDiskDb struct {
	db *sql.DB
}

func (articleDb *articleDiskDb) SaveArticle(
	url, title, content, image string) (*save.Article, error) {
	return nil, nil
}

func (articleDb *articleDiskDb) GetArticle(articleId uint64) (*save.Article, error) {
	return nil, nil
}

func (articleDb *articleDiskDb) DeleteArticle(articleId uint64) error {
	return nil
}
