#include "xhyve.h"

#define ECREATINGVM   1
#define ESETUPMEM     2
#define EINITMSR      3
#define EINITPCI      4
#define EMPTABLEBUILD 5

static int acpi;
static const int BSP = 0;

int
run_xhyve(int ncpus, size_t memsize)
{
	int error, bvmcons, mptgen;
	uint64_t rip;
	int rtc_localtime;

	bvmcons = 0;

	error = xh_vm_create();
	if (error) {
		fprintf(stderr, "Unable to create VM (%d)\n", error);
		return ECREATINGVM;
	}

	error = xh_vm_setup_memory(memsize, VM_MMAP_ALL);
	if (error) {
		fprintf(stderr, "Unable to setup memory (%d)\n", error);
		return ESETUPMEM;
	}

	error = init_msr();
	if (error) {
		fprintf(stderr, "init_msr error %d\n", error);
		return EINITMSR;
	}

	init_mem();
	init_inout();
	pci_irq_init();
	ioapic_init();

	rtc_init(rtc_localtime);
	sci_init();

	/*
	 * Exit if a device emulation finds an error in it's initilization
	 */
	if (init_pci() != 0) {
		return EINITPCI;
	}

	if (bvmcons) {
		init_bvmcons();
	}

	/*
	 * build the guest tables, MP etc.
	 */
	if (mptgen) {
		error = mptable_build(ncpus);
		if (error) {
			return EMPTABLEBUILD;
		}
	}

	error = smbios_build();
	assert(error == 0);

	if (acpi) {
		error = acpi_build(ncpus);
		assert(error == 0);
	}

	rip = 0;

	vcpu_add(BSP, BSP, rip);

	/*
	 * Head off to the main event dispatch loop
	 */
	mevent_dispatch();

	return 0;
}
