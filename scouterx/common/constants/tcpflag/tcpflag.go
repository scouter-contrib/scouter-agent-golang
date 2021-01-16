package tcpflag

const OK byte = 0x01
const NOT_OK byte = 0x02

const HasNEXT byte = 0x03
const NoNEXT byte = 0x04

const FAIL byte = 0x05
const INVALID_SESSION byte = 0x44

const CLUSTER_SEND_NEXT byte = 0x03
const CLUSTER_SEND_STOP byte = 0x04
