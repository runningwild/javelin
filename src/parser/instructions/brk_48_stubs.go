package instruction

import "github.com/runningwild/javelin/machine"


func (i InstructionBRK_48) Encode() uint32 {
  // TODO: fill this out
  return 0
}

func (i InstructionBRK_48) Mask() uint32 {
  // TODO: fill this out
  return 0
}

func (i InstructionBRK_48) Execute(m *machine.Machine) {
  // TODO: fill this out
}

func DecodeBRK_48(v uint32) (*InstructionBRK_48, error) {
  // TODO: fill this out
  return nil, nil
}

