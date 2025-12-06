package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADD_6 struct {
                Xd uint32 `"ADD" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
                Xm uint32 `"," "X" @Integer`
}

type InstructionADD_6 struct {
    Fields inst_ADD_6 `@@`
}
func (i InstructionADD_6) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionADD_6) Xn() uint32 {
    return i.Fields.Xn
}
func (i InstructionADD_6) Xm() uint32 {
    return i.Fields.Xm
}

func init() {
	var p = participle.MustBuild[InstructionADD_6](
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

