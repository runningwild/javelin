package instruction

import "github.com/alecthomas/participle/v2"

type inst_MOVN_253 struct {
                Xd uint32 `"MOVN" "X" @Integer`
                Imm uint32 `"#" @Integer`
      // LSL
      // imm:shift
      // <nil>
}

type InstructionMOVN_253 struct {
    Fields inst_MOVN_253 `@@`
}
func (i InstructionMOVN_253) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionMOVN_253) Imm() uint32 {
    return i.Fields.Imm
}

func init() {
	var p = participle.MustBuild[InstructionMOVN_253](
		participle.Lexer(asmDef),
	)
    parsers["MOVN"] = append(parsers["MOVN"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

