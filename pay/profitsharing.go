package pay

//ProfitSharing 单次分账
func (c *Client) ProfitSharing(key string, params Params) (Params, error) {
	return sendAPI(key, ProfitSharingURL, params, true)
}

//ProfitSharingMulti 多次分账
func (c *Client) ProfitSharingMulti(key string, params Params) (Params, error) {
	return sendAPI(key, ProfitSharingMultiURL, params, true)
}

//ProfitSharingFinish 完结分账
func (c *Client) ProfitSharingFinish(key string, params Params) (Params, error) {
	return sendAPI(key, ProfitSharingFinishURL, params, true)
}

//ProfitSharingReturn 分账回退
func (c *Client) ProfitSharingReturn(key string, params Params) (Params, error) {
	return sendAPI(key, ProfitSharingReturnURL, params, true)
}

//ProfitSharingQuery 查询分账结果
func (c *Client) ProfitSharingQuery(key string, params Params) (Params, error) {
	return sendAPI(key, ProfitSharingQueryURL, params, false)
}

//ProfitSharinggReturQuery 查询分账回退结果
func (c *Client) ProfitSharinggReturQuery(key string, params Params) (Params, error) {
	return sendAPI(key, ProfitSharinggReturQueryURL, params, false)
}

//ProfitSharingAddReceiver 添加分账接收方
func (c *Client) ProfitSharingAddReceiver(key string, params Params) (Params, error) {
	return sendAPI(key, ProfitSharingAddReceiverURL, params, false)
}

//ProfitSharingRemoveReceiver 移除分账接收方
func (c *Client) ProfitSharingRemoveReceiver(key string, params Params) (Params, error) {
	return sendAPI(key, ProfitSharingRemoveReceiverURL, params, false)
}
