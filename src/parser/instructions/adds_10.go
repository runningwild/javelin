package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADDS_10 struct {
                Xd uint32 `"ADDS" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
                Imm uint32 `"#" @Integer`
}

type InstructionADDS_10 struct {
    Fields inst_ADDS_10 `@@`
}
func (i InstructionADDS_10) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionADDS_10) Xn() uint32 {
    return i.Fields.Xn
}
func (i InstructionADDS_10) Imm() uint32 {
    return i.Fields.Imm
}

func init() {
	var p = participle.MustBuild[InstructionADDS_10](
		participle.Lexer(asmDef),
	)
    parsers["ADDS"] = append(parsers["ADDS"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

