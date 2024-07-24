package casbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	"github.com/casbin/gorm-adapter/v2"
	"go.gin.order/src/internal/pojo"
	"log"
)

type CasbinService struct {
	enforcer *casbin.Enforcer
}

func CasbinInit(casbinConfig *pojo.CabinConf) {
	db := casbinConfig.UserName + ":" + casbinConfig.PassWord + "@tcp(" + casbinConfig.HOST + ":" + casbinConfig.Port + ")/"
	//db加库名可以指定使用表或者自动创建表
	//a, aerr := gormadapter.NewAdapter(CasbinConfig.Type, db,true)//自己创建表
	adapter, aerr := gormadapter.NewAdapter(casbinConfig.Type, db)
	if aerr != nil {
		//fmt.Println("权限表为创建，错误原因", aerr)
		log.Fatal("权限表为创建，错误原因", aerr)
		return
	}
	enforcer, err := casbin.NewEnforcer("./auth_model.conf", adapter)
	//fmt.Println("e", e)

	if err != nil {
		fmt.Println("加载模型出现错误", err)
		return
	}
	_ = &CasbinService{enforcer: enforcer}
	log.Printf("权限初始化成功")
	//使用模糊匹配路径
	enforcer.AddFunction("regexMatch", RegexMatchFunc)
}

// 正则匹配函数
func RegexMatchFunc(args ...interface{}) (interface{}, error) {
	return util.RegexMatch(args[0].(string), args[1].(string)), nil
}
func (s *CasbinService) check(sub, obj, act string) {
	ok, _ := s.enforcer.Enforce(sub, obj, act)

	//fmt.Println(er, "err")
	if ok {
		fmt.Printf("%s CAN %s %s in %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s in %s\n", sub, act, obj)
	}
}
func (s *CasbinService) CheckPermission(sub, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, obj, act)
}

func (s *CasbinService) AddPolicy(sub, obj, act string) (bool, error) {
	return s.enforcer.AddPolicy(sub, obj, act)
}

func (s *CasbinService) RemovePolicy(sub, obj, act string) (bool, error) {
	return s.enforcer.RemovePolicy(sub, obj, act)
}

func (s *CasbinService) AddRoleForUser(user, role string) (bool, error) {
	return s.enforcer.AddRoleForUser(user, role)
}

func (s *CasbinService) DeleteRoleForUser(user, role string) (bool, error) {
	return s.enforcer.DeleteRoleForUser(user, role)
}
func (s *CasbinService) GetRolesForUser(user string) ([]string, error) {
	return s.enforcer.GetRolesForUser(user)
}

func (s *CasbinService) GetUsersForRole(role string) ([]string, error) {
	return s.enforcer.GetUsersForRole(role)
}

func (s *CasbinService) GetPermissionsForUser(user string) [][]string {
	return s.enforcer.GetPermissionsForUser(user)
}

func (s *CasbinService) HasRoleForUser(user, role string) (bool, error) {
	roles, err := s.GetRolesForUser(user)
	if err != nil {
		return false, err
	}
	for _, r := range roles {
		if r == role {
			return true, nil
		}
	}
	return false, nil
}
