package constant

var RESP_NIL []byte = []byte("$-1\r\n")
var RESP_OK []byte = []byte("+OK\r\n")
var RESP_ZERO []byte = []byte(":0\r\n")
var RESP_ONE []byte = []byte(":1\r\n")
var TTL_KEY_NOT_EXIST []byte = []byte(":-2\r\n")
var TTL_KEY_EXIST_NO_EXPIRE []byte = []byte(":-1\r\n")

const NO_EXPIRE int64 = -1

const OBJ_TYPE_STRING uint8 = 0
const OBJ_ENCODING_RAW uint8 = 0
const OBJ_ENCODING_INT uint8 = 1
const OBJ_ENCODING_EMBSTR uint8 = 2

const EngineStatus_WAITING = 1
const EngineStatus_BUSY = 2
const EngineStatus_SHUTTING_DOWN = 3
