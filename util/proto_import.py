# Imports MCPE protocol field informations from PocketMine-MP source code
# This a just concept scratch code, so the code could be dirty or inefficient.

import sys, pid_import, traceback, os
from pathlib import Path

if len(sys.argv) < 2:
    print("Usage: %s <protocol directory>" % sys.argv[0])
    sys.exit()

phpType = {
    "Byte": "byte",
    "Short": "uint16",
    "Int": "uint32",
    "Long": "uint64",
    "Float": "float32",
    "Double": "float64",
    "String": "string",
}

batch_hardcoded = """type Batch struct {
	Payloads [][]byte
}

func (i Batch) Pid() byte { return BatchHead } // 0x92

func (i *Batch) Read(buf *bytes.Buffer) {
	i.Payloads = make([][]byte, 0)
	payload, err := util.DecodeDeflate(buf.Read(uint32(buffer.ReadInt(buf, ))))
	if err != nil {
		log.Println("Error while decompressing Batch payload:", err)
		return
	}
	b := bytes.NewBuffer(payload)
	for b.Require(4) {
		size := b.ReadInt()
		pk := b.Read(size)
		if pk[0] == 0x92 {
			panic("Invalid BatchPacket inside BatchPacket")
		}
		i.Payloads = append(i.Payloads, pk)
	}
}

func (i Batch) Write() *bytes.Buffer {
	b := new(bytes.Buffer)
	for _, pk := range i.Payloads {
		b.WriteInt(uint32(len(pk)))
		b.Write(pk)
	}
	payload := util.EncodeDeflate(b.Done())
	buf := new(bytes.Buffer)
	buf.BatchWrite(uint32(len(payload)), payload)
	return buf
}

"""

class Chainable:
    def __init__(self, data):
        self.data = data

    def get_class_code(self):
        return Chainable(self.data[self.data.find("extends DataPacket") + 20 : -2])

    def get_consts(self):
        ret = list()
        for line in self.data.split("\n"):
            line = line.replace("    ", "").replace("\n", "").replace("\t", "").replace(";", "")
            if len(line) < 6:
                continue
            if line[:5] == "const":
                if line[6:16] == "NETWORK_ID":
                    continue
                sp = line[6:].split(" = ")
                if sp[1].isdigit():
                    ret.append((sp[0], int(sp[1])))
        return ret

    def get_encode_func(self):
        start = self.data.find("encode()")
        data = self.data[start+8:]
        if data[2:match_brace(data, ("{", "}"))-2].replace("    ", "").replace("\n", "").replace("\t", "") == "":
            return None
        return list(map(
            lambda x: x.replace("    ", "").replace("\n", "").replace("\t", "").replace(";", ""),
            data[2:match_brace(data, ("{", "}"))-2].split("\n")
        ))

    def get_decode_func(self):
        start = self.data.find("decode()")
        data = self.data[start+8:]
        if data[2:match_brace(data, ("{", "}"))-2].replace("    ", "").replace("\n", "").replace("\t", "") == "":
            return None
        return list(map(
            lambda x: x.replace("    ", "").replace("\n", "").replace("\t", "").replace(";", ""),
            data[2:match_brace(data, ("{", "}"))-2].split("\n")
        ))


constTabs = pid_import.parse_consts()

def get_targets():
    path = Path(sys.argv[1])
    out = list()

    for _, (name, _) in enumerate(constTabs):
        yield (name[:-4], str((path / (name[:-4] + "Packet.php")).resolve()))

def match_brace(code: str, matches:tuple):
    lv = 0
    for i, letter in enumerate(code):
        if letter == matches[0]:
            lv += 1
        if letter == matches[1]:
            lv -= 1
            if lv == 0:
                return i

def read(path: str):
    f = open(path, "r")
    data = f.read()
    f.close()
    return data

def parse_method_invokes(expr):
    # if expr[:expr.find("(")] == "":
    #     expr = "Bytes" + expr
    return expr[:expr.find("(")], expr[expr.find("(")+1:match_brace(expr, ("(", ")"))]

def parse_coder(codelist):
    codes = list()
    for line, code in enumerate(codelist):
        if code == "$this->reset()" or code == "":
            continue
        if code[:7] != "$this->" or (code[7:10] != "get" and code[7:10] != "put"):
            codes.append(["Unknown", code])

        type_, arg = parse_method_invokes(code[10:])
        if arg[:7] != "$this->":
            itype, iarg = parse_method_invokes(arg)
            if itype == "strlen" and iarg[:7] == "$this->":
                type_ = "LengthOf_" + type_
                arg = iarg[7:]
            else:
                codes.append([line, code])
        else:
            arg = arg[7:]

        codes.append([type_, arg])

    return codes

targets = list(get_targets())
def get_go_code(n: int):
    if targets[n][1][:5] == "Batch":
        f.write(batch_hardcoded)
        return
    out = ""
    clscode = Chainable(read(targets[n][1])).get_class_code()
    consts = list(map(
        lambda x: (
            ''.join(map(
                lambda y: y[0].upper() + y[1:].lower(),
                x[0].split("_")
            )),
            x[1]
        )
        , clscode.get_consts()
    ))
    encode = clscode.get_encode_func()
    decode = clscode.get_encode_func()
    fields = list()
    if encode != None:
        fields = parse_coder(encode)
    elif decode != None:
        fields = parse_coder(decode)
    if len(consts) > 0:
        out += "const (\n"
        for const in consts:
            out += "    %s byte = %d\n" % const
        out += ")\n\n"

    if len(fields) > 0:
        out += "type %s struct {\n" % targets[n][0]
        for field in fields:
            field[1] = field[1][0].upper() + field[1][1:]
            field[1] = ''.join(filter(lambda x: x.isalpha(), field[1]))
            if field[0] in phpType:
                out +="    %s %s\n" % (field[1], phpType[field[0]])
        out += "}\n\n"
    else:
        out += "type %s struct{}\n" % targets[n][0]

    out += "// Pid implements Packet interface.\nfunc (i %s) Pid() byte { return %sHead }\n\n" % (targets[n][0], targets[n][0])

    if len(fields) > 0:
        out += "// Read implements Packet interface.\nfunc (i *%s) Read(buf *bytes.Buffer){\n" % targets[n][0]
        for i, field in enumerate(fields):
            if field[0] in phpType:
                out += "    i.%s = buf.Read%s()\n" % (field[1], field[0])
            elif field[0][:9] == "LengthOf_":
                out += "    i.%s = buf.Read(buf.Read%s())\n"%  (field[1], field[0][9:])
            elif field[0] == "" and fields[i-1][0][:9] == "LengthOf_" and fields[i-1][1] == field[1]:
                continue
            else:
                out += "    // Unexpected code:" + ' '.join(field) + "\n"
        out += "}\n\n"
    else:
        out += "// Read implements Packet interface.\nfunc (i *%s) Read(buf *bytes.Buffer){}\n" % targets[n][0]

    out += "// Write implements Packet interface.\nfunc (i %s) Write() *bytes.Buffer{\n    buf := new(bytes.Buffer)\n" % targets[n][0]
    for i, field in enumerate(fields):
        if field[0] in phpType:
            out += "    buf.Write%s(i.%s)\n" % tuple(field)
        elif field[0][:9] == "LengthOf_":
            out += "    buf.Write%s(len(i.%s))\n"%  (field[0][9:], field[1])
        elif field[0] == "" and fields[i-1][0][:9] == "LengthOf_" and fields[i-1][1] == field[1]:
            continue
        else:
            out += "    // Unexpected code:" + ' '.join(field) + "\n"
    out += "    return buf\n}\n\n"
    f.write(out)
if __name__ == "__main__":
    global f
    f = open("out.go", "w")
    f.write("""package lav7

import (
	"encoding/hex"

	"github.com/L7-MCPE/lav7/raknet"
	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/buffer"
)

""")
    f.write(pid_import.format_consts(pid_import.parse_consts()) + "\n\n")
    f.write("var packets = map[byte]Packet{\n")
    print("hh")
    for target in targets:
        f.write("    %sHead: new(%s),\n" % (target[0], target[0]))

    f.write("}\n\n")

    for i in range(len(targets)):
        try:
            get_go_code(i)
        except Exception as e:
            f.write('\n'.join(map(
    lambda x: "// " + x,
    ("""An exception was thrown while parsing/converting PocketMine-MP protocol.
Please read original PHP code and port it manually.
Exception: %s""" % (str(sys.exc_info()[1]) + "\n\n") +
Chainable(read(targets[i][1])).get_class_code().data).split("\n")
)) + "\n\n")
            print(targets[i][0]+"Packet:", e)
            traceback.print_exc()
            print()
            continue

    fmt = os.popen("gofmt -s -e out.go").read()
    f.seek(0)
    f.truncate()
    f.write(fmt)
    f.close()
