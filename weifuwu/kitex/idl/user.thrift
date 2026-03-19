namespace go user

include "common.thrift"

struct RegisterReq {
    1:string username,
    2:string password,
}

struct RegisterResp {
    1:common.Resp resp,
}

service UserService {
    RegisterResp Register(1: RegisterReq req)
}