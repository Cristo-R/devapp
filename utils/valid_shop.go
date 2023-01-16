package utils

import "regexp"

var DevEnvDomain = "preview.shoplazza.com"
var ProdEnvDomain = "myshoplaza.com"

func ValidShop(shop string) bool {
	devShopRegexp := regexp.MustCompile("^[a-zA-Z0-9-]+." + DevEnvDomain + "$")
	prodShopRegexp := regexp.MustCompile("^[a-zA-Z0-9-]+." + ProdEnvDomain + "$")
	if devShopRegexp.MatchString(shop) || prodShopRegexp.MatchString(shop) {
		return true
	}
	return false
}
