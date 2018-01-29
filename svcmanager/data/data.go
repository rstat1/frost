package data

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"time"

	"git.m/svcmanager/common"
	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
)

//DataStore ...
type DataStore struct {
	filename       string
	knownProjects  []string
	DB             *bolt.DB
	queryEngine    *storm.DB
	originalImages *bolt.DB
	Cache          *CacheService
}

//NewDataStoreInstance Do I really need to explain this one?
func NewDataStoreInstance(filename string) *DataStore {
	db, err := bolt.Open(filename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	query, _ := storm.Open(filename, storm.UseDB(db))
	ds := &DataStore{
		DB:          db,
		queryEngine: query,
		filename:    filename,
		Cache:       NewCacheService(),
	}
	db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("System")); err != nil {
			panic(err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("Routes")); err != nil {
			panic(err)
		}
		tx.CreateBucketIfNotExists([]byte("Users"))
		return nil
	})
	return ds
}

//NewUser ...
func (data *DataStore) NewUser(name, password string) (error, User) {
	var groupType string = ""
	var user User
	UUID, _ := uuid.GenerateUUID()
	if user = data.GetUser(name); user.Id != "" {
		return errors.New("user already exists"), user
	} else {
		usersBucket := data.queryEngine.From("Users")
		count, err := usersBucket.Count(&User{})
		if err != nil {
			return err, User{}
		}
		if count == 0 {
			groupType = "root"
		} else if name == "client" {
			groupType = "client"
		} else {
			groupType = "user"
		}
		passHasher := sha512.New512_256()
		hash := passHasher.Sum([]byte(password))
		user = User{
			Id:       UUID,
			PassHash: hex.EncodeToString(hash),
			Group:    groupType,
			Username: name,
		}
		usersBucket.Save(&user)

		data.DB.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(user.Id))
			return err
		})
	}
	return nil, user
}

//GetUser ...
func (data *DataStore) GetUser(username string) User {
	var userInfo User
	usersBucket := data.queryEngine.From("Users")
	if err := usersBucket.One("Username", username, &userInfo); err == nil {
		return userInfo
	} else {
		common.Logger.WithField("username", username).Errorln("couldn't find user with that name")
		return User{}
	}
}

//GetUserById ...
func (data *DataStore) GetUserByID(id string) User {
	var userInfo User
	usersBucket := data.queryEngine.From("Users")
	if err := usersBucket.One("Id", id, &userInfo); err == nil {
		return userInfo
	} else {
		common.Logger.WithField("userid", id).Errorln("couldn't find user with that id")
		return User{}
	}
}

//GetKnownRoutes ...
func (data *DataStore) GetKnownRoutes() ([]KnownRoute, error) {
	var err error
	var knownRoutes []KnownRoute
	routeList := data.queryEngine.From("Routes")
	err = routeList.All(&knownRoutes)

	if err == nil {
		return knownRoutes, nil
	} else {
		return nil, err
	}
}

//AddNewRoute ...
func (data *DataStore) AddNewRoute(service KnownRoute) error {
	common.Logger.Debugln(service)
	if service.AppName == "" {
		return errors.New("no AppName set")
	}
	routes := data.queryEngine.From("Routes")
	if err := routes.Save(&service); err != nil {
		common.Logger.WithField("func", "AddNewRoute").Errorln(err)
		return err
	} else {
		return nil
	}
}

//DeleteRoute ...
func (data *DataStore) DeleteRoute(route string) error {
	var foundRoute KnownRoute
	routes := data.queryEngine.From("Routes")
	if err := routes.One("AppName", route, &foundRoute); err == nil {
		return routes.DeleteStruct(&foundRoute)
	} else {
		return err
	}
	return nil
}

//GetRoute ...
func (data *DataStore) GetRoute(name string) KnownRoute {
	var foundRoute KnownRoute
	routes := data.queryEngine.From("Routes")
	if err := routes.One("AppName", name, &foundRoute); err == nil {
		return foundRoute
	} else {
		common.CreateFailureResponseWithFields(err, 500, logrus.Fields{
			"func":  "GetRoute",
			"route": name,
		})
		return KnownRoute{}
	}
}
