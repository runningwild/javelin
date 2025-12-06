package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADCS_3 struct {
                Xd uint32 `"ADCS" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
                Xm uint32 `"," "X" @Integer`
}

type InstructionADCS_3 struct {
    Fields inst_ADCS_3 `@@`
}
func (i InstructionADCS_3) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionADCS_3) Xn() uint32 {
    return i.Fields.Xn
}
func (i InstructionADCS_3) Xm() uint32 {
    return i.Fields.Xm
}

func init() {
	var p = participle.MustBuild[InstructionADCS_3](
		participle.Lexer(asmDef),
	)
    parsers["ADCS"] = append(parsers["ADCS"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

