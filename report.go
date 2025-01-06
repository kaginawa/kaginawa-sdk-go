package kaginawa

import "time"

// Report represents a status report that enriched by Kaginawa Server.
type Report struct {
	// ID is the device identification string commonly hardware MAC address.
	ID string `json:"id"`

	// Trigger is the reason of report initiation.
	// -1: connected to the SSH server
	// 0: kaginawa started
	// 1 or higher: interval timer in minutes
	Trigger int `json:"trigger"`

	// Runtime is the runtime environment information such as OS name and CPU architecture.
	Runtime string `json:"runtime"`

	// Success is the shorthand of len(Errors) == 0.
	Success bool `json:"success"`

	// Sequence is the number of reports generated since process start.
	Sequence int `json:"seq"`

	// DeviceTime is the initiated UTC Unix timestamp in seconds on the device.
	DeviceTime int64 `json:"device_time"`

	// BootTime is the process started UTC Unix timestamp in seconds on the device.
	BootTime int64 `json:"boot_time"`

	// GenMillis is the report generation time in milliseconds.
	GenMillis int64 `json:"gen_ms"`

	// AgentVersion is the Kaginawa software version.
	AgentVersion string `json:"agent_version"`

	// CustomID is the user-specified device identification string.
	CustomID string `json:"custom_id"`

	// SSHServerHost is the hostname of the connected SSH server.
	SSHServerHost string `json:"ssh_server_host"`

	// SSHRemotePort is the port number of the connected SSH server.
	SSHRemotePort int `json:"ssh_remote_port"`

	// SSHConnectTime is the connected UTC Unix timestamp in seconds to the SSH server.
	SSHConnectTime int64 `json:"ssh_connect_time"`

	// Adapter is the name of network adapter, source of the MAC address.
	Adapter string `json:"adapter"`

	// LocalIPv4 is the local IP v4 address.
	LocalIPv4 string `json:"ip4_local"`

	// LocalIPv6 is the local IP v6 address.
	LocalIPv6 string `json:"ip6_local"`

	// Hostname is the hostname of the device.
	Hostname string `json:"hostname"`

	// RTTMillis is the measured round trip time in milliseconds.
	RTTMillis int64 `json:"rtt_ms"`

	// UploadKBPS is the measured upload throughput in kbps.
	UploadKBPS int64 `json:"upload_bps"`

	// DownloadKBPS is the measured download throughput in kbps.
	DownloadKBPS int64 `json:"download_bps"`

	// DiskTotalBytes is the total disk space in bytes.
	DiskTotalBytes int64 `json:"disk_total_bytes"`

	// DiskUsedBytes is the used disk space in bytes.
	DiskUsedBytes int64 `json:"disk_used_bytes"`

	// DiskLabel is the disk label.
	DiskLabel string `json:"disk_label"`

	// DiskFilesystem is the filesystem name of the disk.
	DiskFilesystem string `json:"disk_filesystem"`

	// DiskMountPoint is the mount point of the disk (default is root).
	DiskMountPoint string `json:"disk_mount_point"`

	// DiskDevice is the device name of the disk.
	DiskDevice string `json:"disk_device"`

	// USBDevices is the list of detected usb devices.
	USBDevices []USBDevice `json:"usb_devices"`

	// BDLocalDevices is the list of detected bluetooth devices.
	BDLocalDevices []string `json:"bd_local_devices"`

	// KernelVersion is the kernel/OS version.
	//
	// Available since:
	// - kaginawa v1.0.0
	// - kaginawa-server v0.0.3
	// - kaginawa-sdk-go v0.1.0
	KernelVersion string `json:"kernel_version"`

	// Errors is the list of report generation errors.
	Errors []string `json:"errors"`

	// Payload is the output of the PayloadCmd
	Payload string `json:"payload"`
	// PayloadCmd is the executed payload command
	PayloadCmd string `json:"payload_cmd"`

	// GlobalIP is the global IP address.
	// This is a server-side attribute.
	GlobalIP string `json:"ip_global"`

	// GlobalHost is the reverse-lookup result for global IP address.
	// This is a server-side attribute.
	GlobalHost string `json:"host_global"`

	// ServerTime is the server-consumed UTC Unix timestamp (seconds).
	// This is a server-side attribute.
	ServerTime int64 `json:"server_time"`
}

// USBDevice defines usb device attributes.
type USBDevice struct {
	Name      string `json:"name"`
	VendorID  string `json:"vendor_id"`
	ProductID string `json:"product_id"`
	Location  string `json:"location"`
}

// Timestamp returns Time object from ServerTime.
func (r Report) Timestamp() time.Time {
	return time.Unix(r.ServerTime, 0)
}

// BootTimestamp returns Time object from BootTime.
func (r Report) BootTimestamp() time.Time {
	return time.Unix(r.BootTime, 0)
}
