package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADD_5 struct {
                Xd uint32 `"ADD" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
                Imm uint32 `"#" @Integer`
}

type InstructionADD_5 struct {
    Fields inst_ADD_5 `@@`
}
func (i InstructionADD_5) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionADD_5) Xn() uint32 {
    return i.Fields.Xn
}
func (i InstructionADD_5) Imm() uint32 {
    return i.Fields.Imm
}

func init() {
	var p = participle.MustBuild[InstructionADD_5](
		participle.Lexer(asmDef),
	)
    parsers["ADD"] = append(parsers["ADD"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

