namespace go user

include "common.thrift"

struct RegisterReq{
1: i16 account,
2: string name,
3: i16 password,
}

struct RegisterRes{
1: common.Resp resp,
}

struct LoginReq{
1: i16 account,
2: i16 password,
}

struct LoginRes{
1: common.Resp req,
2: i16 password,
}

service RegisterService{
    RegisterRes  Register(1:RegisterReq req)
    LoginRes Login(1:LoginReq req)
}