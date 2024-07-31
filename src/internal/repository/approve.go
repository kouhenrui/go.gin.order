package repository

import (
	"go.gin.order/src/config/database"
	"go.gin.order/src/internal/pojo"
	"gorm.io/gorm"
)

type Approve struct {
	db          *gorm.DB
	approval    pojo.Approval
	approver    pojo.Approver
	approvation pojo.ApprovalAction
}

func NewApprovalRepository() *Approve {
	return &Approve{
		db:          database.PostgreClient,
		approval:    pojo.Approval{},
		approver:    pojo.Approver{},
		approvation: pojo.ApprovalAction{},
	}
}

func (a *Approve) CreateApproval(app pojo.Approval) (uint, error) {
	err := a.db.Save(&app).Error
	return app.ID, err
}
func (a *Approve) CreateApprover(apr pojo.Approver) (uint, error) {
	err := a.db.Save(&apr).Error
	return apr.ID, err
}
func (a *Approve) CreateApprovalAction(apa pojo.ApprovalAction) (uint, error) {
	err := a.db.Save(&apa).Error
	return apa.ID, err
}

func (a *Approve) UpdateApproval(id uint, status string) error {
	a.approval.ID = id
	a.approval.Status = status
	return a.db.Updates(&a.approval).Error
}

func (a *Approve) UpdateApprover(approval_id, approver_id uint, status, Comment string) error {
	a.approver.ApprovalID = approval_id
	a.approver.ID = approver_id
	a.approver.Status = status
	a.approver.Comment = Comment
	return a.db.Updates(&a.approver).Error
}
func (a *Approve) GetApproversByApprovalID(id uint) ([]ApproversByApprovalID, error) {
	var approvers []ApproversByApprovalID
	a.approver.ApprovalID = id
	err := a.db.Model(&a.approver).Scan(&approvers).Error
	return approvers, err
}

func (a *Approve) GetApproverOrderByApproverID(id uint) (*pojo.Approver, error) {
	a.approver.ID = id
	err := a.db.Find(&a.approver).Error
	return &a.approver, err
}
