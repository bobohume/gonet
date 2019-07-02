package main

// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
/*
void print(char *str) {
    printf("%s\n", str);
    char s[]="Golden Global View";
	memmove(s,s+7,strlen(s)+1-7);
	printf("%s",s);
}
*/
import "C"
import "unsafe"

func main() {
	s := "Hello Cgo"
	cs := C.CString(s)
	C.print(cs)
	C.free(unsafe.Pointer(cs))
}
