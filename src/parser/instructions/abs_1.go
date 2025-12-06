package instruction

import "github.com/alecthomas/participle/v2"

type inst_ABS_1 struct {
                Xd uint32 `"ABS" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
}

type InstructionABS_1 struct {
    Fields inst_ABS_1 `@@`
}
func (i InstructionABS_1) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionABS_1) Xn() uint32 {
    return i.Fields.Xn
}

func init() {
	var p = participle.MustBuild[InstructionABS_1](
		participle.Lexer(asmDef),
	)
    parsers["ABS"] = append(parsers["ABS"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

