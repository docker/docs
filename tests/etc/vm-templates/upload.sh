#/bin/sh

upload_1010() {
	ovftool \
		--name=test-osx-10.10-template-vm \
		--datastore=Templates \
		--diskMode=thick \
		--noSSLVerify \
		--acceptAllEulas \
		--overwrite \
		--allowExtraConfig \
		--extraConfig:smc.present=TRUE \
	packer/output-osx-10.10/packer-osx-10.10.vmx \
		vi://administrator%40vsphere.local:Docker4theTeam%21@172.16.1.20/Datacenter/host/Cluster1
}

upload_1011(){
	ovftool \
		--name=test-osx-10.11-template-vm \
		--datastore=Templates \
		--diskMode=thick \
		--noSSLVerify \
		--acceptAllEulas \
		--overwrite \
		--allowExtraConfig \
		--extraConfig:smc.present=TRUE \
		packer/output-osx-10.11/packer-osx-10.11.vmx \
		vi://administrator%40vsphere.local:Docker4theTeam%21@172.16.1.20/Datacenter/host/Cluster1
}

upload_1012(){
	ovftool \
		--name=test-osx-10.12-template-vm \
		--datastore=Templates \
		--diskMode=thick \
		--noSSLVerify \
		--acceptAllEulas \
		--overwrite \
		--allowExtraConfig \
		--extraConfig:smc.present=TRUE \
		packer/output-osx-10.12/packer-osx-10.12.vmx \
		vi://administrator%40vsphere.local:Docker4theTeam%21@172.16.1.20/Datacenter/host/Cluster1
}

if [ -z "$1" ]; then
	upload_1010
	upload_1011
	upload_1012
fi

case "$1" in
	"10.10")
		upload_1010
		;;
	"10.11")
		upload_1011
		;;
	"10.12")
		upload_1012
		;;
	*)
		echo "Unsupported release"
		exit 1
esac

