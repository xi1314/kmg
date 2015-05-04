package kmgEmail

/*
从gmail发送邮件
	req := kmgEmail.NewSmtpRequestFromGmail("xxxx@gmail.com", "xxxx")
	req.To = []string{"xxxx@gmail.com"}
	req.Subject = "测试邮件标题2"
	req.Message = "<p>测试邮件内容2</p>"
	req.Socks4aProxyAddr = "127.0.0.1:20000"
	return kmgEmail.SmtpSendEmail(req)

*/
func NewSmtpRequestFromGmail(gmail string, password string) *SmtpRequest {
	return &SmtpRequest{
		SmtpHost:     "smtp.gmail.com",
		SmtpPort:     587,
		SmtpUsername: gmail,
		SmtpPassword: password,
		From:         gmail,
	}
}
