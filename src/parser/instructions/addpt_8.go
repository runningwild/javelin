package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADDPT_8 struct {
                Xd uint32 `"ADDPT" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
                Xm uint32 `"," "X" @Integer`
      // LSL
      // imm:amount
      // <nil>
}

type InstructionADDPT_8 struct {
    Fields inst_ADDPT_8 `@@`
}
func (i InstructionADDPT_8) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionADDPT_8) Xn() uint32 {
    return i.Fields.Xn
}
func (i InstructionADDPT_8) Xm() uint32 {
    return i.Fields.Xm
}

func init() {
	var p = participle.MustBuild[InstructionADDPT_8](
		participle.Lexer(asmDef),
	)
    parsers["ADDPT"] = append(parsers["ADDPT"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

