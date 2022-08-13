package kookvoice

type gatewayResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    gatewayRespData `json:"data"`
}

type gatewayRespData struct {
	GatewayUrl string `json:"gateway_url"`
}

type firstShakeReq struct {
	Request bool              `json:"request"`
	Id      int               `json:"id"`
	Method  string            `json:"method"`
	Data    firstShakeReqData `json:"data"`
}

type firstShakeReqData struct {
}

type secondShakeReq struct {
	Request bool               `json:"request"`
	Id      int                `json:"id"`
	Method  string             `json:"method"`
	Data    secondShakeReqData `json:"data"`
}

type secondShakeReqData struct {
	DisplayName string `json:"displayName"`
}

type BaseShakeReq struct {
	Request bool   `json:"request"`
	Id      int    `json:"id"`
	Method  string `json:"method"`
}

type thirdShakeReq struct {
	Request bool              `json:"request"`
	Id      int               `json:"id"`
	Method  string            `json:"method"`
	Data    thirdShakeReqData `json:"data"`
}

type thirdShakeReqData struct {
	Comedia bool   `json:"comedia"`
	RtcpMux bool   `json:"rtcpMux"`
	Type    string `json:"type"`
}

type thirdShakeResp struct {
	Response bool `json:"response"`
	Id       int  `json:"id"`
	Ok       bool `json:"ok"`
	Data     thirdShakeRespData
}

type thirdShakeRespData struct {
	Id       string `json:"id"`
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	RtcpPort int    `json:"rtcpPort"`
}

type fourthShakeReq struct {
	Request bool               `json:"request"`
	Id      int                `json:"id"`
	Method  string             `json:"method"`
	Data    fourthShakeReqData `json:"data"`
}

type fourthShakeReqData struct {
	AppData       appData       `json:"appData"`
	Kind          string        `json:"kind"`
	PeerId        string        `json:"peerId"`
	RtpParameters rtpParameters `json:"rtpParameters"`
	TransportId   string        `json:"transportId"`
}

type appData struct {
}

type rtpParameters struct {
	Codecs    []codec    `json:"codecs"`
	Encodings []encoding `json:"encodings"`
}

type codec struct {
	Channels    int        `json:"channels"`
	ClockRate   int        `json:"clockRate"`
	MimeType    string     `json:"mimeType"`
	Parameters  parameters `json:"parameters"`
	PayloadType int        `json:"payloadType"`
}

type parameters struct {
	SpropStereo int `json:"sprop-stereo"`
}

type encoding struct {
	Ssrc int `json:"ssrc"`
}
