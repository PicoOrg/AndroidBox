#ifndef _FUZZ_H_
#define _FUZZ_H_
#include <stdint.h>
#include <stdlib.h>
int LLVMFuzzerTestOneInput(const uint8_t *data, size_t size);

#endif