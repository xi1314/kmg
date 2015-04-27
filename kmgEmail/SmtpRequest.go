package kmgEmail

type SmtpRequest struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUsername string //also as stmp username
	SmtpPassword string
	From         string
	To           []string
	Subject      string
	Message      string

	Socks4aProxyAddr string
}
