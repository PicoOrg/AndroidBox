#ifndef _FUZZ_H_
#define _FUZZ_H_
#include <stdint.h>
#include <stdlib.h>
int LLVMFuzzerTestOneInput(char *data, size_t size);

extern void FuzzMain(char* data, int size);
#endif