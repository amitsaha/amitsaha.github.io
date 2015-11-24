/* Demo code showing the usage of the cleanup variable
   attribute. See:http://gcc.gnu.org/onlinedocs/gcc/Variable-Attributes.html
*/

/* Defines two cleanup functions to close and delete a temporary file
   and free a buffer
*/

# include <stdlib.h>
# include <stdio.h>

# define TMP_FILE "/tmp/tmp.file"

void free_buffer(char **buffer)
{
  printf("Freeing buffer\n");
  free(*buffer);
}

void cleanup_file(FILE **fp)
{
  printf("Closing file\n");
  fclose(*fp);

  printf("Deleting the file\n");
  remove(TMP_FILE);
}

int main(int argc, char **argv)
{
  char *buffer __attribute__ ((__cleanup__(free_buffer))) = malloc(20);
  FILE *fp __attribute__ ((__cleanup__(cleanup_file)));

  fp = fopen(TMP_FILE, "w+");
  
  if (fp != NULL)
    fprintf(fp, "%s", "Alinewithnospaces");

  fflush(fp);
  fseek(fp, 0L, SEEK_SET);
  fscanf(fp, "%s", buffer);
  printf("%s\n", buffer);
  
  return 0;
}


