package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADD_4 struct {
                Xd uint32 `"ADD" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
                Rwidth string `@("W"|"X")`
                Rm uint32 `@Integer`
}

type InstructionADD_4 struct {
    Fields inst_ADD_4 `@@`
}
func (i InstructionADD_4) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionADD_4) Xn() uint32 {
    return i.Fields.Xn
}
func (i InstructionADD_4) R() string {
    return i.Fields.Rwidth
}
func (i InstructionADD_4) M() uint32 {
    return i.Fields.Rm
}

func init() {
	var p = participle.MustBuild[InstructionADD_4](
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

