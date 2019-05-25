#include <cups/cups.h>
#include <stdio.h>

int main() {
  int job_id;
  printf("%s\n", cupsGetDefault());

  /* Print a single file */
  job_id = cupsPrintFile(cupsGetDefault(), "./test.txt", "Test Print", 0, NULL);

  if (job_id == 0) puts(cupsLastErrorString());

  printf("job id: %d\n", job_id);

  return 0;
}
