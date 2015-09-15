#include <stdio.h>
#include <stdint.h>

#include <xhyve/support/misc.h>
#include <xhyve/support/atomic.h>
#include <xhyve/support/segments.h>
#include <xhyve/support/cpuset.h>
#include <xhyve/vmm/vmm_api.h>

#include <xhyve/xhyve.h>
#include <xhyve/acpi.h>
#include <xhyve/inout.h>
#include <xhyve/dbgport.h>
#include <xhyve/ioapic.h>
#include <xhyve/mem.h>
#include <xhyve/mevent.h>
#include <xhyve/mptbl.h>
#include <xhyve/pci_emul.h>
#include <xhyve/pci_irq.h>
#include <xhyve/pci_lpc.h>
#include <xhyve/smbiostbl.h>
#include <xhyve/xmsr.h>
#include <xhyve/rtc.h>

#include <xhyve/firmware/kexec.h>
#include <xhyve/firmware/fbsd.h>

int
run_xhyve(int ncpus, size_t memsize);
