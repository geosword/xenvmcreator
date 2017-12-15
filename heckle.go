package main
import ( 
	"fmt"
	"os/exec"
	"log"
	"os"
	"flag"
)

// use to "test" exec_cmd will output the command passed to it, rather than actually execute it
const outputOnly=true

func exec_cmd(cmd string, outputonly bool) string {
	if !outputonly  {
		out, err :=exec.Command("/bin/bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		return string(out)
	} else {
		return cmd	
	}
}

func main() {
	// flag.StringVar(&flagvmname,"name","blah","The name of the VM to create")
	vmtemplatePtr		:= flag.String("template","","The Name of the XenServer template to use")
	vmnamePtr 			:= flag.String("name","blah","The name of the VM to create")
	//vmcpusPtr 		:= flag.Int("cpus",1,"the number of vCPUs to assign to the VM")
	vmemoryPtr	 		:= flag.String("memory","1GiB","The number of MEGABYTES RAM to assign to the VM")
	//vmdisksizePtr		:= flag.Int("disksize",20,"The number of GiB to allocate to the disk")
	vmnetworkPtr		:= flag.String("network","","The name of the network to connect the vm to")
	vmisoPtr			:= flag.String("iso","","The name of the ISO from which to first-time-boot the vm")

	flag.Parse()
	var vm_unwantedoutput =""
	
	// STEP 1
	// create the VM with the required template
	var vm_uuid =""
	vm_uuid = exec_cmd("xe vm-install template=\"" + *vmtemplatePtr + "\" new-name-label=\"" + *vmnamePtr +"\"",outputOnly)
	fmt.Println(vm_uuid)
	// once we've output the command, assign a dummy value to the variables we'll need later on, to make it look vaugely like a real run
	if outputOnly {
		vm_uuid="vm_uuid"
	}
	
	// STEP 2 now we need to disable booting from the automatically assigned disk
	var vm_disk_uuid =""
	var vmcmd_disableboot =""
	// get the uuid of the disk assigned to our new vm
	vm_disk_uuid = exec_cmd("xe vbd-list vm-uuid=" + vm_uuid + " userdevice=0",outputOnly)
	fmt.Println(vm_disk_uuid)
	// now disable booting from it (we want to boot from the cdrom)
	if outputOnly {
		vm_disk_uuid="vm_disk_uuid"
	}	
	vmcmd_disableboot = "xe vbd-param-set uuid=" + vm_disk_uuid + " bootable=false"
	vm_unwantedoutput = exec_cmd(vmcmd_disableboot,outputOnly)
	fmt.Println(vm_unwantedoutput)
	
	// STEP 3 now we add a cd drive and "insert" the cd image
	var vm_cd_uuid =""
	vm_unwantedoutput = exec_cmd("xe vm-cd-add vm=\"" + *vmnamePtr + "\" cd-name=\"" + *vmisoPtr + "\" device=3",outputOnly)
	fmt.Println(vm_unwantedoutput)

	// now we need to list device 3 (the cdrom) of the vm, and get the devices uuid, so we can enable booting on that.
	vm_cd_uuid = exec_cmd("xe vbd-list vm-name-label=\"" + *vmnamePtr + "\" userdevice=3",outputOnly)
	fmt.Println(vm_cd_uuid)	
	if outputOnly {
    	vm_cd_uuid="vm_cd_uuid"
    }
	//xe vbd-param-set  uuid=[device uuid] bootable=true
	vm_unwantedoutput = exec_cmd("xe vbd-param-set  uuid=\"" + vm_cd_uuid  + "\" bootable=true", outputOnly)
	fmt.Println(vm_unwantedoutput)	
	// xe vm-param-set uuid=[VM uuid] other-config:install-repository=cdrom
	vm_unwantedoutput = exec_cmd("xe vm-param-set uuid=" + vm_uuid + " other-config:install-repository=cdrom", outputOnly)
	fmt.Println(vm_unwantedoutput)
	

	// STEP 4 now on to setting up the network
	//xe network-list --minimal name-label="CG 1072"	
	var vm_network_uuid =""
	vm_network_uuid = exec_cmd("xe network-list --minimal name-label=\"" + *vmnetworkPtr + "\"", outputOnly)
	if outputOnly {
    	vm_network_uuid="vm_network_uuid"
    }
	//xe vif-create vm-uuid=[VM uuid] network-uuid=[network uuid] device=0
	vm_unwantedoutput = exec_cmd("xe vif-create vm-uuid=\"" + vm_uuid + "\" network-uuid=\"" + vm_network_uuid + "\" device=0", outputOnly)
	fmt.Println(vm_unwantedoutput)
	
	// STEP 5 now set the RAM accordingly
	// xe vm-memory-limits-set dynamic-max=VM MEMORY dynamic-min=VM MEMORY static-max=VM MEMORY static-min=VM MEMORY name-label="newVM"
	vm_unwantedoutput = exec_cmd("xe vm-memory-limits-set dynamic-max=" + 
			*vmemoryPtr + " dynamic-min=" + *vmemoryPtr + " static-max=" + *vmemoryPtr + 
			" static-min=" + *vmemoryPtr + " name-label=\"" + *vmnamePtr + "\"", outputOnly)
	fmt.Println(vm_unwantedoutput)
}
