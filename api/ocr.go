package api

import "github.com/x-io/wechat/api/wapi"

//OCRByIDCard OCR 识别身份证
func (m *API) OCRByIDCard(key, url string) ([]byte, error) {
	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return nil, err
	}

	return wapi.OCRByIDCard(accessToken, url)
}

//OCRByBankCard OCR 银行卡
func (m *API) OCRByBankCard(key, url string) ([]byte, error) {
	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return nil, err
	}

	return wapi.OCRByBankCard(accessToken, url)
}

//OCRByBizLicense OCR 营业执照
func (m *API) OCRByBizLicense(key, url string) ([]byte, error) {
	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return nil, err
	}

	return wapi.OCRByBizLicense(accessToken, url)
}

//OCRByDrivingLicense OCR 驾驶证
func (m *API) OCRByDrivingLicense(key, url string) ([]byte, error) {
	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return nil, err
	}

	return wapi.OCRByDrivingLicense(accessToken, url)
}

//OCRByVehicleLicense OCR 行驶证
func (m *API) OCRByVehicleLicense(key, url string) ([]byte, error) {
	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return nil, err
	}

	return wapi.OCRByVehicleLicense(accessToken, url)
}

//OCRByPrintedText OCR 文字
func (m *API) OCRByPrintedText(key, url string) ([]byte, error) {
	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return nil, err
	}

	return wapi.OCRByPrintedText(accessToken, url)
}
