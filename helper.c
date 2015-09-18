#include "helper.h"

#define MB (1024UL * 1024)

int firmware_parse(const char *opt) {
	char *fw, *opt1, *opt2, *opt3, *cp;

	fw = strdup(opt);

	if (strncmp(fw, "kexec", strlen("kexec")) == 0) {
		fw_func = kexec;
	} else if (strncmp(fw, "fbsd", strlen("fbsd")) == 0) {
		fw_func = fbsd_load;
	} else {
		goto fail;
	}

	if ((cp = strchr(fw, ',')) != NULL) {
		*cp = '\0';
		opt1 = cp + 1;
	} else {
		goto fail;
	}

	if ((cp = strchr(opt1, ',')) != NULL) {
		*cp = '\0';
		opt2 = cp + 1;
	} else {
		goto fail;
	}

	if ((cp = strchr(opt2, ',')) != NULL) {
		*cp = '\0';
		opt3 = cp + 1;
	} else {
		goto fail;
	}

	opt2 = strlen(opt2) ? opt2 : NULL;
	opt3 = strlen(opt3) ? opt3 : NULL;

	if (fw_func == kexec) {
		kexec_init(opt1, opt2, opt3);
	} else if (fw_func == fbsd_load) {
		/* FIXME: let user set boot-loader serial device */
		fbsd_init(opt1, opt2, opt3, NULL);
	} else {
		goto fail;
	}

	return 0;

fail:
	return -1;
}


int parse_memsize(const char *opt, size_t *ret_memsize)
{
	char *endptr;
	size_t optval;
	int error;

	optval = strtoul(opt, &endptr, 0);
	if (*opt != '\0' && *endptr == '\0') {
		/*
		 * For the sake of backward compatibility if the memory size
		 * specified on the command line is less than a megabyte then
		 * it is interpreted as being in units of MB.
		 */
		if (optval < MB)
			optval *= MB;
		*ret_memsize = optval;
		error = 0;
	} else
		error = expand_number(opt, ((uint64_t *) ret_memsize));

	return (error);
}

int num_vcpus_allowed(void)
{
	return (VM_MAXCPU);
}

int expand_number(const char *buf, uint64_t *num)
{
	char *endptr;
	uintmax_t umaxval;
	uint64_t number;
	unsigned shift;
	int serrno;

	serrno = errno;
	errno = 0;
	umaxval = strtoumax(buf, &endptr, 0);
	if (umaxval > UINT64_MAX)
		errno = ERANGE;
	if (errno != 0)
		return (-1);
	errno = serrno;
	number = umaxval;

	switch (tolower((unsigned char)*endptr)) {
	case 'e':
		shift = 60;
		break;
	case 'p':
		shift = 50;
		break;
	case 't':
		shift = 40;
		break;
	case 'g':
		shift = 30;
		break;
	case 'm':
		shift = 20;
		break;
	case 'k':
		shift = 10;
		break;
	case 'b':
	case '\0': /* No unit. */
		*num = number;
		return (0);
	default:
		/* Unrecognized unit. */
		errno = EINVAL;
		return (-1);
	}

	if ((number << shift) >> shift != number) {
		/* Overflow */
		errno = ERANGE;
		return (-1);
	}
	*num = number << shift;
	return (0);
}
