package models

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/types"
)

const (
	HTTPHeaderPolicyStateEnabled  = 1 // 已启用
	HTTPHeaderPolicyStateDisabled = 0 // 已禁用
)

type HTTPHeaderPolicyDAO dbs.DAO

func NewHTTPHeaderPolicyDAO() *HTTPHeaderPolicyDAO {
	return dbs.NewDAO(&HTTPHeaderPolicyDAO{
		DAOObject: dbs.DAOObject{
			DB:     Tea.Env,
			Table:  "edgeHTTPHeaderPolicies",
			Model:  new(HTTPHeaderPolicy),
			PkName: "id",
		},
	}).(*HTTPHeaderPolicyDAO)
}

var SharedHTTPHeaderPolicyDAO *HTTPHeaderPolicyDAO

func init() {
	dbs.OnReady(func() {
		SharedHTTPHeaderPolicyDAO = NewHTTPHeaderPolicyDAO()
	})
}

// 初始化
func (this *HTTPHeaderPolicyDAO) Init() {
	this.DAOObject.Init()
	this.DAOObject.OnUpdate(func() error {
		return SharedSysEventDAO.CreateEvent(nil, NewServerChangeEvent())
	})
	this.DAOObject.OnInsert(func() error {
		return SharedSysEventDAO.CreateEvent(nil, NewServerChangeEvent())
	})
	this.DAOObject.OnDelete(func() error {
		return SharedSysEventDAO.CreateEvent(nil, NewServerChangeEvent())
	})
}

// 启用条目
func (this *HTTPHeaderPolicyDAO) EnableHTTPHeaderPolicy(tx *dbs.Tx, id int64) error {
	_, err := this.Query(tx).
		Pk(id).
		Set("state", HTTPHeaderPolicyStateEnabled).
		Update()
	return err
}

// 禁用条目
func (this *HTTPHeaderPolicyDAO) DisableHTTPHeaderPolicy(tx *dbs.Tx, id int64) error {
	_, err := this.Query(tx).
		Pk(id).
		Set("state", HTTPHeaderPolicyStateDisabled).
		Update()
	return err
}

// 查找启用中的条目
func (this *HTTPHeaderPolicyDAO) FindEnabledHTTPHeaderPolicy(tx *dbs.Tx, id int64) (*HTTPHeaderPolicy, error) {
	result, err := this.Query(tx).
		Pk(id).
		Attr("state", HTTPHeaderPolicyStateEnabled).
		Find()
	if result == nil {
		return nil, err
	}
	return result.(*HTTPHeaderPolicy), err
}

// 创建策略
func (this *HTTPHeaderPolicyDAO) CreateHeaderPolicy(tx *dbs.Tx) (int64, error) {
	op := NewHTTPHeaderPolicyOperator()
	op.IsOn = true
	op.State = HTTPHeaderPolicyStateEnabled
	err := this.Save(tx, op)
	if err != nil {
		return 0, err
	}
	return types.Int64(op.Id), nil
}

// 修改AddHeaders
func (this *HTTPHeaderPolicyDAO) UpdateAddingHeaders(tx *dbs.Tx, policyId int64, headersJSON []byte) error {
	if policyId <= 0 {
		return errors.New("invalid policyId")
	}

	op := NewHTTPHeaderPolicyOperator()
	op.Id = policyId
	op.AddHeaders = headersJSON
	err := this.Save(tx, op)

	return err
}

// 修改SetHeaders
func (this *HTTPHeaderPolicyDAO) UpdateSettingHeaders(tx *dbs.Tx, policyId int64, headersJSON []byte) error {
	if policyId <= 0 {
		return errors.New("invalid policyId")
	}

	op := NewHTTPHeaderPolicyOperator()
	op.Id = policyId
	op.SetHeaders = headersJSON
	err := this.Save(tx, op)

	return err
}

// 修改ReplaceHeaders
func (this *HTTPHeaderPolicyDAO) UpdateReplacingHeaders(tx *dbs.Tx, policyId int64, headersJSON []byte) error {
	if policyId <= 0 {
		return errors.New("invalid policyId")
	}

	op := NewHTTPHeaderPolicyOperator()
	op.Id = policyId
	op.ReplaceHeaders = headersJSON
	err := this.Save(tx, op)

	return err
}

// 修改AddTrailers
func (this *HTTPHeaderPolicyDAO) UpdateAddingTrailers(tx *dbs.Tx, policyId int64, headersJSON []byte) error {
	if policyId <= 0 {
		return errors.New("invalid policyId")
	}

	op := NewHTTPHeaderPolicyOperator()
	op.Id = policyId
	op.AddTrailers = headersJSON
	err := this.Save(tx, op)

	return err
}

// 修改DeleteHeaders
func (this *HTTPHeaderPolicyDAO) UpdateDeletingHeaders(tx *dbs.Tx, policyId int64, headerNames []string) error {
	if policyId <= 0 {
		return errors.New("invalid policyId")
	}

	namesJSON, err := json.Marshal(headerNames)
	if err != nil {
		return err
	}

	op := NewHTTPHeaderPolicyOperator()
	op.Id = policyId
	op.DeleteHeaders = string(namesJSON)
	err = this.Save(tx, op)

	return err
}

// 组合配置
func (this *HTTPHeaderPolicyDAO) ComposeHeaderPolicyConfig(tx *dbs.Tx, headerPolicyId int64) (*shared.HTTPHeaderPolicy, error) {
	policy, err := this.FindEnabledHTTPHeaderPolicy(tx, headerPolicyId)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, nil
	}

	config := &shared.HTTPHeaderPolicy{}
	config.Id = int64(policy.Id)
	config.IsOn = policy.IsOn == 1

	// AddHeaders
	if len(policy.AddHeaders) > 0 {
		refs := []*shared.HTTPHeaderRef{}
		err = json.Unmarshal([]byte(policy.AddHeaders), &refs)
		if err != nil {
			return nil, err
		}
		if len(refs) > 0 {
			for _, ref := range refs {
				headerConfig, err := SharedHTTPHeaderDAO.ComposeHeaderConfig(tx, ref.HeaderId)
				if err != nil {
					return nil, err
				}
				config.AddHeaders = append(config.AddHeaders, headerConfig)
			}
		}
	}

	// AddTrailers
	if len(policy.AddTrailers) > 0 {
		refs := []*shared.HTTPHeaderRef{}
		err = json.Unmarshal([]byte(policy.AddTrailers), &refs)
		if err != nil {
			return nil, err
		}
		if len(refs) > 0 {
			resultRefs := []*shared.HTTPHeaderRef{}
			for _, ref := range refs {
				headerConfig, err := SharedHTTPHeaderDAO.ComposeHeaderConfig(tx, ref.HeaderId)
				if err != nil {
					return nil, err
				}
				if headerConfig == nil {
					continue
				}
				resultRefs = append(resultRefs, ref)
				config.AddTrailers = append(config.AddTrailers, headerConfig)
			}
			config.AddHeaderRefs = resultRefs
		}
	}

	// SetHeaders
	if len(policy.SetHeaders) > 0 {
		refs := []*shared.HTTPHeaderRef{}
		err = json.Unmarshal([]byte(policy.SetHeaders), &refs)
		if err != nil {
			return nil, err
		}
		if len(refs) > 0 {
			resultRefs := []*shared.HTTPHeaderRef{}
			for _, ref := range refs {
				headerConfig, err := SharedHTTPHeaderDAO.ComposeHeaderConfig(tx, ref.HeaderId)
				if err != nil {
					return nil, err
				}
				if headerConfig == nil {
					continue
				}
				resultRefs = append(resultRefs, ref)
				config.SetHeaders = append(config.SetHeaders, headerConfig)
			}
			config.SetHeaderRefs = resultRefs
		}
	}

	// ReplaceHeaders
	if len(policy.ReplaceHeaders) > 0 {
		refs := []*shared.HTTPHeaderRef{}
		err = json.Unmarshal([]byte(policy.ReplaceHeaders), &refs)
		if err != nil {
			return nil, err
		}
		if len(refs) > 0 {
			resultRefs := []*shared.HTTPHeaderRef{}
			for _, ref := range refs {
				headerConfig, err := SharedHTTPHeaderDAO.ComposeHeaderConfig(tx, ref.HeaderId)
				if err != nil {
					return nil, err
				}
				if headerConfig == nil {
					continue
				}
				resultRefs = append(resultRefs, ref)
				config.ReplaceHeaders = append(config.ReplaceHeaders, headerConfig)
			}
			config.ReplaceHeaderRefs = resultRefs
		}
	}

	// Delete Headers
	if len(policy.DeleteHeaders) > 0 {
		headers := []string{}
		err = json.Unmarshal([]byte(policy.DeleteHeaders), &headers)
		if err != nil {
			return nil, err
		}
		config.DeleteHeaders = headers
	}

	// Expires
	// TODO

	return config, nil
}
