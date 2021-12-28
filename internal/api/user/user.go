package user

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	rpc "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//
//func GetUsersOnlineStatus(c *gin.Context) {
//	params := api.GetUsersOnlineStatusReq{}
//	if err := c.BindJSON(&params); err != nil {
//		log.NewError(params.OperationID, "bind json failed ", err.Error(), c)
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
//		return
//	}
//
//	if params.Secret != config.Config.Secret {
//		log.NewError(params.OperationID, "parse token failed ", params.Secret, config.Config.Secret)
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "secret failed"})
//		return
//	}
//
//	req := &pbRelay.GetUsersOnlineStatusReq{
//		OperationID: params.OperationID,
//		UserIDList:  params.UserIDList,
//	}
//	var wsResult []*rpc.GetUsersOnlineStatusResp_SuccessResult
//	var respResult []*rpc.GetUsersOnlineStatusResp_SuccessResult
//	flag := false
//	log.NewDebug(params.OperationID, "GetUsersOnlineStatus req come here", params.UserIDList)
//
//	grpcCons := getcdv3.GetConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOnlineMessageRelayName)
//	for _, v := range grpcCons {
//		client := rpc.NewOnlineMessageRelayServiceClient(v)
//		reply, err := client.GetUsersOnlineStatus(context.Background(), req)
//		if err != nil {
//			log.NewError(params.OperationID, "GetUsersOnlineStatus rpc  err", req.String(), err.Error())
//			continue
//		} else {
//			if reply.ErrCode == 0 {
//				wsResult = append(wsResult, reply.SuccessResult...)
//			}
//		}
//	}
//	log.NewDebug(params.OperationID, "call GetUsersOnlineStatus rpc server is success", wsResult)
//	//Online data merge of each node
//	for _, v1 := range params.UserIDList {
//		flag = false
//		temp := new(pbRelay.GetUsersOnlineStatusResp_SuccessResult)
//		for _, v2 := range wsResult {
//			if v2.UserID == v1 {
//				flag = true
//				temp.UserID = v1
//				temp.Status = constant.OnlineStatus
//				temp.DetailPlatformStatus = append(temp.DetailPlatformStatus, v2.DetailPlatformStatus...)
//			}
//		}
//		if !flag {
//			temp.UserID = v1
//			temp.Status = constant.OfflineStatus
//		}
//		respResult = append(respResult, temp)
//	}
//	log.NewDebug(params.OperationID, "Finished merged data", respResult)
//	resp := gin.H{"errCode": 0, "errMsg": "", "data": respResult}
//
//	c.JSON(http.StatusOK, resp)
//}

func GetUserInfo(c *gin.Context) {
	params := api.GetUserInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetUserInfoReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(params.OperationID, "GetUserInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	resp := api.GetUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserInfoList: RpcResp.UserInfoList}
	log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func UpdateUserInfo(c *gin.Context) {
	params := api.UpdateUserInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.UpdateUserInfoReq{UserInfo: &open_im_sdk.UserInfo{}}
	utils.CopyStructFields(req.UserInfo, &params)
	req.OperationID = params.OperationID
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(params.OperationID, "UpdateUserInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.UpdateUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "UpdateUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.UpdateUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "UpdateUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}