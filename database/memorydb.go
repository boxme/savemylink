package database

import (
	"errors"
	"savemylink/models"
	"savemylink/save"
	model "savemylink/save/models"
	"savemylink/services"
	"strings"
	"sync/atomic"
)

var id uint64

type MemoryDb struct {
	userMemoryDb    *UserMemoryDb
	articleMemoryDb *ArticleMemoryDb
}

func NewMemoryDb() *MemoryDb {
	return &MemoryDb{
		userMemoryDb: &UserMemoryDb{
			userTable: make(map[string]*models.User),
		},
		articleMemoryDb: &ArticleMemoryDb{
			id:           0,
			articleTable: make(map[string]*model.Article),
		},
	}
}

func (memoryDb *MemoryDb) GetUserDb() services.UserDB {
	return memoryDb.userMemoryDb
}

func (memoryDb *MemoryDb) GetArticleDb() save.ArticleDB {
	return memoryDb.articleMemoryDb
}

func (memoryDb *MemoryDb) Close() {
	memoryDb.userMemoryDb.close()
	memoryDb.articleMemoryDb.close()
}

type ArticleMemoryDb struct {
	id           uint64
	articleTable map[string]*model.Article
}

func (articleDb *ArticleMemoryDb) SaveArticle(
	url, title, content, image string) (*model.Article, error) {
	if _, ok := articleDb.articleTable[url]; ok {
		return nil, errors.New("article is already saved")
	}

	article :=
		model.NewArticle(
			atomic.AddUint64(&articleDb.id, 1), url, title, content, image)
	articleDb.articleTable[url] = article

	return article, nil
}

func (articleDb *ArticleMemoryDb) GetArticle(id uint64) (*model.Article, error) {
	for _, article := range articleDb.articleTable {
		if article.Id == id {
			return article, nil
		}
	}

	return nil, errors.New("Article not found")
}

func (articleDb *ArticleMemoryDb) DeleteArticle(id uint64) error {
	article, err := articleDb.GetArticle(id)
	if err == nil {
		delete(articleDb.articleTable, article.Url)
		return nil
	}

	return err
}

func (articleDb *ArticleMemoryDb) close() {
	articleDb.articleTable = make(map[string]*model.Article)
}

type UserMemoryDb struct {
	userTable map[string]*models.User
}

func (db *UserMemoryDb) CreateUser(email, password, userToken string) (*models.User, error) {
	if _, ok := db.userTable[email]; ok {
		return nil, errors.New(email + " is already registered")
	}

	user := models.NewUser(atomic.AddUint64(&id, 1), email, password, userToken)
	db.userTable[email] = user

	return user, nil
}

func (db *UserMemoryDb) GetById(id uint64) (*models.User, error) {
	for _, user := range db.userTable {
		if user.Id == id {
			return user, nil
		}
	}

	return nil, errors.New("User not found")
}

func (db *UserMemoryDb) GetByEmail(email string) (*models.User, error) {
	user, ok := db.userTable[email]
	if !ok {
		return nil, errors.New(email + " is not registered")
	}

	return user, nil
}

func (db *UserMemoryDb) GetByUserToken(userToken string) (*models.User, error) {
	for _, user := range db.userTable {
		if strings.Compare(userToken, user.Token) == 0 {
			return user, nil
		}
	}

	return nil, errors.New("User not found")
}

func (db *UserMemoryDb) LoginUser(email, userToken string) (*models.User, error) {
	user, ok := db.userTable[email]
	if !ok {
		return nil, errors.New(email + " is not registered")
	}

	updatedUser := models.NewUser(user.Id, user.Email, user.Password, userToken)
	db.userTable[email] = updatedUser

	return updatedUser, nil
}

func (db *UserMemoryDb) LogoutUser(userToken string) error {
	for _, user := range db.userTable {
		if strings.Compare(userToken, user.Token) == 0 {
			user.Token = ""
			return nil
		}
	}
	return errors.New("user is not found")
}

func (db *UserMemoryDb) close() {
	db.userTable = make(map[string]*models.User)
}
