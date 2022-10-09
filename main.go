package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type OpCode uint32

const (
	OP_MOVE     OpCode = iota /* A B	R(A) := R(B)					*/
	OP_LOADK                  /* A Bx	R(A) := Kst(Bx)					*/
	OP_LOADKX                 /* A 	R(A) := Kst(extra arg)				*/
	OP_LOADBOOL               /* A B C	R(A) := (Bool)B; if (C) pc++			*/
	OP_LOADNIL                /* A B	R(A), R(A+1), ..., R(A+B) := nil		*/
	OP_GETUPVAL               /* A B	R(A) := UpValue[B]				*/

	OP_GETTABUP /* A B C	R(A) := UpValue[B][RK(C)]			*/
	OP_GETTABLE /* A B C	R(A) := R(B)[RK(C)]				*/

	OP_SETTABUP /* A B C	UpValue[A][RK(B)] := RK(C)			*/
	OP_SETUPVAL /* A B	UpValue[B] := R(A)				*/
	OP_SETTABLE /* A B C	R(A)[RK(B)] := RK(C)				*/

	OP_NEWTABLE /* A B C	R(A) := {} (size = B,C)				*/

	OP_SELF /* A B C	R(A+1) := R(B); R(A) := R(B)[RK(C)]		*/

	OP_ADD /* A B C	R(A) := RK(B) + RK(C)				*/
	OP_SUB /* A B C	R(A) := RK(B) - RK(C)				*/
	OP_MUL /* A B C	R(A) := RK(B) * RK(C)				*/
	OP_DIV /* A B C	R(A) := RK(B) / RK(C)				*/
	OP_MOD /* A B C	R(A) := RK(B) % RK(C)				*/
	OP_POW /* A B C	R(A) := RK(B) ^ RK(C)				*/
	OP_UNM /* A B	R(A) := -R(B)					*/
	OP_NOT /* A B	R(A) := not R(B)				*/
	OP_LEN /* A B	R(A) := length of R(B)				*/

	OP_CONCAT /* A B C	R(A) := R(B).. ... ..R(C)			*/

	OP_JMP /* A sBx	pc+=sBx; if (A) close all upvalues >= R(A) + 1	*/
	OP_EQ  /* A B C	if ((RK(B) == RK(C)) ~= A) then pc++		*/
	OP_LT  /* A B C	if ((RK(B) <  RK(C)) ~= A) then pc++		*/
	OP_LE  /* A B C	if ((RK(B) <= RK(C)) ~= A) then pc++		*/

	OP_TEST    /* A C	if not (R(A) <=> C) then pc++			*/
	OP_TESTSET /* A B C	if (R(B) <=> C) then R(A) := R(B) else pc++	*/

	OP_CALL     /* A B C	R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1)) */
	OP_TAILCALL /* A B C	return R(A)(R(A+1), ... ,R(A+B-1))		*/
	OP_RETURN   /* A B	return R(A), ... ,R(A+B-2)	(see note)	*/

	OP_FORLOOP /* A sBx	R(A)+=R(A+2);
	   if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }*/
	OP_FORPREP /* A sBx	R(A)-=R(A+2); pc+=sBx				*/

	OP_TFORCALL /* A C	R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));	*/
	OP_TFORLOOP /* A sBx	if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx }*/

	OP_SETLIST /* A B C	R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B	*/

	OP_CLOSURE /* A Bx	R(A) := closure(KPROTO[Bx])			*/

	OP_VARARG /* A B	R(A), R(A+1), ..., R(A+B-2) = vararg		*/

	OP_EXTRAARG /* Ax	extra (larger) argument for previous opcode	*/
)

func main() {
	fileBytes, err := os.ReadFile("luac.out") // TODO: pass as argument
	if err != nil {
		panic(err)
	}

	fmt.Printf("First byte value 0x%x\n", fileBytes[0])
	fmt.Printf("Found %s\n", string(fileBytes[1:4]))
	fmt.Printf("Version number 0x%x", fileBytes[4])
	isLittleEndian := fileBytes[6] == 1
	if isLittleEndian {
		fmt.Println(" Little endian")
	} else {
		fmt.Println(" Big endian")
	}
	integer_size := fileBytes[7]
	fmt.Printf("Integer Size: %d\n", integer_size)
	instruction_size := int(fileBytes[9])
	fmt.Println("Instruction Size", instruction_size)

	functionHeader := fileBytes[18:29] // 18th byte is the start of the "main chunk"
	fmt.Printf("line number %d\n", functionHeader[0])

	pc := 29                                                                                  // Start of the main chunk
	length_of_instructions := binary.LittleEndian.Uint32(fileBytes[pc : pc+instruction_size]) // TODO: Determine from endianess in header
	pc += instruction_size                                                                    // TODO: set this based on instruction size in the header

	for index := 0; index < int(length_of_instructions); index++ {
		bytecode_instruction := fileBytes[pc : pc+instruction_size]
		instruction_as_int := binary.LittleEndian.Uint32(bytecode_instruction)

		instruction := instruction_as_int & 0b111111

		pc += instruction_size
		fmt.Println("=============")
		switch OpCode(instruction) {
		case OP_MOVE:
			fmt.Println("OP Code: OP_MOVE")
			a := (bytecode_instruction[0] & 0b11000000 >> 6) + (bytecode_instruction[1] & 0b00111111 << 2)
			b := (bytecode_instruction[1] & 0b11000000 >> 6) + (bytecode_instruction[2] & 0b01111111 << 2)
			c := (bytecode_instruction[2] & 0b10000000 >> 7) + (bytecode_instruction[3] << 1)
			fmt.Printf("A %d: B %d: C: %d\n", a, b, c)
			break
		case OP_LOADK:
			fmt.Println("OP Code: OP_LOADK")
			a := (instruction_as_int >> 6) & 0xFF
			bx := (instruction_as_int >> 14) & 0x3FFFF
			fmt.Printf("A %d: Bx %x: \n", a, bx)
		case OP_GETTABUP:

		case OP_ADD:
			fmt.Println("OP Code: OP_ADD")
			break
		default:
			fmt.Println("Unrecognized Opcode: ", instruction)
			break
		}
	}

}
