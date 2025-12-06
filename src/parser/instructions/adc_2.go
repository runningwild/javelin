package instruction

import "github.com/alecthomas/participle/v2"

type inst_ADC_2 struct {
                Xd uint32 `"ADC" "X" @Integer`
                Xn uint32 `"," "X" @Integer`
                Xm uint32 `"," "X" @Integer`
}

type InstructionADC_2 struct {
    Fields inst_ADC_2 `@@`
}
func (i InstructionADC_2) Xd() uint32 {
    return i.Fields.Xd
}
func (i InstructionADC_2) Xn() uint32 {
    return i.Fields.Xn
}
func (i InstructionADC_2) Xm() uint32 {
    return i.Fields.Xm
}

func init() {
	var p = participle.MustBuild[InstructionADC_2](
		participle.Lexer(asmDef),
	)
    parsers["ADC"] = append(parsers["ADC"], func(s string) (Instruction, error) {
		inst, err := p.ParseString("", s)
        if err != nil {
            return nil, nil
        }
        return inst, nil
    })
}

