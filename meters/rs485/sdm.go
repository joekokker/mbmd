package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register(NewSDMProducer)
}

const (
	METERTYPE_SDM = "SDM"
)

type SDMProducer struct {
	Opcodes
}

func NewSDMProducer() Producer {
	/**
	 * Opcodes as defined by Eastron SDM630.
	 * See http://bg-etech.de/download/manual/SDM630Register.pdf
	 * This is to a large extent a superset of all SDM devices, however there are
	 * subtle differences (see 220, 230). Some opcodes might not work on some devices.
	 */
	ops := Opcodes{
		VoltageL1:     0x0000, // 220, 230
		VoltageL2:     0x0002,
		VoltageL3:     0x0004,
		VoltageL1toL2: 0x00c8,  // Phase  to  Phase  Voltage
		VoltageL2toL3: 0x00ca,  // Phase  to  Phase  Voltage
		VoltageL3toL1: 0x00cc,  // Phase  to  Phase  Voltage
		VoltageLtoL:   0x00ce,  // Phase  to  Phase  Mean  Voltage
		CurrentL1:     0x0006,
		CurrentL2:     0x0008,
		CurrentL3:     0x000A,
		CurrentN:      0x00e0,
		Current:       0x0030,
		PowerL1:       0x000C,
		PowerL2:       0x000E,
		PowerL3:       0x0010,
		Power:         0x0034,
		ApparentPowerL1:   0x0012,  // Phase Apparent Power
		ApparentPowerL2:   0x0014,  // Phase Apparent Power
		ApparentPowerL3:   0x0016,  // Phase Apparent Power
		ApparentPower:     0x0038,  // Total apparent power
		ReactivePowerL1:   0x0018,  // Phase Reactive Power
		ReactivePowerL2:   0x001a,  // Phase Reactive Power
		ReactivePowerL3:   0x001c,  // Phase Reactive Power
		ReactivePower:     0x003c,  // Total reactive power
		ImportPower:   0x0054,
		ImportL1:      0x015a,
		ImportL2:      0x015c,
		ImportL3:      0x015e,
		Import:        0x0048, // 220, 230
		ExportL1:      0x0160,
		ExportL2:      0x0162,
		ExportL3:      0x0164,
		Export:        0x004a, // 220, 230
		SumL1:         0x0166,
		SumL2:         0x0168,
		SumL3:         0x016a,
		Sum:           0x0156, // 220
		CosphiL1:      0x001e, //      230
		CosphiL2:      0x0020,
		CosphiL3:      0x0022,
		Cosphi:        0x003e, // Total Power Factor (+/-) Positive for capacitive and negative for inductive.
		THDL1:         0x00ea, // voltage
		THDL2:         0x00ec, // voltage
		THDL3:         0x00ee, // voltage
		THD:           0x00f8, // voltage
		Frequency:     0x0046,
		THDL1Current:      0x00f0,  // Line current, THD 
		THDL2Current:      0x00f2,  // Line current, THD 
		THDL3Current:      0x00f4,  // Line current, THD
		THDCurrent:        0x00fa,  // Line current, Mean THD
		//ApparentImportPower: 0x0064,
	}
	return &SDMProducer{Opcodes: ops}
}

func (p *SDMProducer) Type() string {
	return METERTYPE_SDM
}

func (p *SDMProducer) Description() string {
	return "Eastron SDM630"
}

func (p *SDMProducer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *SDMProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *SDMProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
