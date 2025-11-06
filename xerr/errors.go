package xerr

/**
错误码定义
A开头 客户端错误
B开头 服务端错误
C开头 三方服务错误

0-99      HTTP 协议定义的错误
100-999   框架定义的错误
1000-9999 用户自定义错误
*/

var (
	BadRequest          = New2(400, "A0001", "Bad Request")
	Unauthorized        = New2(401, "A0002", "Unauthorized")
	Forbidden           = New2(403, "A0003", "Forbidden")
	NotFound            = New2(404, "A0004", "Not Found")
	MethodNotAllowed    = New2(405, "A0005", "Method Not Allowed")
	InternalServerError = New2(500, "B0001", "Internal Server Error")

	ForbiddenIp     = New2(400, "A0100", "forbidden ip")
	ParametersError = New2(400, "A0101", "parameters error")
)
