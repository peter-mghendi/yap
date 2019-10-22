package handler

import (
	"github.com/l3njo/yap/model"
	"github.com/mikespook/gorbac"
)

// RBAC is an instance of the Role-Based Access Control
var (
	RBAC                  *gorbac.RBAC
	permissionPostOps     gorbac.Permission
	permissionUserOps     gorbac.Permission
	permissionDraftOps    gorbac.Permission
	permissionReactionOps gorbac.Permission
)

// InitRBAC initializes the Role-Based Access Control
func InitRBAC() error {
	var err error
	rbac := gorbac.New()

	roleReader := gorbac.NewStdRole(string(model.UserReader))
	roleEditor := gorbac.NewStdRole(string(model.UserEditor))
	roleKeeper := gorbac.NewStdRole(string(model.UserKeeper))

	permissionPostOps = gorbac.NewStdPermission("postOps")         // Publish, Retract, Delete, Edit released posts | Delete reactions
	permissionUserOps = gorbac.NewStdPermission("userOps")         // Delete, Assign user
	permissionDraftOps = gorbac.NewStdPermission("draftOps")       // Create, Delete draft, Edit draft
	permissionReactionOps = gorbac.NewStdPermission("reactionOps") // Create, Delete reaction

	_ = roleKeeper.Assign(permissionPostOps)
	_ = roleKeeper.Assign(permissionUserOps)
	_ = roleEditor.Assign(permissionDraftOps)
	_ = roleReader.Assign(permissionReactionOps)

	_ = rbac.Add(roleReader)
	_ = rbac.Add(roleEditor)
	_ = rbac.Add(roleKeeper)

	err = rbac.SetParent(string(model.UserEditor), string(model.UserReader))
	err = rbac.SetParent(string(model.UserKeeper), string(model.UserEditor))
	if err != nil {
		return err
	}

	RBAC = rbac
	return err
}
