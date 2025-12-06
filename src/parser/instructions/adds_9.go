package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADDS_9 struct {
                Xd uint32 `"ADDS" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
                Rwidth string `@("W"|"X")`
                Rm uint32 `@Integer`
}

type InstructionADDS_9 struct {
    Fields inst_ADDS_9 `@@`
}
func (i InstructionADDS_9) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionADDS_9) Xn() uint32 {
    return i.Fields.Xn
}
func (i InstructionADDS_9) R() string {
    return i.Fields.Rwidth
}
func (i InstructionADDS_9) M() uint32 {
    return i.Fields.Rm
}

func init() {
	var p = participle.MustBuild[InstructionADDS_9](
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

