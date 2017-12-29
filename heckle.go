package main
import ( 
	"fmt"
	"os/exec"
	"log"
	"os"
	"flag"
	"strconv"
	"strings"
	"log/syslog"
	"regexp"
	//"encoding/csv"
)

// use to "test" exec_cmd will output the command passed to it, rather than actually execute it
// TODO make outputOnly a command line parameter. 
// TODO find a more elegant way ( ? : notation ) for if outputOnly { ...} on the print blah stuff
// const outputOnly=false

var version string
var outputOnly bool

func exec_cmd(cmd string, outputonly bool) string {
	if !outputonly  {
		log.Print("Running [" + cmd + "]")
		out, err :=exec.Command("/bin/bash", "-c", cmd).Output()
		output := strings.TrimSpace(string(out))
		log.Print("TRIMMED Output was [" + output + "]")
		if err != nil {
			// TODO roll back whats been done if a vital step fails.
			log.Fatal(err)
			os.Exit(1)
		}
		// trim off the CR at the end of any output
		return output
	} else {
		return cmd
	}
}

func createvm(template string, name string, cpus int,memory string, disksize string,network string, iso string) string {
	var vm_unwantedoutput =""
	if outputOnly {
		fmt.Println("outputting only")
	} else {
		fmt.Println("executing commands")
	}
	
	// STEP 0 validate the inputs
	sizeCheck := regexp.MustCompile(`[0-9]+[GMK]iB`)
	matches := sizeCheck.FindAllString(memory,-1)
	if matches == nil {
		fmt.Println("Memory must be a number followed by (G|M|K)iB")
		os.Exit(1)
	}


	matches = sizeCheck.FindAllString(disksize,-1)
	if matches == nil {
		fmt.Println("Disk size must be a number followed by (G|M|K)iB")
		os.Exit(1)
	}

	if cpus < 1 {
		fmt.Println("CPUs must be a positive integer")
		os.Exit(1)
	}

	// STEP 1
	// create the VM with the required template
	var vm_uuid =""
	vm_uuid = exec_cmd("xe vm-install template=\"" + template + "\" new-name-label=\"" + name + "\"",outputOnly)
	// once we've output the command, assign a dummy value to the variables we'll need later on, to make it look vaugely like a real run
	if outputOnly {
		fmt.Println(vm_uuid)
		vm_uuid="vm_uuid"
	}
	
	// STEP 2 now we need to disable booting from the automatically assigned disk
	var vm_disk_uuid = exec_cmd("xe vbd-list --minimal vm-uuid=" + vm_uuid + " userdevice=0",outputOnly)
	// get the uuid of the disk assigned to our new vm
	// now disable booting from it (we want to boot from the cdrom)
	if outputOnly {
		fmt.Println(vm_disk_uuid)
		vm_disk_uuid="vm_disk_uuid"
	}
	vm_unwantedoutput = exec_cmd("xe vbd-param-set uuid=" + vm_disk_uuid + " bootable=false",outputOnly)
	if outputOnly {
		fmt.Println(vm_unwantedoutput)
	}
	// STEP 3 make the drive the size we want it
	// we do this here because there will only be one VDI/VBD to parse
	// xe vbd-list params=vdi-uuid --minimal vm-uuid=
	var vm_vdi_uuid = exec_cmd("xe vbd-list params=vdi-uuid --minimal vm-uuid=" + vm_uuid,outputOnly)
	if outputOnly {
		fmt.Println(vm_vdi_uuid)
		vm_vdi_uuid="vm_vdi_uuid"
	}
	//xe vdi-resize uuid=[VDI uuid] disk-size=20GiB
	vm_unwantedoutput = exec_cmd("xe vdi-resize uuid=" + vm_vdi_uuid + " disk-size=" + disksize , outputOnly)
	if outputOnly {
		fmt.Println(vm_unwantedoutput)
	}

	// STEP 4 now we add a cd drive and "insert" the cd image
	
	vm_unwantedoutput = exec_cmd("xe vm-cd-add uuid=" + vm_uuid + " cd-name=\"" + iso + "\" device=3",outputOnly)
	if outputOnly {
		fmt.Println(vm_unwantedoutput)
	}

	// now we need to list device 3 (the cdrom) of the vm, and get the devices uuid, so we can enable booting on that.
	var vm_cd_uuid = exec_cmd("xe vbd-list --minimal vm-uuid=" + vm_uuid + " userdevice=3",outputOnly)

	
	if outputOnly {
    	fmt.Println(vm_cd_uuid)	
    	vm_cd_uuid="vm_cd_uuid"
    }
	//xe vbd-param-set  uuid=[device uuid] bootable=true
	vm_unwantedoutput = exec_cmd("xe vbd-param-set uuid=" + vm_cd_uuid  + " bootable=true", outputOnly)
	if outputOnly {
		fmt.Println(vm_unwantedoutput)	
	}
	// xe vm-param-set uuid=[VM uuid] other-config:install-repository=cdrom
	vm_unwantedoutput = exec_cmd("xe vm-param-set uuid=" + vm_uuid + " other-config:install-repository=cdrom", outputOnly)
	if outputOnly {
		fmt.Println(vm_unwantedoutput)
	}

	// STEP 5 now on to setting up the network
	//xe network-list --minimal name-label="CG 1072"	
	var vm_network_uuid = exec_cmd("xe network-list --minimal name-label=\"" + network + "\"", outputOnly)
	
	if outputOnly {
		fmt.Println(vm_network_uuid)
    	vm_network_uuid="vm_network_uuid"
    }
    
	//xe vif-create vm-uuid=[VM uuid] network-uuid=[network uuid] device=0
	vm_unwantedoutput = exec_cmd("xe vif-create vm-uuid=" + vm_uuid + " network-uuid=" + vm_network_uuid + " device=0", outputOnly)
	if outputOnly {	
		fmt.Println(vm_unwantedoutput)
	}
	
	// STEP 6 now set the RAM accordingly
	// xe vm-memory-limits-set dynamic-max=VM MEMORY dynamic-min=VM MEMORY static-max=VM MEMORY static-min=VM MEMORY name-label="newVM"
	vm_unwantedoutput = exec_cmd("xe vm-memory-limits-set dynamic-max=" + 
				memory + " dynamic-min=" + memory + " static-max=" + memory + 
					" static-min=" + memory + " uuid=" + vm_uuid, outputOnly)
	if outputOnly {	
		fmt.Println(vm_unwantedoutput)
	}

	// STEP 7 set the number of VCPUs
	vm_unwantedoutput = exec_cmd("xe vm-param-set uuid=" + vm_uuid + " platform:cores-per-socket=1 VCPUs-max=" + strconv.Itoa(cpus), outputOnly)
	if outputOnly {	
		fmt.Println(vm_unwantedoutput)
	}
	vm_unwantedoutput = exec_cmd("xe vm-param-set uuid=" + vm_uuid + " platform:cores-per-socket=1 VCPUs-at-startup=" + strconv.Itoa(cpus), outputOnly)
	if outputOnly {	
		fmt.Println(vm_unwantedoutput)
	}

	// STEP 8 set the boot parameters so that it goes and gets our preseed file and doesnt ask any questions
	// xe vm-param-set PV-args="auto priority=critical keymap=gb locale=en_GB hostname=preseedtest url=http://10.0.1.10/preseed-stretch.cfg -- quiet console=hvc0" vm=VMNAME
	vm_unwantedoutput = exec_cmd("xe vm-param-set PV-args=\"auto priority=critical keymap=gb locale=en_GB hostname=preseedtest url=http://10.0.1.10/preseed-stretch.cfg -- quiet console=hvc0\" uuid=" + vm_uuid, outputOnly)
	if outputOnly {	
		fmt.Println(vm_unwantedoutput)
	}
	
	return vm_uuid

	
}

func startvm(vm_uuid string) {
	var vm_unwantedoutput string
	vm_unwantedoutput = exec_cmd("xe vm-start uuid=" + vm_uuid, outputOnly)
	if outputOnly {	
		fmt.Println(vm_unwantedoutput)
	}
}

func main() {
	vmtemplatePtr		:= flag.String("template","","The Name of the XenServer template to use")
	vmnamePtr 			:= flag.String("name","blah","The name of the VM to create")
	vmcpusPtr 			:= flag.Int("cpus",1,"the number of vCPUs to assign to the VM")
	vmemoryPtr	 		:= flag.String("memory","1GiB","The number of MEGABYTES RAM to assign to the VM")
	vmdisksizePtr		:= flag.String("disksize","10GiB","The number of GiB to allocate to the disk")
	vmnetworkPtr		:= flag.String("network","","The name of the network to connect the vm to")
	vmisoPtr			:= flag.String("iso","","The name of the ISO from which to first-time-boot the vm")
	vmversionPtr		:= flag.Bool("version",false,"show the version and date of the build")
	vmstartPtr			:= flag.Bool("start",false,"If set it will start the vm once created")
	// vmmanifest			:= flag.String("manifest","","A CSV file containing the template,name,cpus,memory,disksize,network,iso values for multiple hosts")
	vmoutputonlyPtr		:= flag.Bool("outputonly",false,"If set it will only output the commands it would execute, naturally without the correct parameter values set.")

	flag.Parse()
	outputOnly = *vmoutputonlyPtr
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "heckle")

	if outputOnly {
		fmt.Println("outputonly commands will not execute")
	}

	if *vmversionPtr {
		fmt.Println("version " + version)
		os.Exit(0)
	}
	
    if e == nil {
        log.SetOutput(logwriter)
    }

    var vm_uuid string
    
    vm_uuid = createvm(*vmtemplatePtr, *vmnamePtr, *vmcpusPtr, *vmemoryPtr, *vmdisksizePtr, *vmnetworkPtr, *vmisoPtr)
    // READY!!!! 
	if *vmstartPtr {
		startvm(vm_uuid)
	}
}
