package instruction

import "github.com/alecthomas/participle/v2"

type inst_MOVZ_254 struct {
                Xd uint32 `"MOVZ" "X" @Integer`
                Imm uint32 `"#" @Integer`
      // LSL
      // imm:shift
      // <nil>
}

type InstructionMOVZ_254 struct {
    Fields inst_MOVZ_254 `@@`
}
func (i InstructionMOVZ_254) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionMOVZ_254) Imm() uint32 {
    return i.Fields.Imm
}

func init() {
	var p = participle.MustBuild[InstructionMOVZ_254](
		participle.Lexer(asmDef),
	)
    parsers["MOVZ"] = append(parsers["MOVZ"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

