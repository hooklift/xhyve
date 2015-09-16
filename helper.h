#ifndef helpers_h
#define helpers_h 1

#include <errno.h>
#include <inttypes.h>
#include <ctype.h>
#include <stdlib.h>

#include <xhyve/acpi.h>
#include <xhyve/vmm/vmm_api.h>
#include <xhyve/inout.h>
#include <xhyve/mem.h>
#include <xhyve/xmsr.h>
#include <xhyve/pci_emul.h>
#include <xhyve/ioapic.h>
#include <xhyve/pci_lpc.h>
#include <xhyve/mevent.h>
#include <xhyve/mptbl.h>
#include <xhyve/pci_irq.h>
#include <xhyve/rtc.h>
#include <xhyve/smbiostbl.h>
#include <xhyve/xhyve.h>

#include <xhyve/firmware/kexec.h>
#include <xhyve/firmware/fbsd.h>

int firmware_parse(const char *opt);
int parse_memsize(const char *opt, size_t *ret_memsize);
int num_vcpus_allowed(void);
int expand_number(const char *buf, uint64_t *num);
#endif
