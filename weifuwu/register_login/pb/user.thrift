namespace go user

include "common.thrift"

struct RegisterReq{
1: string account,
2: string name,
3: string password,
}

struct RegisterRes{
1: common.Resp resp,
}

struct LoginReq{
1: string account,
2: string password,
}

struct LoginRes{
1: common.Resp req,
2: string password,
}

service RegisterService{
    RegisterRes  Register(1:RegisterReq req)
    LoginRes Login(1:LoginReq req)
}