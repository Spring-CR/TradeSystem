package enum

type CxlRejResponseTo string

const (
	CxlRejResponseTo_None    CxlRejResponseTo = ""
	CxlRejResponseTo_Cancel  CxlRejResponseTo = "1"
	CxlRejResponseTo_Replace CxlRejResponseTo = "2"
)
