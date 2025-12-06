package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADRP_13 struct {
                Xd uint32 `"ADRP" "X" @Integer`
                Label string `@Ident`
}

type InstructionADRP_13 struct {
    Fields inst_ADRP_13 `@@`
}
func (i InstructionADRP_13) Xd() uint32 {
    return i.Fields.Xd
}

func init() {
	var p = participle.MustBuild[InstructionADRP_13](
		participle.Lexer(asmDef),
	)
    parsers["ADRP"] = append(parsers["ADRP"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

