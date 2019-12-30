package wapi

import (
	"fmt"

	"github.com/x-io/wechat/util"
)

const (
	idcardURL         = "https://api.weixin.qq.com/cv/ocr/idcard?img_url=%s&access_token=%s"
	bankcardURL       = "https://api.weixin.qq.com/cv/ocr/bankcard?img_url=%s&access_token=%s"
	bizLicenseURL     = "https://api.weixin.qq.com/cv/ocr/bizlicense?img_url=%s&access_token=%s"
	printedTextURL    = "https://api.weixin.qq.com/cv/ocr/comm??img_url=%s&access_token=%s"
	vehicleLicenseURL = "https://api.weixin.qq.com/cv/ocr/driving?img_url=%s&access_token=%s"
	drivingLicenseURL = "https://api.weixin.qq.com/cv/ocr/drivinglicense?img_url=%s&access_token=%s"
)

//OCRByIDCard OCR 身份证
func OCRByIDCard(accessToken, url string) (response []byte, err error) {
	uri := fmt.Sprintf(idcardURL, url, accessToken)
	return util.HTTPPost(uri, nil)
}

//OCRByBankCard OCR 银行卡
func OCRByBankCard(accessToken, url string) (response []byte, err error) {
	uri := fmt.Sprintf(bankcardURL, url, accessToken)
	return util.HTTPPost(uri, nil)
}

//OCRByBizLicense OCR 营业执照
func OCRByBizLicense(accessToken, url string) (response []byte, err error) {
	uri := fmt.Sprintf(bizLicenseURL, url, accessToken)
	return util.HTTPPost(uri, nil)
}

//OCRByDrivingLicense OCR 驾驶证
func OCRByDrivingLicense(accessToken, url string) (response []byte, err error) {
	uri := fmt.Sprintf(drivingLicenseURL, url, accessToken)
	return util.HTTPPost(uri, nil)
}

//OCRByVehicleLicense OCR 行驶证
func OCRByVehicleLicense(accessToken, url string) (response []byte, err error) {
	uri := fmt.Sprintf(vehicleLicenseURL, url, accessToken)
	return util.HTTPPost(uri, nil)
}

//OCRByPrintedText OCR 文字
func OCRByPrintedText(accessToken, url string) (response []byte, err error) {
	uri := fmt.Sprintf(printedTextURL, url, accessToken)
	return util.HTTPPost(uri, nil)
}
