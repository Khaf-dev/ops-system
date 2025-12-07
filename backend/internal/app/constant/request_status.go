package constants

type RequestStatus string

const (
	RequestPending  RequestStatus = "pending"
	RequestApproved RequestStatus = "approved"
	RequestRejected RequestStatus = "rejected"
	RequestCanceled RequestStatus = "canceled"
	RequestInReview RequestStatus = "in_review"
)
