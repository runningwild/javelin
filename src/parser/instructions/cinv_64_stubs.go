package instruction

import "github.com/runningwild/javelin/machine"


func (i InstructionCINV_64) Encode() uint32 {
  // TODO: fill this out
  return 0
}

func (i InstructionCINV_64) Mask() uint32 {
  // TODO: fill this out
  return 0
}

func (i InstructionCINV_64) Execute(m *machine.Machine) {
  // TODO: fill this out
}

func DecodeCINV_64(v uint32) (*InstructionCINV_64, error) {
  // TODO: fill this out
  return nil, nil
}

