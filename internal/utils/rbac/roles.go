package rbac

import (
	"database/sql"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

type RolesManager interface {
	Enforce(sub, obj, act string) (bool, error)
	UpdatePermissionsForRole(sub string, permissions [][]string) (bool, error)
	GetAllRole() ([]string, error)
	GetPermissionsForRole(sub string) ([][]string, error)
}

type RoleManager struct {
	db       *sql.DB
	enforcer *casbin.Enforcer
}

var _ RolesManager = (*RoleManager)(nil)

func NewRolesManager(db *sql.DB) *RoleManager {
	// Init casbin mysql adapters
	a, err := sqladapter.NewAdapter(db, "mysql", "rule_engine_casbin_rule")
	if err != nil {
		panic(err)
	}

	m := getCasbinModel()
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}
	return &RoleManager{
		db:       db,
		enforcer: e,
	}
}

func getCasbinModel() model.Model {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == \"1\"")

	// m = r.sub == p.sub && r.obj == p.obj && r.act == p.act || r.sub == "root"
	// g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act

	return m
}

func (r *RoleManager) Enforce(sub, obj, act string) (bool, error) {
	return r.enforcer.Enforce(sub, obj, act)
}

func (r *RoleManager) UpdatePermissionsForRole(sub string, permissions [][]string) (bool, error) {
	_, err := r.enforcer.DeletePermissionsForUser(sub)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(permissions); i++ {
		_, err := r.enforcer.AddPermissionsForUser(sub, permissions[i])
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (r *RoleManager) GetAllRole() ([]string, error) {
	return r.enforcer.GetAllSubjects()
}

func (r *RoleManager) GetPermissionsForRole(sub string) ([][]string, error) {
	return r.enforcer.GetPermissionsForUser(sub)
}
