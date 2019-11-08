package pay

import (
	"errors"
)

//Transfer 转账到钱包
func (c *Client) Transfer(key string, params Params) (Params, error) {
	if c.isSandbox {
		return nil, errors.New("该方法不支持沙盒模式")
	}
	return sendTransfer(key, transferURL, params, true)
}

//TransferBank 转账到钱包
func (c *Client) TransferBank(key string, params Params) (Params, error) {
	if c.isSandbox {
		return nil, errors.New("该方法不支持沙盒模式")
	}
	return sendTransfer(key, transferBankURL, params, true)
}
