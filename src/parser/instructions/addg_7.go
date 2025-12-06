package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADDG_7 struct {
                Xd uint32 `"ADDG" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
                Uimm6 uint32 `"#" @Integer`
                Uimm4 uint32 `"#" @Integer`
}

type InstructionADDG_7 struct {
    Fields inst_ADDG_7 `@@`
}
func (i InstructionADDG_7) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionADDG_7) Xn() uint32 {
    return i.Fields.Xn
}
func (i InstructionADDG_7) Uimm6() uint32 {
    return i.Fields.Uimm6
}
func (i InstructionADDG_7) Uimm4() uint32 {
    return i.Fields.Uimm4
}

func init() {
	var p = participle.MustBuild[InstructionADDG_7](
		participle.Lexer(asmDef),
	)
    parsers["ADDG"] = append(parsers["ADDG"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

