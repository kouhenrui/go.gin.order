package service

import (
	"errors"
	"go.gin.order/pkg/msg"
	"go.gin.order/src/config/dto"
	"go.gin.order/src/internal/pojo"
	"go.gin.order/src/internal/repository"
	"log"
)

type ApprovalService struct {
	approveRep *repository.Approve
}

func NewApprovalService() *ApprovalService {
	return &ApprovalService{approveRep: repository.NewApprovalRepository()}
}
func (a *ApprovalService) CreateApprove(req *dto.CreateApprovalReq) error {
	log.Println(req, "req")
	status := "pending"
	var approval = pojo.Approval{
		Title:       req.Title,
		Description: req.Description,
		Status:      status, //req.Status,
		Type:        req.Type,
		CreatedBy:   "",
	}
	var approve_ids = []uint{}
	approve_id, err := a.approveRep.CreateApproval(approval)
	if err != nil {
		return err
	}
	log.Println(req.CreateApprovers, "approvers")
	//存储审核人
	for _, approvers := range req.CreateApprovers {
		var approver = pojo.Approver{
			ApprovalID: approve_id,
			Name:       approvers.Name,
			Email:      approvers.Email,
			Status:     status,
			Order:      approvers.Order,
		}
		appriver_id, err := a.approveRep.CreateApprover(approver)
		if err != nil {
			return err
		}
		approve_ids = append(approve_ids, appriver_id)
		//approver.ApprovalID = approve_id
		//approver.Name = approvers.Name
		//approver.Email = approvers.Email
		//approver.Order = approvers.Order
		//approver.Status = "pending"
	}
	//审核日志
	for _, approveids := range approve_ids {
		var approvalaction = pojo.ApprovalAction{
			ApprovalID: approve_id,
			ApproverID: approveids,
			Action:     status,
			Comment:    "",
		}
		_, err = a.approveRep.CreateApprovalAction(approvalaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ApprovalService) UpdateApprove(req dto.UpdateApproval, approver_id uint) error {
	approver, err := a.approveRep.GetApproverOrderByApproverID(approver_id)
	if err != nil {
		return err
	}
	order := approver.Order
	approvers, err := a.approveRep.GetApproversByApprovalID(req.ApprovalId)
	if err != nil {
		return err
	}
	for _, value := range approvers {
		if value.ApproverID != approver_id && value.Order < order && len(value.Status) < 0 {
			return errors.New(msg.APPROVALORDERERROR)
		}
	}
	err = a.approveRep.UpdateApprover(req.ApprovalId, approver_id, req.Action, req.Comment)
	if err != nil {
		return err
	}
	var approvalaction = pojo.ApprovalAction{
		ApprovalID: req.ApprovalId,
		ApproverID: approver_id,
		Action:     req.Action,
		Comment:    req.Comment,
	}
	_, err = a.approveRep.CreateApprovalAction(approvalaction)
	if err != nil {
		return err
	}
	return nil
}
