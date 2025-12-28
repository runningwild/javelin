package instruction

import "github.com/alecthomas/participle/v2"

type inst_MOVK_252 struct {
                Xd uint32 `"MOVK" "X" @Integer`
                Imm uint32 `"#" @Integer`
      // LSL
      // imm:shift
      // <nil>
}

type InstructionMOVK_252 struct {
    Fields inst_MOVK_252 `@@`
}
func (i InstructionMOVK_252) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionMOVK_252) Imm() uint32 {
    return i.Fields.Imm
}

func init() {
	var p = participle.MustBuild[InstructionMOVK_252](
		participle.Lexer(asmDef),
	)
    parsers["MOVK"] = append(parsers["MOVK"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

