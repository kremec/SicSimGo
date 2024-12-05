prog    START  0

        . LOADING FROM MEMORY
        LDCH y          . Load character
        LDX #10         . Immediate addressing
        LDA x           . Direct addressing
                        . Indirect addressing
                        . Indexed addressing

        . STORING TO MEMORY
        STA z

        . OPERATIONS
        ADD y           . Addition
        RMO A,B         . Register to register
        DIVR B,A        . Division
                        . Bit operations

        . SUBROUTINES
        JSUB subrt      . Call subroutine and return from it

        . DISASSEMBLY FIX
        J io
testb   BYTE 0x00

        . IO
io      RD stdin
        WD outdat

halt    J halt


subrt   +LDA #16
        DIVR B,A
        STA z
        RSUB


. SYMBOLS
x       WORD 9
y       WORD 4
z       RESW 1

stdin   WORD 0x000000
outdat  WORD 0x0000AA

        END prog