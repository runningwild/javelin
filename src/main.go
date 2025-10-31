package main

import (
	"fmt"
	"os"
)

func main() {
	addAsm := `
add r2,r3, r5
add v15.16b, v2.16b, v5.16b
`
	p, err := ParseProgram(addAsm)
	if err != nil {
		fmt.Printf("Failed to parse ADD: %v", err)
		os.Exit(1)
	}
	for _, i := range p.Instructions {
		fmt.Printf("%v\n", i)
	}

	return

	// fmt.Println("Starting ARMv8-M Simulator...")

	// m := NewMachine()

	// // Initialize some general-purpose registers and memory for testing
	// m.R[0] = 0x00001000 // R0 will be used as a base address
	// m.R[1] = 0x00000005 // R1 = 5
	// m.R[2] = 0x00000003 // R2 = 3

	// fmt.Println("\n--- Initial CPU State (GP Registers) ---")
	// m.PrintRegisters()
	// m.PrintMemoryRegion(m.R[0], 16) // Print 16 bytes from R0's address

	// assemblyInstructionsGP := []string{
	// 	"STR R1, [R0, #0]", // Store R1 (5) at Mem[R0 + 0]
	// 	"STR R2, [R0, #4]", // Store R2 (3) at Mem[R0 + 4]
	// 	"LDR R3, [R0, #0]", // Load value from Mem[0x00001000] into R3
	// 	"LDR R4, [R0, #4]", // Load value from Mem[0x00001004] into R4
	// 	"ADD R5, R3, R4",   // R5 = R3 + R4
	// 	"STR R5, [R0, #8]", // Store R5 into Mem[0x00001008]
	// }

	// fmt.Println("\n--- Executing General-Purpose Instructions ---")
	// for _, asm := range assemblyInstructionsGP {
	// 	fmt.Printf("Executing: %s\n", asm)
	// 	inst, err := ParseInstruction(asm)
	// 	if err != nil {
	// 		fmt.Printf("Error parsing instruction '%s': %v\n", asm, err)
	// 		return
	// 	}
	// 	inst.Execute(m)
	// }

	// fmt.Println("\n--- CPU State After GP Instructions ---")
	// m.PrintRegisters()
	// m.PrintMemoryRegion(m.R[0], 16) // Print 16 bytes from R0's address

	// // --- NEON Instruction Test ---
	// fmt.Println("\n--- Initializing Memory for NEON Test ---")
	// // Put some data into memory for NEON operations
	// // Mem[0x1000] = 1, 2, 3, 4 (uint32s)
	// // Mem[0x1010] = 5, 6, 7, 8 (uint32s)
	// binary.LittleEndian.PutUint32(m.Mem[0x1000:0x1004], 1)
	// binary.LittleEndian.PutUint32(m.Mem[0x1004:0x1008], 2)
	// binary.LittleEndian.PutUint32(m.Mem[0x1008:0x100C], 3)
	// binary.LittleEndian.PutUint32(m.Mem[0x100C:0x1010], 4)

	// binary.LittleEndian.PutUint32(m.Mem[0x1010:0x1014], 5)
	// binary.LittleEndian.PutUint32(m.Mem[0x1014:0x1018], 6)
	// binary.LittleEndian.PutUint32(m.Mem[0x1018:0x101C], 7)
	// binary.LittleEndian.PutUint32(m.Mem[0x101C:0x1020], 8)

	// m.R[0] = 0x1000 // Base address for first vector
	// m.R[1] = 0x1010 // Base address for second vector

	// fmt.Println("\n--- CPU State Before NEON Instructions ---")
	// m.PrintRegisters()
	// m.PrintMemoryRegion(0x1000, 32) // Print 32 bytes for NEON data

	// assemblyInstructionsNEON := []string{
	// 	"VLD1.32 {D0, D1}, [R0]", // Load 1,2,3,4 into D0, D1
	// 	"VLD1.32 {D2, D3}, [R1]", // Load 5,6,7,8 into D2, D3
	// 	"VADD.I32 D4, D0, D2",    // D4 = D0 + D2 (1+5, 2+6)
	// 	"VADD.I32 D5, D1, D3",    // D5 = D1 + D3 (3+7, 4+8)
	// 	"VST1.32 {D4, D5}, [R0]", // Store D4, D5 back to [R0]
	// }

	// fmt.Println("\n--- Executing NEON Instructions ---")
	// for _, asm := range assemblyInstructionsNEON {
	// 	fmt.Printf("Executing: %s\n", asm)
	// 	inst, err := ParseInstruction(asm)
	// 	if err != nil {
	// 		fmt.Printf("Error parsing instruction '%s': %v\n", asm, err)
	// 		return
	// 	}
	// 	inst.Execute(m)
	// }

	// fmt.Println("\n--- Final CPU State After NEON Instructions ---")
	// m.PrintRegisters()
	// m.PrintMemoryRegion(0x1000, 32) // Print 32 bytes for NEON data

	// fmt.Println("Simulator finished.")
}
