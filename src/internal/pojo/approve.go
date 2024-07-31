package pojo

import "gorm.io/gorm"

// Approval 表示申请信息
type Approval struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(255);not null;comment:'申请标题'"`
	Description string `json:"description" gorm:"type:text;comment:'申请描述'"`
	Status      string `json:"status" gorm:"type:approval_status;default:'pending';comment:'申请状态'"`
	Type        string `json:"type" gorm:"type:varchar(50);not null;comment:'申请类型，例如 leave, expense'"`
	CreatedBy   string `json:"created_by" gorm:"type:varchar(255);not null;comment:'创建者邮箱或用户名'"`
}

// Approver 表示审批人信息
type Approver struct {
	gorm.Model
	ApprovalID uint   `json:"approvalId" gorm:"not null;comment:'申请ID，关联到Approval表'"`
	Name       string `json:"name" gorm:"type:varchar(255);not null;comment:'审批人姓名'"`
	Email      string `json:"email" gorm:"type:varchar(255);not null;comment:'审批人邮箱'"`
	Status     string `json:"status" gorm:"type:approval_status;default:'pending';comment:'审批状态'"`
	Order      int    `json:"order" gorm:"type:int;not null;comment:'审批顺序'"`
	Comment    string `json:"comment" gorm:"type:text;comment:'审批人的附加评论'"`
}

// ApprovalAction 表示审批记录
type ApprovalAction struct {
	gorm.Model
	ApprovalID uint   `json:"approval_id" gorm:"not null;comment:'申请ID，关联到approval表'"`
	ApproverID uint   `json:"approver_id" gorm:"not null;comment:'审批人ID，关联到approver表'"`
	Action     string `json:"action" gorm:"type:approval_status;not null;comment:'审批动作'"`
	Comment    string `json:"comment" gorm:"type:text;comment:'审批附加评论'"`
}
