package models

import (
	"regexp"
	"time"

	"github.com/0987363/mgo"
	"github.com/0987363/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

// RegexpUser is regexp with username
var RegexpUser = regexp.MustCompile(`^[A-Za-z0-9.+-_@]{3,30}$`)

// RegexpPwd is regexp with sha256 password
var RegexpPwd = regexp.MustCompile(`^[A-Za-z0-9]{1,64}$`)

type User struct {
	ID             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserName       string        `json:"username" bson:"username,omitempty"`
	Password       string        `json:"password,omitempty" bson:"-"`
	HashedPassword string        `json:"-" bson:"hashed_password,omitempty"`
	Token          string        `json:"-" bson:"token,omitempty"`
	Expiry         *time.Time    `json:"-" bson:"expiry,omitempty"`
}

func (user *User) UpdateToken(db *mgo.Session) error {
	c := db.DB(VSub).C(UserCollection)

	err := c.UpdateId(user.ID, bson.M{"$set": bson.M{"token": user.Token, "expiry": user.Expiry}})
	if err != nil {
		return err
	}
	return nil
}

func (user *User) Create(db *mgo.Session) error {
	c := db.DB(VSub).C(UserCollection)

	info, err := c.Upsert(bson.M{"username": user.UserName}, bson.M{"$setOnInsert": user})
	if err != nil {
		return err
	}
	if info.UpsertedId == nil {
		return ErrIsExist
	}

	return nil
}

// Update a user infomation
func (user *User) Update(db *mgo.Session) error {
	c := db.DB(VSub).C(UserCollection)

	err := c.UpdateId(user.ID, bson.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserByID(db *mgo.Session, id bson.ObjectId) error {
	c := db.DB(VSub).C(UserCollection)
	if err := c.RemoveId(id); err != nil {
		return err
	}
	return nil
}

func (user *User) HashPassword() bool {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return false
	}
	user.HashedPassword = string(hashedPassword)
	return true
}

// PasswordVerify verify original password with bcrypt password
func PasswordVerify(password, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

// FindUserByUserName find a user exists by username
func FindUserByUserName(db *mgo.Session, name string) (*User, error) {
	c := db.DB(VSub).C(UserCollection)

	var user User
	if err := c.Find(bson.M{"username": name}).One(&user); err != nil {
		return nil, ErrorConvert(err)
	}
	return &user, nil
}

// FindUserByID find a user exists by id
func FindUserByID(db *mgo.Session, id bson.ObjectId) *User {
	c := db.DB(VSub).C(UserCollection)

	user := User{}
	err := c.Find(bson.M{"_id": id}).One(&user)
	if err != nil {
		return nil
	}
	return &user
}
