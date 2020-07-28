package data

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"time"

	"go.alargerobot.dev/frost/common"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

//DataStore ...
type DataStore struct {
	filename       string
	knownProjects  []string
	DB             *bbolt.DB
	queryEngine    *storm.DB
	originalImages *bbolt.DB
	Cache          *CacheService
}

//NewDataStoreInstance Do I really need to explain this one?
func NewDataStoreInstance(filename string) *DataStore {
	db, err := bbolt.Open(filename, 0600, &bbolt.Options{Timeout: 1 * time.Second})
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
	db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("System")); err != nil {
			panic(err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("Services")); err != nil {
			panic(err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("Users")); err != nil {
			panic(err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("ExtraRoutes")); err != nil {
			panic(err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("Config")); err != nil {
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
		return firstRunState
	} else {
		common.Logger.Errorln(err)
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
		common.CreateFailureResponseWithFields(err, 500, logrus.Fields{"func": "GetRoute", "route": name})
		return ServiceDetails{}, err
	}
}

//DoesUserHavePermission ...
func (data *DataStore) DoesUserHavePermission(username, service, permission string) bool {
	var serviceAccess map[string]map[string]bool
	if username == "root" {
		return true
	} else {
		permMap := data.queryEngine.From("SitePermissionMappings")
		if err := permMap.Get("SitePermissionMappings", username, &serviceAccess); err == nil {
			return serviceAccess[service][permission]
		} else {
			common.Logger.Errorln(err)
			return false
		}
	}
}

//GetServiceByID ...
func (data *DataStore) GetServiceByID(id string) (ServiceDetails, error) {
	var serviceDetails ServiceDetails
	c := data.queryEngine.From("Services")
	err := c.One("ServiceID", id, &serviceDetails)
	if err == nil {
		return serviceDetails, nil
	} else {
		common.CreateFailureResponse(err, "GetServiceByID", 500)
		return ServiceDetails{}, err
	}
}

//GetAllExtraRoutes ...
func (data *DataStore) GetAllExtraRoutes() ([]ExtraRoute, error) {
	var extraRoutes []ExtraRoute
	extras := data.queryEngine.From("ExtraRoutes")
	if e2 := extras.All(&extraRoutes); e2 != nil {
		common.CreateFailureResponse(e2, "GetAllExtraRoutes", 500)
		return nil, e2
	}
	return extraRoutes, nil
}

//GetExtraRoutesForAPIName ...
func (data *DataStore) GetExtraRoutesForAPIName(name string) ([]ExtraRoute, error) {
	var extraRoutes []ExtraRoute
	extras := data.queryEngine.From("ExtraRoutes")
	if err := extras.Select(q.Eq("APIName", name)).Find(&extraRoutes); err != nil {
		common.CreateFailureResponse(err, "GetAllExtraRoutes", 500)
		return nil, err
	}
	return extraRoutes, nil
}

//GetUserPermissionMap ...
func (data *DataStore) GetUserPermissionMap(username string) (map[string]map[string]bool, error) {
	var serviceAccess map[string]map[string]bool
	permMap := data.queryEngine.From("SitePermissionMappings")
	if err := permMap.Get("SitePermissionMappings", username, &serviceAccess); err == nil {
		return serviceAccess, nil
	} else {
		return nil, err
	}
}

//GetServiceConfigValue ...
func (data *DataStore) GetServiceConfigValue(key, serviceName string) (value string, e error) {
	conf := data.queryEngine.From("Config")
	e = conf.Get(serviceName, key, &value)
	return value, common.LogError("", e)
}

//SetFirstRunState ...
func (data *DataStore) SetFirstRunState() {
	firstRun := data.queryEngine.From("System")
	firstRun.Set("System", "firstrun", false)
}

//SetConfigValue ...
func (data *DataStore) SetConfigValue(key, service string, value interface{}) error {
	e := data.DB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.Bucket([]byte("Config")).CreateBucketIfNotExists([]byte(service))
		return err
	})
	if e != nil {
		return common.LogError("", e)
	}
	conf := data.queryEngine.From("Config")
	return conf.Set(service, key, value)
}

//AddNewRoute ...
func (data *DataStore) AddNewRoute(service ServiceDetails) error {
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
	}
	return nil
}

//AddExtraRoute ...
func (data *DataStore) AddExtraRoute(newRoute ExtraRoute) error {
	if newRoute.APIName == "" || newRoute.FullURL == "" {
		return errors.New("missing some important info")
	}
	extras := data.queryEngine.From("ExtraRoutes")
	if err := extras.Save(&newRoute); err != nil {
		common.Logger.WithField("func", "AddExtraRoute").Errorln(err)
		return err
	}
	return nil
}

//DeleteRoute ...
func (data *DataStore) DeleteRoute(route string, deleteForUpdate bool) error {
	var foundRoute ServiceDetails
	routes := data.queryEngine.From("Services")
	extras := data.queryEngine.From("ExtraRoutes")

	if e := routes.One("AppName", route, &foundRoute); e != nil {
		common.Logger.WithField("func", "DeleteRoute(find)").Errorln(e)
		return e
	}
	if err := routes.DeleteStruct(&foundRoute); err != nil {
		common.Logger.WithField("func", "DeleteRoute(delete)").Errorln(err)
		return err
	}
	if deleteForUpdate == false {
		if e2 := extras.Select(q.Eq("AppName", route)).Delete(new(ExtraRoute)); e2 != nil {
			common.Logger.WithField("func", "DeleteRoute(delete)").Errorln(e2)
		}
	}
	return nil
}

//DeleteUser ...
func (data *DataStore) DeleteUser(user string) error {
	usersBucket := data.queryEngine.From("Users")
	if user := data.GetUser(user); user.Username != "" {
		return usersBucket.DeleteStruct(&user)
	} else {
		return errors.New("no such user")
	}
}

//DeleteExtraRoute ...
func (data *DataStore) DeleteExtraRoute(route string) error {
	var extraRoute []ExtraRoute
	common.Logger.Debugln(route)
	extras := data.queryEngine.From("ExtraRoutes")
	if err := extras.Select(q.Eq("APIRoute", route)).Find(&extraRoute); err == nil {
		return extras.DeleteStruct(&extraRoute[0])
	} else {
		return err
	}
}

//UpdateRoute ...
func (data *DataStore) UpdateRoute(service ServiceDetails, serviceName string) error {
	routes := data.queryEngine.From("Services")
	if e := data.DeleteRoute(serviceName, true); e != nil {
		common.Logger.WithFields(logrus.Fields{"func": "UpdateRoute", "action": "delete"}).Errorln(e)
		return e
	} else {
		return routes.Save(&service)
	}
}

//UpdateSysConfig ...
func (data *DataStore) UpdateSysConfig(propChange ServiceEdit) error {
	system := data.queryEngine.From("System")
	return system.Set("System", "key", &propChange.NewValue)
}

//UpdateUserPermissions ...
func (data *DataStore) UpdateUserPermissions(newPermission PermissionChange) error {
	if perms, err := data.GetUserPermissionMap(newPermission.Username); err == nil {
		servicePerms := perms[newPermission.ServiceName]
		if servicePerms != nil {
			servicePerms[newPermission.Name] = newPermission.NewValue
		} else {
			servicePerms = make(map[string]bool)
			servicePerms[newPermission.Name] = newPermission.NewValue
		}
		perms[newPermission.ServiceName] = servicePerms
		common.Logger.Debugln(perms)
		permMap := data.queryEngine.From("SitePermissionMappings")
		return permMap.Set("SitePermissionMappings", newPermission.Username, perms)
	} else {
		return err
	}
}

//UpdateUser ...
func (data *DataStore) UpdateUser(user User) error {
	usersBucket := data.queryEngine.From("Users")
	return usersBucket.Update(&user)
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
