package binchunk

// brew install lua@5.3
// brew link --overwrite lua@5.3
// xxd luac.out

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type binaryChunk struct {
	header
	sizeUpvalues byte // ?
	mainFunc     *Prototype
}

type header struct {
	signature       [4]byte
	version         byte
	format          byte
	luacData        [6]byte
	cintSize        byte
	sizetSize       byte
	instrcutionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luacInt         int64   //lua 整型
	luacNum         float64 //lua 浮点型
}

// function prototype
type Prototype struct {
	Source          string // debug
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte     // 固定参数个数
	IsVararg        byte     // 是否有变长参数
	MaxStackSize    byte     // 寄存器数量
	Code            []uint32 // 指令表 每条指令占 4 个字节
	Constants       []interface{}
	Upvalues        []Upvalue    // 每个元素占用 2 个字节
	Protos          []*Prototype // 子函数原型表
	LineInfo        []uint32     // 行号表
	LocVars         []LocVar     // 局部变量表
	UpvalueNames    []string     // debug
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()
	reader.readByte() // size_upvalues
	return reader.readProto("")
}
