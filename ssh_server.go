package kaginawa

// SSHServer represents a SSH server entry that managed by Kaginawa Server.
type SSHServer struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Key      string `json:"key"`
	Password string `json:"password"`
}
