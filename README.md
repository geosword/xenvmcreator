# xenvmcreator
A go script to create VMs on Citrix Xenserver Hypervisors

Usage of /tmp/go-build085687383/command-line-arguments/_obj/exe/heckle:
  -cpus int
        the number of vCPUs to assign to the VM (default 1)
  -disksize string
        The number of GiB to allocate to the disk (default "10GiB")
  -iso string
        The name of the ISO from which to first-time-boot the vm
  -memory string
        The number of MEGABYTES RAM to assign to the VM (default "1GiB")
  -name string
        The name of the VM to create (default "blah")
  -network string
        The name of the network to connect the vm to
  -outputonly
        If set it will only output the commands it would execute, naturally without the correct parameter values set.
  -start
        If set it will start the vm once created
  -template string
        The Name of the XenServer template to use
  -version
        show the version and date of the build
