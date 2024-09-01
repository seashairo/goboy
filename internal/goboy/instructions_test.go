package goboy

import (
	"encoding/json"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestInstructions_00(t *testing.T) { testFile(t, "00.json") }
func TestInstructions_01(t *testing.T) { testFile(t, "01.json") }
func TestInstructions_02(t *testing.T) { testFile(t, "02.json") }
func TestInstructions_03(t *testing.T) { testFile(t, "03.json") }
func TestInstructions_04(t *testing.T) { testFile(t, "04.json") }
func TestInstructions_05(t *testing.T) { testFile(t, "05.json") }
func TestInstructions_06(t *testing.T) { testFile(t, "06.json") }
func TestInstructions_07(t *testing.T) { testFile(t, "07.json") }
func TestInstructions_08(t *testing.T) { testFile(t, "08.json") }
func TestInstructions_09(t *testing.T) { testFile(t, "09.json") }
func TestInstructions_0a(t *testing.T) { testFile(t, "0a.json") }
func TestInstructions_0b(t *testing.T) { testFile(t, "0b.json") }
func TestInstructions_0c(t *testing.T) { testFile(t, "0c.json") }
func TestInstructions_0d(t *testing.T) { testFile(t, "0d.json") }
func TestInstructions_0e(t *testing.T) { testFile(t, "0e.json") }
func TestInstructions_0f(t *testing.T) { testFile(t, "0f.json") }

func TestInstructions_11(t *testing.T) { testFile(t, "11.json") }
func TestInstructions_12(t *testing.T) { testFile(t, "12.json") }
func TestInstructions_13(t *testing.T) { testFile(t, "13.json") }
func TestInstructions_14(t *testing.T) { testFile(t, "14.json") }
func TestInstructions_15(t *testing.T) { testFile(t, "15.json") }
func TestInstructions_16(t *testing.T) { testFile(t, "16.json") }
func TestInstructions_17(t *testing.T) { testFile(t, "17.json") }
func TestInstructions_18(t *testing.T) { testFile(t, "18.json") }
func TestInstructions_19(t *testing.T) { testFile(t, "19.json") }
func TestInstructions_1a(t *testing.T) { testFile(t, "1a.json") }
func TestInstructions_1b(t *testing.T) { testFile(t, "1b.json") }
func TestInstructions_1c(t *testing.T) { testFile(t, "1c.json") }
func TestInstructions_1d(t *testing.T) { testFile(t, "1d.json") }
func TestInstructions_1e(t *testing.T) { testFile(t, "1e.json") }
func TestInstructions_1f(t *testing.T) { testFile(t, "1f.json") }

func TestInstructions_20(t *testing.T) { testFile(t, "20.json") }
func TestInstructions_21(t *testing.T) { testFile(t, "21.json") }
func TestInstructions_22(t *testing.T) { testFile(t, "22.json") }
func TestInstructions_23(t *testing.T) { testFile(t, "23.json") }
func TestInstructions_24(t *testing.T) { testFile(t, "24.json") }
func TestInstructions_25(t *testing.T) { testFile(t, "25.json") }
func TestInstructions_26(t *testing.T) { testFile(t, "26.json") }
func TestInstructions_27(t *testing.T) { testFile(t, "27.json") }
func TestInstructions_28(t *testing.T) { testFile(t, "28.json") }
func TestInstructions_29(t *testing.T) { testFile(t, "29.json") }
func TestInstructions_2a(t *testing.T) { testFile(t, "2a.json") }
func TestInstructions_2b(t *testing.T) { testFile(t, "2b.json") }
func TestInstructions_2c(t *testing.T) { testFile(t, "2c.json") }
func TestInstructions_2d(t *testing.T) { testFile(t, "2d.json") }
func TestInstructions_2e(t *testing.T) { testFile(t, "2e.json") }
func TestInstructions_2f(t *testing.T) { testFile(t, "2f.json") }

func TestInstructions_30(t *testing.T) { testFile(t, "30.json") }
func TestInstructions_31(t *testing.T) { testFile(t, "31.json") }
func TestInstructions_32(t *testing.T) { testFile(t, "32.json") }
func TestInstructions_33(t *testing.T) { testFile(t, "33.json") }
func TestInstructions_34(t *testing.T) { testFile(t, "34.json") }
func TestInstructions_35(t *testing.T) { testFile(t, "35.json") }
func TestInstructions_36(t *testing.T) { testFile(t, "36.json") }
func TestInstructions_37(t *testing.T) { testFile(t, "37.json") }
func TestInstructions_38(t *testing.T) { testFile(t, "38.json") }
func TestInstructions_39(t *testing.T) { testFile(t, "39.json") }
func TestInstructions_3a(t *testing.T) { testFile(t, "3a.json") }
func TestInstructions_3b(t *testing.T) { testFile(t, "3b.json") }
func TestInstructions_3c(t *testing.T) { testFile(t, "3c.json") }
func TestInstructions_3d(t *testing.T) { testFile(t, "3d.json") }
func TestInstructions_3e(t *testing.T) { testFile(t, "3e.json") }
func TestInstructions_3f(t *testing.T) { testFile(t, "3f.json") }

func TestInstructions_40(t *testing.T) { testFile(t, "40.json") }
func TestInstructions_41(t *testing.T) { testFile(t, "41.json") }
func TestInstructions_42(t *testing.T) { testFile(t, "42.json") }
func TestInstructions_43(t *testing.T) { testFile(t, "43.json") }
func TestInstructions_44(t *testing.T) { testFile(t, "44.json") }
func TestInstructions_45(t *testing.T) { testFile(t, "45.json") }
func TestInstructions_46(t *testing.T) { testFile(t, "46.json") }
func TestInstructions_47(t *testing.T) { testFile(t, "47.json") }
func TestInstructions_48(t *testing.T) { testFile(t, "48.json") }
func TestInstructions_49(t *testing.T) { testFile(t, "49.json") }
func TestInstructions_4a(t *testing.T) { testFile(t, "4a.json") }
func TestInstructions_4b(t *testing.T) { testFile(t, "4b.json") }
func TestInstructions_4c(t *testing.T) { testFile(t, "4c.json") }
func TestInstructions_4d(t *testing.T) { testFile(t, "4d.json") }
func TestInstructions_4e(t *testing.T) { testFile(t, "4e.json") }
func TestInstructions_4f(t *testing.T) { testFile(t, "4f.json") }

func TestInstructions_50(t *testing.T) { testFile(t, "50.json") }
func TestInstructions_51(t *testing.T) { testFile(t, "51.json") }
func TestInstructions_52(t *testing.T) { testFile(t, "52.json") }
func TestInstructions_53(t *testing.T) { testFile(t, "53.json") }
func TestInstructions_54(t *testing.T) { testFile(t, "54.json") }
func TestInstructions_55(t *testing.T) { testFile(t, "55.json") }
func TestInstructions_56(t *testing.T) { testFile(t, "56.json") }
func TestInstructions_57(t *testing.T) { testFile(t, "57.json") }
func TestInstructions_58(t *testing.T) { testFile(t, "58.json") }
func TestInstructions_59(t *testing.T) { testFile(t, "59.json") }
func TestInstructions_5a(t *testing.T) { testFile(t, "5a.json") }
func TestInstructions_5b(t *testing.T) { testFile(t, "5b.json") }
func TestInstructions_5c(t *testing.T) { testFile(t, "5c.json") }
func TestInstructions_5d(t *testing.T) { testFile(t, "5d.json") }
func TestInstructions_5e(t *testing.T) { testFile(t, "5e.json") }
func TestInstructions_5f(t *testing.T) { testFile(t, "5f.json") }

func TestInstructions_60(t *testing.T) { testFile(t, "60.json") }
func TestInstructions_61(t *testing.T) { testFile(t, "61.json") }
func TestInstructions_62(t *testing.T) { testFile(t, "62.json") }
func TestInstructions_63(t *testing.T) { testFile(t, "63.json") }
func TestInstructions_64(t *testing.T) { testFile(t, "64.json") }
func TestInstructions_65(t *testing.T) { testFile(t, "65.json") }
func TestInstructions_66(t *testing.T) { testFile(t, "66.json") }
func TestInstructions_67(t *testing.T) { testFile(t, "67.json") }
func TestInstructions_68(t *testing.T) { testFile(t, "68.json") }
func TestInstructions_69(t *testing.T) { testFile(t, "69.json") }
func TestInstructions_6a(t *testing.T) { testFile(t, "6a.json") }
func TestInstructions_6b(t *testing.T) { testFile(t, "6b.json") }
func TestInstructions_6c(t *testing.T) { testFile(t, "6c.json") }
func TestInstructions_6d(t *testing.T) { testFile(t, "6d.json") }
func TestInstructions_6e(t *testing.T) { testFile(t, "6e.json") }
func TestInstructions_6f(t *testing.T) { testFile(t, "6f.json") }

func TestInstructions_70(t *testing.T) { testFile(t, "70.json") }
func TestInstructions_71(t *testing.T) { testFile(t, "71.json") }
func TestInstructions_72(t *testing.T) { testFile(t, "72.json") }
func TestInstructions_73(t *testing.T) { testFile(t, "73.json") }
func TestInstructions_74(t *testing.T) { testFile(t, "74.json") }
func TestInstructions_75(t *testing.T) { testFile(t, "75.json") }
func TestInstructions_77(t *testing.T) { testFile(t, "77.json") }
func TestInstructions_78(t *testing.T) { testFile(t, "78.json") }
func TestInstructions_79(t *testing.T) { testFile(t, "79.json") }
func TestInstructions_7a(t *testing.T) { testFile(t, "7a.json") }
func TestInstructions_7b(t *testing.T) { testFile(t, "7b.json") }
func TestInstructions_7c(t *testing.T) { testFile(t, "7c.json") }
func TestInstructions_7d(t *testing.T) { testFile(t, "7d.json") }
func TestInstructions_7e(t *testing.T) { testFile(t, "7e.json") }
func TestInstructions_7f(t *testing.T) { testFile(t, "7f.json") }

func TestInstructions_80(t *testing.T) { testFile(t, "80.json") }
func TestInstructions_81(t *testing.T) { testFile(t, "81.json") }
func TestInstructions_82(t *testing.T) { testFile(t, "82.json") }
func TestInstructions_83(t *testing.T) { testFile(t, "83.json") }
func TestInstructions_84(t *testing.T) { testFile(t, "84.json") }
func TestInstructions_85(t *testing.T) { testFile(t, "85.json") }
func TestInstructions_86(t *testing.T) { testFile(t, "86.json") }
func TestInstructions_87(t *testing.T) { testFile(t, "87.json") }
func TestInstructions_88(t *testing.T) { testFile(t, "88.json") }
func TestInstructions_89(t *testing.T) { testFile(t, "89.json") }
func TestInstructions_8a(t *testing.T) { testFile(t, "8a.json") }
func TestInstructions_8b(t *testing.T) { testFile(t, "8b.json") }
func TestInstructions_8c(t *testing.T) { testFile(t, "8c.json") }
func TestInstructions_8d(t *testing.T) { testFile(t, "8d.json") }
func TestInstructions_8e(t *testing.T) { testFile(t, "8e.json") }
func TestInstructions_8f(t *testing.T) { testFile(t, "8f.json") }

func TestInstructions_90(t *testing.T) { testFile(t, "90.json") }
func TestInstructions_91(t *testing.T) { testFile(t, "91.json") }
func TestInstructions_92(t *testing.T) { testFile(t, "92.json") }
func TestInstructions_93(t *testing.T) { testFile(t, "93.json") }
func TestInstructions_94(t *testing.T) { testFile(t, "94.json") }
func TestInstructions_95(t *testing.T) { testFile(t, "95.json") }
func TestInstructions_96(t *testing.T) { testFile(t, "96.json") }
func TestInstructions_97(t *testing.T) { testFile(t, "97.json") }
func TestInstructions_98(t *testing.T) { testFile(t, "98.json") }
func TestInstructions_99(t *testing.T) { testFile(t, "99.json") }
func TestInstructions_9a(t *testing.T) { testFile(t, "9a.json") }
func TestInstructions_9b(t *testing.T) { testFile(t, "9b.json") }
func TestInstructions_9c(t *testing.T) { testFile(t, "9c.json") }
func TestInstructions_9d(t *testing.T) { testFile(t, "9d.json") }
func TestInstructions_9e(t *testing.T) { testFile(t, "9e.json") }
func TestInstructions_9f(t *testing.T) { testFile(t, "9f.json") }

func TestInstructions_a0(t *testing.T) { testFile(t, "a0.json") }
func TestInstructions_a1(t *testing.T) { testFile(t, "a1.json") }
func TestInstructions_a2(t *testing.T) { testFile(t, "a2.json") }
func TestInstructions_a3(t *testing.T) { testFile(t, "a3.json") }
func TestInstructions_a4(t *testing.T) { testFile(t, "a4.json") }
func TestInstructions_a5(t *testing.T) { testFile(t, "a5.json") }
func TestInstructions_a6(t *testing.T) { testFile(t, "a6.json") }
func TestInstructions_a7(t *testing.T) { testFile(t, "a7.json") }
func TestInstructions_a8(t *testing.T) { testFile(t, "a8.json") }
func TestInstructions_a9(t *testing.T) { testFile(t, "a9.json") }
func TestInstructions_aa(t *testing.T) { testFile(t, "aa.json") }
func TestInstructions_ab(t *testing.T) { testFile(t, "ab.json") }
func TestInstructions_ac(t *testing.T) { testFile(t, "ac.json") }
func TestInstructions_ad(t *testing.T) { testFile(t, "ad.json") }
func TestInstructions_ae(t *testing.T) { testFile(t, "ae.json") }
func TestInstructions_af(t *testing.T) { testFile(t, "af.json") }

func TestInstructions_b0(t *testing.T) { testFile(t, "b0.json") }
func TestInstructions_b1(t *testing.T) { testFile(t, "b1.json") }
func TestInstructions_b2(t *testing.T) { testFile(t, "b2.json") }
func TestInstructions_b3(t *testing.T) { testFile(t, "b3.json") }
func TestInstructions_b4(t *testing.T) { testFile(t, "b4.json") }
func TestInstructions_b5(t *testing.T) { testFile(t, "b5.json") }
func TestInstructions_b6(t *testing.T) { testFile(t, "b6.json") }
func TestInstructions_b7(t *testing.T) { testFile(t, "b7.json") }
func TestInstructions_b8(t *testing.T) { testFile(t, "b8.json") }
func TestInstructions_b9(t *testing.T) { testFile(t, "b9.json") }
func TestInstructions_ba(t *testing.T) { testFile(t, "ba.json") }
func TestInstructions_bb(t *testing.T) { testFile(t, "bb.json") }
func TestInstructions_bc(t *testing.T) { testFile(t, "bc.json") }
func TestInstructions_bd(t *testing.T) { testFile(t, "bd.json") }
func TestInstructions_be(t *testing.T) { testFile(t, "be.json") }
func TestInstructions_bf(t *testing.T) { testFile(t, "bf.json") }

func TestInstructions_c0(t *testing.T) { testFile(t, "c0.json") }
func TestInstructions_c1(t *testing.T) { testFile(t, "c1.json") }
func TestInstructions_c2(t *testing.T) { testFile(t, "c2.json") }
func TestInstructions_c3(t *testing.T) { testFile(t, "c3.json") }
func TestInstructions_c4(t *testing.T) { testFile(t, "c4.json") }
func TestInstructions_c5(t *testing.T) { testFile(t, "c5.json") }
func TestInstructions_c6(t *testing.T) { testFile(t, "c6.json") }
func TestInstructions_c7(t *testing.T) { testFile(t, "c7.json") }
func TestInstructions_c8(t *testing.T) { testFile(t, "c8.json") }
func TestInstructions_c9(t *testing.T) { testFile(t, "c9.json") }
func TestInstructions_ca(t *testing.T) { testFile(t, "ca.json") }
func TestInstructions_cb(t *testing.T) { testFile(t, "cb.json") }
func TestInstructions_cc(t *testing.T) { testFile(t, "cc.json") }
func TestInstructions_cd(t *testing.T) { testFile(t, "cd.json") }
func TestInstructions_ce(t *testing.T) { testFile(t, "ce.json") }
func TestInstructions_cf(t *testing.T) { testFile(t, "cf.json") }

func TestInstructions_d0(t *testing.T) { testFile(t, "d0.json") }
func TestInstructions_d1(t *testing.T) { testFile(t, "d1.json") }
func TestInstructions_d2(t *testing.T) { testFile(t, "d2.json") }
func TestInstructions_d4(t *testing.T) { testFile(t, "d4.json") }
func TestInstructions_d5(t *testing.T) { testFile(t, "d5.json") }
func TestInstructions_d6(t *testing.T) { testFile(t, "d6.json") }
func TestInstructions_d7(t *testing.T) { testFile(t, "d7.json") }
func TestInstructions_d8(t *testing.T) { testFile(t, "d8.json") }
func TestInstructions_d9(t *testing.T) { testFile(t, "d9.json") }
func TestInstructions_da(t *testing.T) { testFile(t, "da.json") }
func TestInstructions_dc(t *testing.T) { testFile(t, "dc.json") }
func TestInstructions_de(t *testing.T) { testFile(t, "de.json") }
func TestInstructions_df(t *testing.T) { testFile(t, "df.json") }

func TestInstructions_e0(t *testing.T) { testFile(t, "e0.json") }
func TestInstructions_e1(t *testing.T) { testFile(t, "e1.json") }
func TestInstructions_e2(t *testing.T) { testFile(t, "e2.json") }
func TestInstructions_e5(t *testing.T) { testFile(t, "e5.json") }
func TestInstructions_e6(t *testing.T) { testFile(t, "e6.json") }
func TestInstructions_e7(t *testing.T) { testFile(t, "e7.json") }
func TestInstructions_e8(t *testing.T) { testFile(t, "e8.json") }
func TestInstructions_e9(t *testing.T) { testFile(t, "e9.json") }
func TestInstructions_ea(t *testing.T) { testFile(t, "ea.json") }
func TestInstructions_ee(t *testing.T) { testFile(t, "ee.json") }
func TestInstructions_ef(t *testing.T) { testFile(t, "ef.json") }

func TestInstructions_f0(t *testing.T) { testFile(t, "f0.json") }
func TestInstructions_f1(t *testing.T) { testFile(t, "f1.json") }
func TestInstructions_f2(t *testing.T) { testFile(t, "f2.json") }
func TestInstructions_f5(t *testing.T) { testFile(t, "f5.json") }
func TestInstructions_f6(t *testing.T) { testFile(t, "f6.json") }
func TestInstructions_f7(t *testing.T) { testFile(t, "f7.json") }
func TestInstructions_f8(t *testing.T) { testFile(t, "f8.json") }
func TestInstructions_f9(t *testing.T) { testFile(t, "f9.json") }
func TestInstructions_fa(t *testing.T) { testFile(t, "fa.json") }
func TestInstructions_fe(t *testing.T) { testFile(t, "fe.json") }
func TestInstructions_ff(t *testing.T) { testFile(t, "ff.json") }

func testFile(t *testing.T, filename string) {
	gameboy := NewGameBoy()
	gameboy.bus = NewRAM(0x10000, 0)
	gameboy.cpu.bus = gameboy.bus

	testCases := loadJson(t, filename)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			executeTest(t, gameboy, testCase)
		})
	}
}

func executeTest(t *testing.T, gameboy *GameBoy, testCase TestCase) {
	init := testCase.Initial
	registers := gameboy.cpu.registers
	registers.a = init.A
	registers.f = init.F
	registers.b = init.B
	registers.c = init.C
	registers.d = init.D
	registers.e = init.E
	registers.h = init.H
	registers.l = init.L
	registers.pc = init.PC - 1
	registers.sp = init.SP

	// Load initial RAM state
	for _, ramEntry := range init.RAM {
		address := ramEntry[0]
		value := ramEntry[1]
		gameboy.bus.writeByte(address, byte(value))
	}

	gameboy.cpu.Tick()

	// Get the expected final state
	final := testCase.Final

	// Compare final CPU register values
	if registers.a != final.A {
		t.Errorf("Register A: expected %d, got %d", final.A, registers.a)
	}
	if registers.f != final.F {
		t.Errorf("Register F: expected %d, got %d", final.F, registers.f)
	}
	if registers.b != final.B {
		t.Errorf("Register B: expected %d, got %d", final.B, registers.b)
	}
	if registers.c != final.C {
		t.Errorf("Register C: expected %d, got %d", final.C, registers.c)
	}
	if registers.d != final.D {
		t.Errorf("Register D: expected %d, got %d", final.D, registers.d)
	}
	if registers.e != final.E {
		t.Errorf("Register E: expected %d, got %d", final.E, registers.e)
	}
	if registers.h != final.H {
		t.Errorf("Register H: expected %d, got %d", final.H, registers.h)
	}
	if registers.l != final.L {
		t.Errorf("Register L: expected %d, got %d", final.L, registers.l)
	}

	// Compare final PC and SP values
	if registers.pc != final.PC-1 {
		t.Errorf("Program Counter: expected %d, got %d", final.PC-1, registers.pc)
	}
	if registers.sp != final.SP {
		t.Errorf("Stack Pointer: expected %d, got %d", final.SP, registers.sp)
	}

	// Compare final RAM contents
	for _, ramEntry := range final.RAM {
		address := ramEntry[0]
		expectedValue := byte(ramEntry[1])
		actualValue := gameboy.bus.readByte(address)
		if actualValue != expectedValue {
			t.Errorf("RAM at address 0x%04X: expected %d, got %d", address, expectedValue, actualValue)
		}
	}
}

type State struct {
	A   byte       `json:"a"`
	B   byte       `json:"b"`
	C   byte       `json:"c"`
	D   byte       `json:"d"`
	E   byte       `json:"e"`
	F   byte       `json:"f"`
	H   byte       `json:"h"`
	L   byte       `json:"l"`
	PC  uint16     `json:"pc"`
	SP  uint16     `json:"sp"`
	RAM [][]uint16 `json:"ram"`
}

type TestCase struct {
	Name    string          `json:"name"`
	Initial State           `json:"initial"`
	Final   State           `json:"final"`
	Cycles  [][]interface{} `json:"cycles"`
}

func loadJson(t *testing.T, filename string) []TestCase {
	file, err := os.Open("data/test_data/" + filename)
	if err != nil {
		t.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var testCases []TestCase
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&testCases)
	if err != nil {
		t.Fatalf("Error decoding JSON: %v", err)
	}

	if len(testCases) == 0 {
		t.Fatalf("Expected at least one test case, found none")
	}

	return testCases
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
