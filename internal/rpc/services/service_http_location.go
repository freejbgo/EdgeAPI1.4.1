package services

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAPI/internal/db/models"
	rpcutils "github.com/TeaOSLab/EdgeAPI/internal/rpc/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

type HTTPLocationService struct {
	BaseService
}

// 创建路径规则
func (this *HTTPLocationService) CreateHTTPLocation(ctx context.Context, req *pb.CreateHTTPLocationRequest) (*pb.CreateHTTPLocationResponse, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	locationId, err := models.SharedHTTPLocationDAO.CreateLocation(tx, req.ParentId, req.Name, req.Pattern, req.Description, req.IsBreak)
	if err != nil {
		return nil, err
	}

	return &pb.CreateHTTPLocationResponse{LocationId: locationId}, nil
}

// 修改路径规则
func (this *HTTPLocationService) UpdateHTTPLocation(ctx context.Context, req *pb.UpdateHTTPLocationRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	err = models.SharedHTTPLocationDAO.UpdateLocation(tx, req.LocationId, req.Name, req.Pattern, req.Description, req.IsOn, req.IsBreak)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// 查找路径规则配置
func (this *HTTPLocationService) FindEnabledHTTPLocationConfig(ctx context.Context, req *pb.FindEnabledHTTPLocationConfigRequest) (*pb.FindEnabledHTTPLocationConfigResponse, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	config, err := models.SharedHTTPLocationDAO.ComposeLocationConfig(tx, req.LocationId)
	if err != nil {
		return nil, err
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	return &pb.FindEnabledHTTPLocationConfigResponse{LocationJSON: configJSON}, nil
}

// 删除路径规则
func (this *HTTPLocationService) DeleteHTTPLocation(ctx context.Context, req *pb.DeleteHTTPLocationRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	err = models.SharedHTTPLocationDAO.DisableHTTPLocation(tx, req.LocationId)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// 查找反向代理设置
func (this *HTTPLocationService) FindAndInitHTTPLocationReverseProxyConfig(ctx context.Context, req *pb.FindAndInitHTTPLocationReverseProxyConfigRequest) (*pb.FindAndInitHTTPLocationReverseProxyConfigResponse, error) {
	// 校验请求
	adminId, userId, err := this.ValidateAdminAndUser(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	reverseProxyRef, err := models.SharedHTTPLocationDAO.FindLocationReverseProxy(tx, req.LocationId)
	if err != nil {
		return nil, err
	}
	if reverseProxyRef == nil || reverseProxyRef.ReverseProxyId <= 0 {
		reverseProxyId, err := models.SharedReverseProxyDAO.CreateReverseProxy(tx, adminId, userId, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		reverseProxyRef = &serverconfigs.ReverseProxyRef{
			IsOn:           false,
			ReverseProxyId: reverseProxyId,
		}
		reverseProxyJSON, err := json.Marshal(reverseProxyRef)
		if err != nil {
			return nil, err
		}
		err = models.SharedHTTPLocationDAO.UpdateLocationReverseProxy(tx, req.LocationId, reverseProxyJSON)
		if err != nil {
			return nil, err
		}
	}

	reverseProxyConfig, err := models.SharedReverseProxyDAO.ComposeReverseProxyConfig(tx, reverseProxyRef.ReverseProxyId)
	if err != nil {
		return nil, err
	}

	refJSON, err := json.Marshal(reverseProxyRef)
	if err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(reverseProxyConfig)
	if err != nil {
		return nil, err
	}
	return &pb.FindAndInitHTTPLocationReverseProxyConfigResponse{
		ReverseProxyJSON:    configJSON,
		ReverseProxyRefJSON: refJSON,
	}, nil
}

// 初始化Web设置
func (this *HTTPLocationService) FindAndInitHTTPLocationWebConfig(ctx context.Context, req *pb.FindAndInitHTTPLocationWebConfigRequest) (*pb.FindAndInitHTTPLocationWebConfigResponse, error) {
	// 校验请求
	adminId, userId, err := this.ValidateAdminAndUser(ctx, 0, 0)
	if err != nil {
		return nil, rpcutils.Wrap("ValidateRequest()", err)
	}

	tx := this.NullTx()

	webId, err := models.SharedHTTPLocationDAO.FindLocationWebId(tx, req.LocationId)
	if err != nil {
		return nil, rpcutils.Wrap("FindLocationWebId()", err)
	}

	if webId <= 0 {
		webId, err = models.SharedHTTPWebDAO.CreateWeb(tx, adminId, userId, nil)
		if err != nil {
			return nil, rpcutils.Wrap("CreateWeb()", err)
		}
		err = models.SharedHTTPLocationDAO.UpdateLocationWeb(tx, req.LocationId, webId)
		if err != nil {
			return nil, rpcutils.Wrap("UpdateLocationWeb()", err)
		}
	}

	config, err := models.SharedHTTPWebDAO.ComposeWebConfig(tx, webId)
	if err != nil {
		return nil, rpcutils.Wrap("ComposeWebConfig()", err)
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, rpcutils.Wrap("json.Marshal()", err)
	}
	return &pb.FindAndInitHTTPLocationWebConfigResponse{
		WebJSON: configJSON,
	}, nil
}

// 修改反向代理设置
func (this *HTTPLocationService) UpdateHTTPLocationReverseProxy(ctx context.Context, req *pb.UpdateHTTPLocationReverseProxyRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	err = models.SharedHTTPLocationDAO.UpdateLocationReverseProxy(tx, req.LocationId, req.ReverseProxyJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}
