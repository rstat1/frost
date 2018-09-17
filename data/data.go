package data

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"time"

	"git.m/svcman/common"
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
		if _, err := tx.CreateBucketIfNotExists([]byte("Services")); err != nil {
			panic(err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("Users")); err != nil {
			panic(err)
		}
		return nil
	})
	ds.FrostInit()
	return ds
}

//FrostInit ...
func (data *DataStore) FrostInit() {
	if data.GetFirstRunState() {
		system := data.queryEngine.From("System")
		if err := system.Set("System", "id", common.RandomID(32)); err != nil {
			panic(err)
		}
		if err := system.Set("System", "key", common.RandomID(48)); err != nil {
			panic(err)
		}
	}
}

//NewUser ...
func (data *DataStore) NewUser(request AuthRequest, p []ServiceAuth) (User, error) {
	var user User
	UUID, _ := uuid.GenerateUUID()
	if user = data.GetUser(request.Username); user.Id != "" {
		return user, errors.New("user already exists")
	} else {
		usersBucket := data.queryEngine.From("Users")
		passHasher := sha512.New512_256()
		hash := passHasher.Sum([]byte(request.Password))
		user = User{
			Id:       UUID,
			PassHash: hex.EncodeToString(hash),
			Username: request.Username,
		}
		usersBucket.Save(&user)
		data.makeUserPermissionMap(request.Username, p)
	}
	return user, nil
}

//GetInstanceDetails ...
func (data *DataStore) GetInstanceDetails() (string, string) {
	var id string
	var key string

	system := data.queryEngine.From("System")
	system.Get("System", "id", &id)
	system.Get("System", "key", &key)

	return id, key
}

//GetUser ...
func (data *DataStore) GetUser(username string) User {
	var userInfo User
	usersBucket := data.queryEngine.From("Users")
	if err := usersBucket.One("Username", username, &userInfo); err == nil {
		return userInfo
	} else {
		return User{}
	}
}

//GetUserByID ...
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

//GetAllUserNames ...
func (data *DataStore) GetAllUserNames() ([]string, error) {
	var users []User
	var names []string
	usersBucket := data.queryEngine.From("Users")
	if err := usersBucket.All(&users); err == nil {
		for _, v := range users {
			names = append(names, v.Username)
		}
		return names, nil
	} else {
		return []string{}, err
	}
}

//GetServiceDetailss ...
func (data *DataStore) GetServiceDetailss() ([]ServiceDetails, error) {
	var err error
	var knownServices []ServiceDetails
	routeList := data.queryEngine.From("Services")
	err = routeList.All(&knownServices)

	if err == nil {
		return knownServices, nil
	} else {
		return nil, err
	}
}

//GetFirstRunState ...
func (data *DataStore) GetFirstRunState() bool {
	var firstRunState bool
	firstRun := data.queryEngine.From("System")
	if err := firstRun.Get("System", "firstrun", &firstRunState); err == nil {
		common.Logger.Debugln(err)
		return firstRunState
	} else {
		return true
	}
}

//GetRoute ...
func (data *DataStore) GetRoute(name string) (ServiceDetails, error) {
	var foundRoute ServiceDetails
	routes := data.queryEngine.From("Services")
	if err := routes.One("AppName", name, &foundRoute); err == nil {
		return foundRoute, nil
	} else {
		common.CreateFailureResponseWithFields(err, 500, logrus.Fields{
			"func":  "GetRoute",
			"route": name,
		})
		return ServiceDetails{}, err
	}
}

//DoesUserHavePermission ...
func (data *DataStore) DoesUserHavePermission(username, service, permission string) bool {
	var serviceAccess map[string]map[string]bool

	common.Logger.WithFields(logrus.Fields{"app": service, "username": username}).Debugln("checking permission...")

	permMap := data.queryEngine.From("SitePermissionMappings")
	if err := permMap.Get("SitePermissionMappings", username, &serviceAccess); err == nil {
		return serviceAccess[service][permission]
	} else {
		common.Logger.Errorln(err)
		return false
	}
}

//GetServiceDetails ...
func (data *DataStore) GetServiceDetails(name string) (ServiceDetails, error) {
	var serviceDetails ServiceDetails
	c := data.queryEngine.From("Services")
	if err := c.One("ServiceName", name, &serviceDetails); err == nil {
		return serviceDetails, nil
	} else {
		common.CreateFailureResponse(err, "GetServiceDetails", 500)
		return ServiceDetails{}, err
	}
}

//GetServiceByID ...
func (data *DataStore) GetServiceByID(id string) (ServiceDetails, error) {
	var serviceDetails ServiceDetails
	c := data.queryEngine.From("Services")
	err := c.One("ServiceID", id, &serviceDetails)
	if err == nil {
		common.Logger.Debugln(serviceDetails.ServiceID)

		return serviceDetails, nil
	} else {
		common.CreateFailureResponse(err, "GetServiceByID", 500)
		return ServiceDetails{}, err
	}
}

//SetFirstRunState ...
func (data *DataStore) SetFirstRunState() {
	firstRun := data.queryEngine.From("System")
	firstRun.Set("System", "firstrun", false)
}

//AddNewRoute ...
func (data *DataStore) AddNewRoute(service ServiceDetails) error {
	common.Logger.Debugln(service)
	if service.ServiceID == "" {
		service.ServiceID = common.RandomID(32)
		service.ServiceKey = common.RandomID(48)
	}
	if service.AppName == "" {
		return errors.New("no AppName set")
	}
	routes := data.queryEngine.From("Services")
	if err := routes.Save(&service); err != nil {
		common.Logger.WithField("func", "AddNewRoute").Errorln(err)
		return err
	} else {
		return nil
	}
}

//DeleteRoute ...
func (data *DataStore) DeleteRoute(route string) error {
	var foundRoute ServiceDetails
	routes := data.queryEngine.From("Services")
	if err := routes.One("AppName", route, &foundRoute); err == nil {
		return routes.DeleteStruct(&foundRoute)
	} else {
		return err
	}
	return nil
}

//DeleteUser ...
func (data *DataStore) DeleteUser(user string) error {
	usersBucket := data.queryEngine.From("Users")
	if user := data.GetUser(user); user.Username != "" {
		usersBucket.DeleteStruct(user)
		return nil
	} else {
		return errors.New("no such user")
	}
}

func (data *DataStore) makeUserPermissionMap(username string, p []ServiceAuth) {
	var serviceAccess map[string]map[string]bool
	serviceAccess = make(map[string]map[string]bool)
	for _, v := range p {
		permMap := make(map[string]bool)
		for _, v2 := range v.Permissions {
			permMap[v2.Name] = v2.Value
		}
		serviceAccess[v.Service] = permMap
	}
	permMap := data.queryEngine.From("SitePermissionMappings")
	permMap.Set("SitePermissionMappings", username, serviceAccess)
}
