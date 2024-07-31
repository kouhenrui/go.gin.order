package dto

type Res struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type LogBody struct {
	SpendTime string      `json:"spend_time"`
	Path      string      `json:"path"`
	Method    string      `json:"method"`
	Status    int         `json:"status"`
	Proto     string      `json:"proto"`
	Ip        string      `json:"ip"`
	Body      string      `json:"body"`
	Query     string      `json:"query"`
	Message   interface{} `json:"message"`
}
type Captcha struct {
	Id      string `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
}

type TokenClaims struct {
	Id    uint   `json:"id,omitempty" `
	Name  string `json:"name,omitempty" binding:"-"`
	Phone string `json:"phone,omitempty" binding:"-"`
	Type  string `json:"type"`
	Email string `json:"email"`
	Role  int32  `json:"role"`
}
type TokenAndExp struct {
	Token   string `json:"token,omitempty"`
	ExpTime string `json:"exp_time,omitempty"`
}
type Cookie struct {
	Name   string
	Value  string
	MaxAge int
	Path   string
	Domain string
	Https  bool
	Http   bool
}
type User struct {
	Id       uint   `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Role     int32  `json:"role,omitempty"`
	Gender   string `json:"gender"`
	RoleName string `json:"role_name,omitempty"`
	Type     string `json:"type,omitempty"`
}

type CreateApprovalReq struct {
	Title       string `json:"title,omitempty" validate:"required" msg:"审批标题"`
	Description string `json:"description,omitempty" msg:"审批内容"`
	Status      string `json:"status,omitempty" validate:"required" msg:"审批状态"`
	Type        string `json:"type,omitempty" validate:"required" msg:"类型"`
	//ApproverId      uint   `json:"approver_id,omitempty" msg:"创建人id是uint类型"`
	CreateApprovers []CreateApprover `json:"create_approves"`
	Action          string           `json:"action,omitempty" msg:"审核日志行为（创建，审批，驳回）"`
	Comment         string           `json:"comment,omitempty" msg:"备注"`
}
type CreateApprover struct {
	Name  string `json:"name,omitempty"  msg:"审核员名称"`
	Email string `json:"email,omitempty" msg:"审核员邮箱"`
	Order int    `json:"order,omitempty" msg:"审核顺序"`
}

type UpdateApproval struct {
	ApprovalId uint   `json:"approval_id"`
	Action     string `json:"action"` // approve, reject
	Comment    string `json:"comment"`
}
