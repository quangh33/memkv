package constant

var RespNil []byte = []byte("$-1\r\n")
var RespOk []byte = []byte("+OK\r\n")
var RespZero []byte = []byte(":0\r\n")
var RespOne []byte = []byte(":1\r\n")
var TtlKeyNotExist []byte = []byte(":-2\r\n")
var TtlKeyExistNoExpire []byte = []byte(":-1\r\n")

const NoExpire int64 = -1

const ObjTypeString uint8 = 0
const ObjEncodingRaw uint8 = 0
const ObjEncodingInt uint8 = 1

const EngineStatusWaiting = 1
const EngineStatusBusy = 2
const EngineStatusShuttingDown = 3
