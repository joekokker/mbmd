package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register(NewElectrexProducer)
}

const (
	METERTYPE_ELECTREX = "ELECTREX"
)

type ElectrexProducer struct {
	Opcodes
}

func NewElectrexProducer() Producer {
	/**
	 * Opcodes as defined by Electrex.
	 * See https://www.electrex.it/en/download/Manuals/User_manual_(extended)_Flash_96.pdf
	 * The registers are numbered with integers. Convert to HEX for the Opcodes in mbmd.
	 */
	ops := Opcodes{
		VoltageL1:         0x00d6,  // Phase  to  Neutral  Voltage
		VoltageL2:         0x00d8,  // Phase  to  Neutral  Voltage
		VoltageL3:         0x00da,  // Phase  to  Neutral  Voltage
		Voltage:           0x0106,  // Phase  to  Neutral  Mean  Voltage
		VoltageL1toL2:     0x00dc,  //  Phase  to  Phase  Voltage
		VoltageL2toL3:     0x00de,  // Phase  to  Phase  Voltage
		VoltageL3toL1:     0x00e0,  // Phase  to  Phase  Voltage
		VoltageLtoL:       0x0108,  // Phase  to  Phase  Mean  Voltage
		CurrentL1:         0x00e2,  // Line current
		CurrentL2:         0x00e4,  // Line current
		CurrentL3:         0x00e6,  // Line current
		CurrentN:          0x00e8  // Neutral Current
		Current:           0x010a,  //  Three phase current
		PowerL1:           0x00ea,  // Phase Active Power (+/-)
		PowerL2:           0x00ec,  // Phase Active Power (+/-)
		PowerL3:           0x00ee,  // Phase Active Power (+/-)
		Power:             0x010c,  // Total Active Power (+/-)
		ApparentPowerL1:   0x00f6,  // Phase Apparent Power
		ApparentPowerL2:   0x00f8,  // Phase Apparent Power
		ApparentPowerL3:   0x00fa,  // Phase Apparent Power
		ApparentPower:     0x0110,  // Total apparent power
		ReactivePowerL1    0x00f0,  // Phase Reactive Power (+/-)
		ReactivePowerL2    0x00f2,  // Phase Reactive Power (+/-)
		ReactivePowerL3    0x00f4,  // Phase Reactive Power (+/-)
		ReactivePower:     0x010e,  // Total reactive power (+/-)
		Import:            0x0114,  // Total import Active Power, AVG
		Export:            0x011c,  // Total export Active Power, AVG
		ApparentImport:    0x011a,  // Total import apparent power,AVG
		ApparentExport:    0x0122,  // Total export apparent power,AVG
		//ReactiveImportInductive:      0x0116,  // Total import reactive power,AVG
		//ReactiveImportCapacitive:     0x0118,  // Total import reactive power,AVG
		//ReactiveExportInductive:      0x011e,  // Total export reactive power,AVG
		//ReactiveExportCapacitive:     0x0120,  // Total export reactive power,AVG
		THDL1:             0x00c8,  // Phase to neutral Voltage, THD
		THDL2:             0x00ca,  // Phase to neutral Voltage, THD
		THDL3:             0x00cc,  // Phase to neutral Voltage, THD
		THD:               0x0102,  // Phase Voltage, Mean THD
		Frequency:         0x00d4,  // Voltage Input Frequency 
		THDL1Current:      0x00ce,  // Line current, THD 
		THDL2Current:      0x00d0,  // Line current, THD 
		THDL3Current:      0x00d2,  // Line current, THD
		THDCurrent:        0x0104,  // Line current, Mean THD
		PowerFactorL1:     0x00fc,  // Phase Power Factor (+/-)
		PowerFactorL2:     0x00fe,  // Phase Power Factor (+/-) 
		PowerFactorL3:     0x0100,  // Phase Power Factor (+/-) 
		PowerFactor:       0x0112,  // Total Power Factor (+/-) 
		//ApparentImportPower: 0x0064,
	}
	return &ElectrexProducer{Opcodes: ops}
}

func (p *ElectrexProducer) Type() string {
	return METERTYPE_ELECTREX
}

func (p *ElectrexProducer) Description() string {
	return "Electrex Flash N X3M meters"
}

func (p *ElectrexProducer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *ElectrexProducer) Probe() Operation {
	return p.snip(Voltage)
}

func (p *ElectrexProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
