package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADDS_11 struct {
                Xd uint32 `"ADDS" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
                Xm uint32 `"," "X" @Integer`
}

type InstructionADDS_11 struct {
    Fields inst_ADDS_11 `@@`
}
func (i InstructionADDS_11) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionADDS_11) Xn() uint32 {
    return i.Fields.Xn
}
func (i InstructionADDS_11) Xm() uint32 {
    return i.Fields.Xm
}

func init() {
	var p = participle.MustBuild[InstructionADDS_11](
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

