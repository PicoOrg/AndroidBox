#include "fuzz.h"

int LLVMFuzzerTestOneInput(char *data, size_t size) {
    FuzzMain(data, size);
    return 0;
}