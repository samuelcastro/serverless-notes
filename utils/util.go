package util

func GetResponseHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}
}

func GetUserId(headers map[string]string) string {
	return headers["app_user_id"]
}

func GetUserName(headers map[string]string) string {
	return headers["app_user_name"]
}
