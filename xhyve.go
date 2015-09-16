// +build darwin

package main

// #cgo CFLAGS: -I${SRCDIR}/vendor/xhyve/include
// #cgo LDFLAGS: -L${SRCDIR} -lxhyve -arch x86_64 -framework Hypervisor -framework vmnet
// #include "helper.h"
import "C"

import (
	"errors"
	"runtime"
	"strconv"
	"unsafe"

	"github.com/satori/go.uuid"
)

var (
	// ErrPCIDevice is returned when an error was found parsing PCI devices.
	ErrPCIDevice = errors.New("error parsing PCI device")
	// ErrLPCDevice is returned when an error was found parsing LPC device options.
	ErrLPCDevice = errors.New("error parsing LPC devices")
	// ErrInvalidMemsize is returned if memorize size is invalid.
	ErrInvalidMemsize = errors.New("invalid memory size")
	// ErrInvalidBootParams is returne when kexec or fbsd params are invalid.
	ErrInvalidBootParams = errors.New("boot parameters are invalid")
	// ErrCreatingVM is returned when xhyve was unable to create the virtual machine.
	ErrCreatingVM = errors.New("unable to create VM")
	// ErrMaxNumVCPUExceeded is returned when the number of vcpus requested for the guest
	// exceeds the limit imposed by xhyve.
	ErrMaxNumVCPUExceeded = errors.New("maximum number of vcpus requested is too high")
	// ErrSettingUpMemory is returned when an error was returned by xhyve when trying
	// to setup guest memory.
	ErrSettingUpMemory = errors.New("unable to setup memory for guest vm")
	// ErrInitializingMSR is returned when xhyve is unable to initialize MSR table
	ErrInitializingMSR = errors.New("unable to initialize MSR table")
	// ErrInitializingPCI is returned when xhyve is unable to initialize PCI emulation
	ErrInitializingPCI = errors.New("unable to initialize PCI emulation")
	// ErrBuildingMPTTable is returned when xhyve is unable to build MPT table
	ErrBuildingMPTTable = errors.New("unable to build MPT table")
	// ErrBuildingSMBIOS is returned when xhyve is unable to build smbios
	ErrBuildingSMBIOS = errors.New("unable to build smbios")
	// ErrBuildingACPI is returned when xhyve is unable to build ACPI
	ErrBuildingACPI = errors.New("unable to build ACPI")
)

// XHyveParams defines parameters needed by xhyve to boot up virtual machines.
type XHyveParams struct {
	// Number of CPUs to assigned to the guest vm.
	VCPUs int
	// Memory in megabytes to assign to guest vm.
	Memory string
	// PCI devices to attach to the guest vm, including bus and slot.
	// Example: []string{"2:0,virtio-net", "0:0,hostbridge"}
	PCIDevs []string
	// LPC devices to attach to the guest vm.
	LPCDevs []string // -l com1,stdio
	// Whether to create ACPI tables or not.
	ACPI *bool
	// Universal identifier for the guest vm.
	UUID string
	// Whether to use localtime or UTC in Real Time Clock.
	RTCLocaltime *bool
	// Either kexec or fbsd params. Format:
	// kexec,kernel image,initrd,"cmdline"
	// fbsd,userboot,boot volume,"kernel env"
	BootParams string
	// Whether to enable or disable bvm console
	BVMConsole *bool
	// Whether to enable or disable mpt table generation
	MPTGen *bool
}

func setDefaults(p *XHyveParams) {
	if p.VCPUs < 1 {
		p.VCPUs = 1
	}

	memsize, err := strconv.Atoi(p.Memory)
	if memsize < 256 || err != nil {
		p.Memory = "256"
	}

	// if len(p.PCISlots) == 0 {
	// 	p.PCISlots = []string{
	// 		"2:0,virtio-net",
	// 		"0:0,hostbridge",
	// 		"31,lpc",
	// 	}
	// }
	//
	// if len(p.LPCDevs) == 0 {
	// 	p.LPCDevs = []string{
	// 		"com1",
	// 		"stdio",
	// 	}
	// }

	if p.UUID == "" {
		p.UUID = uuid.NewV4().String()
	}

	if p.ACPI == nil {
		p.ACPI = new(bool)
	}

	if p.RTCLocaltime == nil {
		p.RTCLocaltime = new(bool)
	}

	if p.BVMConsole == nil {
		p.BVMConsole = new(bool)
	}

	if p.MPTGen == nil {
		p.MPTGen = new(bool)
		*p.MPTGen = true
	}
}

func init() {
	runtime.LockOSThread()
}

// RunXHyve runs xhyve hypervisor with the given parameters.
func RunXHyve(p XHyveParams) error {
	setDefaults(&p)

	for _, d := range p.PCIDevs {
		device := C.CString(d)
		// defer is not adviced to have within a loop but we are not expecting a lot of PCI devices.
		defer C.free(unsafe.Pointer(device))
		if err := C.pci_parse_slot(device); err != 0 {
			return ErrPCIDevice
		}
	}

	for _, d := range p.LPCDevs {
		device := C.CString(d)
		// defer is not adviced to have within a loop but we are not expecting a lot of PCI devices.
		defer C.free(unsafe.Pointer(device))
		if err := C.lpc_device_parse(device); err != 0 {
			return ErrLPCDevice
		}
	}

	bootParams := C.CString(p.BootParams)
	defer C.free(unsafe.Pointer(bootParams))

	if err := C.firmware_parse(bootParams); err != 0 {
		return ErrInvalidBootParams
	}

	if err := C.xh_vm_create(); err != 0 {
		return ErrCreatingVM
	}

	maxVCPUs := C.num_vcpus_allowed()
	vcpus := C.int(p.VCPUs)
	if vcpus > maxVCPUs {
		return ErrMaxNumVCPUExceeded
	}

	var memsize C.size_t
	reqMemsize := C.CString(p.Memory)
	defer C.free(unsafe.Pointer(reqMemsize))
	if err := C.parse_memsize(reqMemsize, &memsize); err != 0 {
		return ErrInvalidMemsize
	}

	if err := C.xh_vm_setup_memory(memsize, C.VM_MMAP_ALL); err != 0 {
		return ErrSettingUpMemory
	}

	if err := C.init_msr(); err != 0 {
		return ErrInitializingMSR
	}

	C.init_mem()
	C.init_inout()
	C.pci_irq_init()
	C.ioapic_init()

	// Uses UTC by default.
	var rtcmode C.int
	if *p.RTCLocaltime {
		rtcmode = C.int(1)
	}
	C.rtc_init(rtcmode)
	C.sci_init()

	if err := C.init_pci(); err != 0 {
		return ErrInitializingPCI
	}

	if *p.BVMConsole {
		C.init_bvmcons()
	}

	if *p.MPTGen {
		if err := C.mptable_build(vcpus); err != 0 {
			return ErrBuildingMPTTable
		}
	}

	if err := C.smbios_build(); err != 0 {
		return ErrBuildingSMBIOS
	}

	if *p.ACPI {
		if err := C.acpi_build(vcpus); err != 0 {
			return ErrBuildingACPI
		}
	}

	var bsp C.int
	var rip C.uint64_t
	C.vcpu_add(bsp, bsp, rip)

	C.mevent_dispatch()

	return nil
}
