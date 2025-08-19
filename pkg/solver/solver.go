package solver

import (
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/bytecodealliance/wasmtime-go/v35"
)

type Solver struct {
	store    *wasmtime.Store
	memory   *wasmtime.Memory
	instance *wasmtime.Instance

	alloc      *wasmtime.Extern
	allocFn    *wasmtime.Func
	solve      *wasmtime.Extern
	solveFn    *wasmtime.Func
	stackPtr   *wasmtime.Extern
	stackPtrFn *wasmtime.Func
}

//go:embed sha3_wasm_bg.7b9ca65ddd.wasm
var wasmBytes []byte

func New() (*Solver, error) {
	engine := wasmtime.NewEngine()

	module, err := wasmtime.NewModule(engine, wasmBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create module: %w", err)
	}

	store := wasmtime.NewStore(engine)
	linker := wasmtime.NewLinker(engine)

	err = linker.DefineWasi()
	if err != nil {
		return nil, fmt.Errorf("failed to define wasi: %w", err)
	}

	instance, err := linker.Instantiate(store, module)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %w", err)
	}

	memory := instance.GetExport(store, "memory").Memory()
	if memory == nil {
		return nil, errors.New("failed to get memory export")
	}

	// init functions
	alloc := instance.GetExport(store, "__wbindgen_export_0")
	if alloc == nil {
		return nil, errors.New("failed to get memory export")
	}
	allocFn := alloc.Func()

	stackPtr := instance.GetExport(store, "__wbindgen_add_to_stack_pointer")
	if stackPtr == nil {
		return nil, errors.New("stack pointer function not found")
	}
	stackPtrFn := stackPtr.Func()

	solve := instance.GetExport(store, "wasm_solve")
	if solve == nil {
		return nil, errors.New("solve function not found")
	}
	solveFn := solve.Func()

	s := &Solver{}

	s.memory = memory
	s.store = store
	s.instance = instance
	s.alloc = alloc
	s.allocFn = allocFn
	s.solve = solve
	s.solveFn = solveFn
	s.stackPtr = stackPtr
	s.stackPtrFn = stackPtrFn

	return s, nil
}

func (s *Solver) Close() {
	s.store.Close()
	s.solve.Close()
	s.alloc.Close()
	s.stackPtr.Close()
}

func (s *Solver) writeToMemory(text string) (int32, int32, error) {
	textBytes := []byte(text)
	length := int32(len(textBytes))

	// Allocate memory
	result, err := s.allocFn.Call(s.store, length, 1)
	if err != nil {
		return 0, 0, fmt.Errorf("allocation failed: %w", err)
	}
	ptr := result.(int32)

	// Write to memory
	mem := s.memory.UnsafeData(s.store)
	copy(mem[ptr:ptr+length], textBytes)

	return ptr, length, nil
}

func (s *Solver) CalculateHash(challenge, salt string, difficulty, expireAt int) (int64, error) {
	prefix := fmt.Sprintf("%s_%d_", salt, expireAt)

	retptrRaw, err := s.stackPtrFn.Call(s.store, -16)
	if err != nil {
		return 0, fmt.Errorf("stack pointer adjustment failed: %w", err)
	}
	retptr := retptrRaw.(int32)

	// Write challenge and prefix to memory
	challengePtr, challengeLen, err := s.writeToMemory(challenge)
	if err != nil {
		return 0, fmt.Errorf("challenge write failed: %w", err)
	}

	prefixPtr, prefixLen, err := s.writeToMemory(prefix)
	if err != nil {
		return 0, fmt.Errorf("prefix write failed: %w", err)
	}

	// Get solve function
	_, err = s.solveFn.Call(s.store,
		retptr,
		challengePtr,
		challengeLen,
		prefixPtr,
		prefixLen,
		float64(difficulty),
	)
	if err != nil {
		return 0, fmt.Errorf("solve function failed: %w", err)
	}

	// Read result from memory
	mem := s.memory.UnsafeData(s.store)
	status := binary.LittleEndian.Uint32(mem[retptr : retptr+4])

	if status == 0 {
		return 0, errors.New("no solution found")
	}

	valueBytes := mem[retptr+8 : retptr+16]
	value := binary.LittleEndian.Uint64(valueBytes)
	floatValue := math.Float64frombits(value) // Convert bytes to float64

	// Convert float to integer (matches Python's int() behavior)
	answer := int64(floatValue)

	// Reset stackPtr pointer
	_, err = s.stackPtrFn.Call(s.store, 16)
	if err != nil {
		return 0, fmt.Errorf("stackPtr pointer reset failed: %w", err)
	}

	return answer, nil
}
