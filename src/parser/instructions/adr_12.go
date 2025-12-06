package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADR_12 struct {
                Xd uint32 `"ADR" "X" @Integer`
                Label string `@Ident`
}

type InstructionADR_12 struct {
    Fields inst_ADR_12 `@@`
}
func (i InstructionADR_12) Xd() uint32 {
    return i.Fields.Xd
}

func init() {
	var p = participle.MustBuild[InstructionADR_12](
		participle.Lexer(asmDef),
	)
    parsers["ADR"] = append(parsers["ADR"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

