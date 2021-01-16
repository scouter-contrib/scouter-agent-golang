package texttype

const (
	ERROR TextType = "error"
	APICALL TextType = "apicall"
	METHOD TextType = "method"
	SERVICE TextType = "service"
	SQL TextType = "sql"
	OBJECT TextType = "object"
	REFERER TextType = "referer"
	USER_AGENT TextType = "ua"
	GROUP TextType = "group"
	CITY TextType = "city"
	SQL_TABLES TextType = "table"
	MARIA TextType = "maria"
	LOGIN TextType = "login"
	DESC TextType = "desc"
	WEB TextType = "web"
	HASH_MSG TextType = "hmsg"
	STACK_ELEMENT TextType = "stackelem"
)

type TextType string
