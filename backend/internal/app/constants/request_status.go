package constants

type RequestStatus string
type Decision string
type StepMode string

const (

	// Request / Flow Status
	RequestPending  RequestStatus = "pending"
	RequestApproved RequestStatus = "approved"
	RequestRejected RequestStatus = "rejected"
	RequestCanceled RequestStatus = "canceled"
	RequestInReview RequestStatus = "in_review"

	// Approval decision for single approval row
	DecisionPending  Decision = "pending"
	DecisionApproved Decision = "approved"
	DecisionRejected Decision = "rejected"

	// Step modes
	ModeAND StepMode = "AND"
	ModeOR  StepMode = "OR"
)
