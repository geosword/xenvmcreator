package main
import ( 
	"fmt"
	"os/exec"
	"log"
	"os"
	"flag"
)

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
	vmnamePtr 		:= flag.String("name","blah","The name of the VM to create")
	//vmcpusPtr 		:= flag.Int("cpus",1,"the number of vCPUs to assign to the VM")
	//vmemoryPtr 		:= flag.Int("memory",1024,"The number of MEGABYTES RAM to assign to the VM")
	//vmdisksizePtr		:= flag.Int("disksize",20,"The number of GiB to allocate to the disk")
	//vmnetworkPtr		:= flag.String("network","","The name of the network to connect the vm to")
	//vmisoPtr		:= flag.String("iso","","The name of the ISO from which to first-time-boot the vm")

	

	//fmt.Println("hello world")
	
	flag.Parse()
	fmt.Println("vmname is :", *vmnamePtr )
	var vmcmd_install="xe vm-install template=\"" + *vmtemplatePtr + "\" new-name-label=\"" + *vmnamePtr +"\""
	output :=""
	output = exec_cmd(vmcmd_install,true)
	fmt.Println(output)
}
