package tidy

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -ltidy -L/usr/local/lib/libtidy.a
#include <tidy.h>
#include <buffio.h>
#include <errno.h>
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func Tidy(htmlSource string) (string, os.Error) {
	input := C.CString(htmlSource)
	defer C.free(unsafe.Pointer(input))

	var output C.TidyBuffer
	defer C.tidyBufFree( &output )

	var errbuf C.TidyBuffer
  	defer C.tidyBufFree( &errbuf )

	var rc C.int = -1
	var ok C.Bool

	var tdoc C.TidyDoc = C.tidyCreate()	// Initialize "document"
	defer C.tidyRelease( tdoc )

	ok = C.tidyOptSetBool( tdoc, C.TidyXhtmlOut, C.yes )  // Convert to XHTML

	if ok == 1 {
		rc = C.tidySetErrorBuffer( tdoc, &errbuf )	// Capture diagnostics
	}

	if rc >= 0 {
		in := _Ctypedef_tmbchar(*input)
    	rc = C.tidyParseString( tdoc, &in )	// Parse the input
    }

	if rc >= 0 {
	    rc = C.tidyCleanAndRepair( tdoc )	// Tidy it up!	
	}

	if rc >= 0 {
    	rc = C.tidyRunDiagnostics( tdoc )	// Kvetch
    }

	if rc > 1 {		// If error, force output.
		if C.tidyOptSetBool(tdoc, C.TidyForceOutput, C.yes) == 0 {
			rc = -1
		}
	}

	if rc >= 0 {
    	rc = C.tidySaveBuffer( tdoc, &output )	// Pretty Print
    }

	if rc >= 0 {
    	out := _Ctype_char(*output.bp)
    	if rc > 0 {
    		err := _Ctype_char(*errbuf.bp)
    		return C.GoString(&out), os.NewError(C.GoString(&err))
      	}
		return C.GoString(&out), nil
  	}
    return "", os.NewSyscallError(fmt.Sprintf( "A severe error (%d) occurred.\n", int(rc) ), int(rc))
}