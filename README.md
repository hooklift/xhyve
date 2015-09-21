# Experimental Go bindings for xhyve.
Work in progress.

### Building
0. git submodule --init 
1. make patch-xhyve
2. make
3. sudo ./nnn -m 1024M -c 1 -A -s 0:0,hostbridge -s 31,lpc -l com1,stdio -s 2:0,virtio-net -U 6BCE442E-4359-4BD9-84F7-EDFB8EC6D2EF -f 'kexec,imgs/stable.766.3.0.coreos_production_pxe.vmlinuz,imgs/stable.766.3.0.coreos_production_pxe_image.cpio.gz,earlyprintk=serial console=ttyS0 coreos.autologin'
